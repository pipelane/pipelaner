/*
 * Copyright (c) 2024 Alexey Khokhlov
 */

package pipelaner

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/pipelane/pipelaner/gen/source/input"
	"github.com/pipelane/pipelaner/pipeline/components"
	"github.com/pipelane/pipelaner/pipeline/source"
	"github.com/pipelane/pipelaner/sources/generator/pipelaner/server"
	"github.com/pipelane/pipelaner/sources/shared/grpc_server"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"github.com/pipelane/pipelaner/sources/shared/proto/service"
)

type Pipelaner struct {
	components.Logger
	srv *server.PipelanerServer
}

func init() {
	source.RegisterInput("pipelaner", &Pipelaner{})
}

func (p *Pipelaner) Init(cfg input.Input) error {
	var opts []grpc.ServerOption
	c, ok := cfg.(input.Pipelaner)
	if !ok {
		return errors.New("invalid input config type")
	}
	if tls := c.GetCommonConfig().Tls; tls != nil {
		cred, errs := credentials.NewServerTLSFromFile(tls.CertFile, tls.KeyFile)
		if errs != nil {
			p.Log().Fatal().Err(errs).Msg("Failed to generate credentials")
		}
		opts = []grpc.ServerOption{grpc.Creds(cred)}
	}
	l := p.Log().With().Logger()
	serv := grpc_server.NewServer(&grpc_server.ServerConfig{
		Host:           c.GetCommonConfig().Host,
		Port:           c.GetCommonConfig().Port,
		ConnectionType: c.GetConnectionType().String(),
		UnixSocketPath: c.GetUnixSocketPath(),
		Opts:           opts,
	}, &l)
	p.srv = server.NewServer(&l, c.GetOutputBufferSize())
	err := serv.Serve(func(s *grpc.Server) {
		service.RegisterPipelanerServer(s, p.srv)
	})
	return err
}

func (p *Pipelaner) Generate(ctx context.Context, input chan<- any) {
	for m := range p.srv.Recv() {
		select {
		case <-ctx.Done():
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
					p.Log().Error().Err(err).Msg("Error invalid unmarshal json data")
					continue
				}
				input <- v
			}
		}
	}
}
