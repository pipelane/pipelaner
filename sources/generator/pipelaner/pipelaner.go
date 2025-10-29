/*
 * Copyright (c) 2024 Alexey Khokhlov
 */

package pipelaner

import (
	"context"
	"encoding/json"
	"fmt"

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
	c, ok := cfg.(input.Pipelaner)
	if !ok {
		return fmt.Errorf("invalid input config type: %T", cfg)
	}

	var (
		opts []grpc.ServerOption
		host string
		port uint
	)

	if commonCfg := c.GetCommonConfig(); commonCfg != nil {
		host, port = commonCfg.Host, commonCfg.Port
		if tls := commonCfg.Tls; tls != nil {
			cred, err := credentials.NewServerTLSFromFile(tls.CertFile, tls.KeyFile)
			if err != nil {
				p.Log().Fatal().Err(err).Msg("generate tls credentials")
			}
			opts = append(opts, grpc.Creds(cred))
		}
	}
	l := p.Log().With().Logger()
	serv := grpc_server.NewServer(&grpc_server.ServerConfig{
		Host:           host,
		Port:           port,
		ConnectionType: c.GetConnectionType().String(),
		UnixSocketPath: c.GetUnixSocketPath(),
		Opts:           opts,
	}, &l)
	p.srv = server.NewServer(&l, c.GetOutputBufferSize())
	go func() {
		err := serv.Serve(context.Background(), func(s *grpc.Server) {
			service.RegisterPipelanerServer(s, p.srv)
		})
		if err != nil {
			p.Log().Fatal().Err(err).Msg("Failed to start server")
		}
	}()
	return nil
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
