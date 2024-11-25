/*
 * Copyright (c) 2024 Alexey Khokhlov
 */

package pipelaner

/*import (
	"github.com/goccy/go-json"
	"github.com/pipelane/pipelaner/source/generator/pipelaner/server"
	grpc_server "github.com/pipelane/pipelaner/source/shared/grpc_server"

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
	cfg *pipelaner.BaseLaneConfig
	srv *server.PipelanerServer
}

func init() {
	pipelaner.RegisterGenerator("pipelaner", &Pipelaner{})
}

func (p *Pipelaner) Init(ctx *pipelaner.Context) error {
	p.cfg = ctx.LaneItem().Config()

	v := &Config{}
	err := p.cfg.ParseExtended(v)
	if err != nil {
		return err
	}
	host := "localhost"
	if v.Host != nil {
		host = *v.Host
	}

	l := ctx.Logger()
	var opts []grpc.ServerOption
	if v.TLS {
		cred, errs := credentials.NewServerTLSFromFile(v.CertFile, v.KeyFile)
		if errs != nil {
			l.Fatal().Err(errs).Msg("Failed to generate credentials")
		}
		opts = []grpc.ServerOption{grpc.Creds(cred)}
	}

	serv := grpc_server.NewServer(&grpc_server.ServerConfig{
		Host:           host,
		Port:           v.Port,
		ConnectionType: v.ConnectionType,
		UnixSocketPath: v.UnixSocketPath,
		Opts:           opts,
	}, &l)
	p.srv = server.NewServer(&l, p.cfg.OutputBufferSize)
	serv.Serve(func(s *grpc.Server) {
		service.RegisterPipelanerServer(s, p.srv)
	})

	return nil
}

func (p *Pipelaner) Generate(ctx *pipelaner.Context, input chan<- any) {
	l := ctx.Logger()
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
					l.Error().Err(err).Msg("Error invalid unmarshal json data")
					continue
				}
				input <- v
			}
		}
	}
}*/
