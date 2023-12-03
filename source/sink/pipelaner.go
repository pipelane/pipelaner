/*
 * Copyright (c) 2023 Alexey Khokhlov
 */

package sink

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/pipelane/pipelaner"
	"github.com/pipelane/pipelaner/internal/service"
)

var (
	ErrInvalidDataFormat = errors.New("ErrInvalidDataFormat")
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
	client service.PipelanerClient
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
	var opts []grpc.DialOption
	if v.Tls {
		creds, err := credentials.NewClientTLSFromFile(v.CertFile, "")
		if err != nil {
			p.logger.Fatal().Msgf("Failed to create TLS credentials: %v", err)
		}
		opts = append(opts, grpc.WithTransportCredentials(creds))
	} else {
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}
	conn, err := grpc.Dial(fmt.Sprintf("%s:%d", host, v.Port), opts...)
	if err != nil {
		return err
	}
	client := service.NewPipelanerClient(conn)
	p.client = client

	return nil
}

func (p *Pipelaner) Sink(ctx context.Context, val any) {
	var m *service.Message
	switch v := val.(type) {
	case string:
		m = &service.Message{
			Data: &service.Message_StringValue{
				StringValue: v,
			},
		}
	case []byte:
		m = &service.Message{
			Data: &service.Message_BytesValue{
				BytesValue: v,
			},
		}
	default:
		b, err := json.Marshal(v)
		if err != nil {
			p.logger.Error().Err(err).Msg("Grpc sink failed")
			return
		}
		m = &service.Message{
			Data: &service.Message_BytesValue{
				BytesValue: b,
			},
		}
	}
	_, err := p.client.Sink(ctx, m)
	if err != nil {
		p.logger.Error().Err(err).Msg("Grpc sing failed")
	}
}
