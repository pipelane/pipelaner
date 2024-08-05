/*
 * Copyright (c) 2023 Alexey Khokhlov
 */

package server

import (
	"context"

	"github.com/rs/zerolog"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/pipelane/pipelaner/source/shared/proto/service"
)

type PipelanerServer struct {
	service.UnimplementedPipelanerServer
	logger zerolog.Logger
	buffer chan *service.Message
}

func (s *PipelanerServer) Sink(_ context.Context, message *service.Message) (*emptypb.Empty, error) {
	s.buffer <- message
	return &emptypb.Empty{}, nil
}

func NewServer(logger zerolog.Logger, bufferSize int64) *PipelanerServer {
	s := &PipelanerServer{logger: logger, buffer: make(chan *service.Message, bufferSize)}
	return s
}

func (s *PipelanerServer) Recv() <-chan *service.Message {
	return s.buffer
}
