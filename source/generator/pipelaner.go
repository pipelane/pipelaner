/*
 * Copyright (c) 2023 Alexey Khokhlov
 */

package generator

import (
	"context"
	"fmt"
	"net"

	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"github.com/pipelane/pipelaner"
	"github.com/pipelane/pipelaner/internal/service"
	"github.com/pipelane/pipelaner/source/generator/server"
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

func (p *Pipelaner) Init(cfg *pipelaner.BaseLaneConfig) error {
	p.cfg = cfg
	p.logger = pipelaner.NewLogger()
	v := &GrpcCfg{}
	err := cfg.ParseExtended(v)
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

func (p *Pipelaner) Generate(ctx context.Context, input chan<- any) {
	for m := range p.srv.Recv() {
		select {
		case <-ctx.Done():
			break
		default:
			input <- m
		}
	}
}
