/*
 * Copyright (c) 2024 Alexey Khokhlov
 */

package pipelaner

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/pipelane/pipelaner/gen/source/sink"
	"github.com/pipelane/pipelaner/pipeline/components"
	"github.com/pipelane/pipelaner/pipeline/node"
	"github.com/pipelane/pipelaner/pipeline/source"
	service2 "github.com/pipelane/pipelaner/sources/shared/proto/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

func init() {
	source.RegisterSink("pipelaner", &Pipelaner{})
}

type Pipelaner struct {
	components.Logger
	client service2.PipelanerClient
}

func (p *Pipelaner) Init(cfg sink.Sink) error {
	pipelanerCfg, ok := cfg.(sink.Pipelaner)
	if !ok {
		return fmt.Errorf("invalid pipelaner config type: %T", cfg)
	}
	var opts []grpc.DialOption
	if tls := pipelanerCfg.GetCommonConfig().Tls; tls != nil {
		cred, err := credentials.NewClientTLSFromFile(tls.CertFile, "")
		if err != nil {
			return err
		}
		opts = append(opts, grpc.WithTransportCredentials(cred))
	} else {
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}
	conn, err := grpc.NewClient(
		fmt.Sprintf(
			"%s:%d",
			pipelanerCfg.GetCommonConfig().Host,
			pipelanerCfg.GetCommonConfig().Port,
		), opts...)
	if err != nil {
		return err
	}
	client := service2.NewPipelanerClient(conn)
	p.client = client
	return nil
}

func (p *Pipelaner) Sink(val any) error {
	var m *service2.Message
	switch v := val.(type) {
	case node.AtomicMessage:
		err := p.Sink(v.Data())
		if err != nil {
			v.Error() <- v
			return err
		}
		v.Success() <- v
		return nil
	case string:
		m = &service2.Message{
			Data: &service2.Message_StringValue{
				StringValue: v,
			},
		}
	case []byte:
		m = &service2.Message{
			Data: &service2.Message_BytesValue{
				BytesValue: v,
			},
		}
	default:
		b, err := json.Marshal(v)
		if err != nil {
			p.Log().Error().Err(err).Msg("Grpc sink failed")
			return err
		}
		m = &service2.Message{
			Data: &service2.Message_JsonValue{
				JsonValue: b,
			},
		}
	}
	_, err := p.client.Sink(context.Background(), m)
	if err != nil {
		p.Log().Error().Err(err).Msg("Grpc sing failed")
		return err
	}
	return nil
}
