/*
 * Copyright (c) 2023 Alexey Khokhlov
 */

package generator

import (
	"fmt"
	"net"

	"github.com/goccy/go-json"

	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"pipelaner"
	"pipelaner/internal/service"
	"pipelaner/source/generator/server"
)

type GrpcCfg struct {
	Host     *string `pipelane:"host"`
	Port     int     `pipelane:"port"`
	Tls      bool    `pipelane:"tls"`
	CertFile string  `pipelane:"cert"`
	KeyFile  string  `pipelane:"key"`
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
	v := &GrpcCfg{}
	err := p.cfg.ParseExtended(v)
	if err != nil {
		return err
	}
	host := "localhost"
	if v.Host != nil {
		host = *v.Host
	}
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", host, v.Port))
	if err != nil {
		p.logger.Fatal().Err(err).Msgf("Failed to listen %s:%d", host, v.Port)
	}
	var opts []grpc.ServerOption
	if v.Tls {
		cred, err := credentials.NewServerTLSFromFile(v.CertFile, v.KeyFile)
		if err != nil {
			p.logger.Fatal().Err(err).Msg("Failed to generate credentials")
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
