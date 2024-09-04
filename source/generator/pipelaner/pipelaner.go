/*
 * Copyright (c) 2024 Alexey Khokhlov
 */

package pipelaner

import (
	"fmt"
	"net"
	"os"
	"syscall"

	"github.com/pipelane/pipelaner/source/generator/pipelaner/server"

	"github.com/goccy/go-json"

	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"github.com/pipelane/pipelaner"
	"github.com/pipelane/pipelaner/source/shared/proto/service"
)

type Config struct {
	Host           *string `pipelane:"host"`
	Port           int     `pipelane:"port"`
	TLS            bool    `pipelane:"tls"`
	CertFile       string  `pipelane:"cert"`
	KeyFile        string  `pipelane:"key"`
	ConnectionType string  `pipelane:"connection_type"`
	UnixSocketPath string  `pipelane:"unix_socket_path"`
}

type Pipelaner struct {
	logger zerolog.Logger
	cfg    *pipelaner.BaseLaneConfig
	srv    *server.PipelanerServer
}

func init() {
	pipelaner.RegisterGenerator("pipelaner", &Pipelaner{})
}

func (p *Pipelaner) Init(ctx *pipelaner.Context) error {
	p.cfg = ctx.LaneItem().Config()
	p.logger = pipelaner.NewLogger()
	v := &Config{}
	err := p.cfg.ParseExtended(v)
	if err != nil {
		return err
	}
	host := "localhost"
	if v.Host != nil {
		host = *v.Host
	}
	lis := p.createListener(v, host)
	var opts []grpc.ServerOption
	if v.TLS {
		cred, errs := credentials.NewServerTLSFromFile(v.CertFile, v.KeyFile)
		if errs != nil {
			p.logger.Fatal().Err(errs).Msg("Failed to generate credentials")
		}
		opts = []grpc.ServerOption{grpc.Creds(cred)}
	}

	grpcServer := grpc.NewServer(opts...)
	p.srv = server.NewServer(p.logger, p.cfg.BufferSize)
	service.RegisterPipelanerServer(grpcServer, p.srv)
	go func() {
		err = grpcServer.Serve(lis)
		if err != nil {
			p.logger.Fatal().Err(err).Msg("Failed run server")
		}
	}()
	return nil
}

func (p *Pipelaner) Generate(ctx *pipelaner.Context, input chan<- any) {
	for m := range p.srv.Recv() {
		select {
		case <-ctx.Context().Done():
			break
		default:
			switch {
			case m.GetBytesValue() != nil:
				input <- m.GetBytesValue()
			case m.GetStringValue() != "":
				input <- m.GetStringValue()
			case m.GetJsonValue() != nil:
				var v map[string]any
				err := json.Unmarshal(m.GetJsonValue(), &v)
				if err != nil {
					p.logger.Error().Err(err).Msg("Error invalid unmarshal json data")
					continue
				}
				input <- v
			}
		}
	}
}

func (p *Pipelaner) createListener(v *Config, host string) net.Listener {
	if v.ConnectionType == "unix" {
		if err := syscall.Unlink(v.UnixSocketPath); err != nil && !os.IsNotExist(err) {
			p.logger.Fatal().Err(err).Msgf("Failed to unlink Unix socket %s", v.UnixSocketPath)
		}

		unixListener, err := net.ListenUnix("unix", &net.UnixAddr{Name: v.UnixSocketPath, Net: "unix"})
		if err != nil {
			p.logger.Fatal().Err(err).Msgf("Failed to listen on Unix socket %s", v.UnixSocketPath)
		}
		return unixListener
	}

	tcpListener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", host, v.Port))
	if err != nil {
		p.logger.Fatal().Err(err).Msgf("Failed to listen on TCP %s:%d", host, v.Port)
	}
	return tcpListener
}
