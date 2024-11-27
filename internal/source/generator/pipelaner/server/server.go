/*
 * Copyright (c) 2023 Alexey Khokhlov
 */

package server

import (
	"context"
	"errors"
	"io"

	service2 "github.com/pipelane/pipelaner/internal/source/shared/proto/service"
	"github.com/rs/zerolog"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type PipelanerServer struct {
	service2.UnimplementedPipelanerServer
	logger *zerolog.Logger
	buffer chan *service2.Message
}

func (s *PipelanerServer) Sink(_ context.Context, message *service2.Message) (*emptypb.Empty, error) {
	s.buffer <- message
	return &emptypb.Empty{}, nil
}

func (s *PipelanerServer) SinkStream(stream service2.Pipelaner_SinkStreamServer) error {
	for {
		message, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			s.logger.Error().Err(err).Msg("Error receiving stream request")
			return status.Errorf(codes.Internal, "cannot receive stream request: %v", err)
		}
		s.buffer <- message
	}
	return nil
}

func NewServer(logger *zerolog.Logger, bufferSize int64) *PipelanerServer {
	s := &PipelanerServer{logger: logger, buffer: make(chan *service2.Message, bufferSize)}
	return s
}

func (s *PipelanerServer) Recv() <-chan *service2.Message {
	return s.buffer
}
