/*
 * Copyright (c) 2024 Alexey Khokhlov
 */

package pipelaner

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/pipelane/pipelaner/gen/source/sink"
	"github.com/pipelane/pipelaner/internal/pipeline/source"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/pipelane/pipelaner/internal/source/shared/proto/service"
)

func init() {
	source.RegisterSink("pipelaner", &Pipelaner{})
}

type Pipelaner struct {
	logger *zerolog.Logger
	client service.PipelanerClient
}

func (p *Pipelaner) Init(cfg sink.Sink) error {
	pipelanerCfg, ok := cfg.(sink.Pipelaner)
	if !ok {
		return fmt.Errorf("invalid pipelaner config type: %T", cfg)
	}
	var opts []grpc.DialOption
	if pipelanerCfg.GetTls() {
		creds, err := credentials.NewClientTLSFromFile(*pipelanerCfg.GetCertFile(), "")
		if err != nil {
			return err
		}
		opts = append(opts, grpc.WithTransportCredentials(creds))
	} else {
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}
	conn, err := grpc.NewClient(fmt.Sprintf("%s:%d", pipelanerCfg.GetHost(), pipelanerCfg.GetPort()), opts...)
	if err != nil {
		return err
	}
	client := service.NewPipelanerClient(conn)
	p.client = client
	return nil
}

func (p *Pipelaner) Sink(val any) {
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
			Data: &service.Message_JsonValue{
				JsonValue: b,
			},
		}
	}
	_, err := p.client.Sink(context.Background(), m)
	if err != nil {
		p.logger.Error().Err(err).Msg("Grpc sing failed")
	}
}
