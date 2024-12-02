/*
 * Copyright (c) 2023 Alexey Khokhlov
 */

package server

import (
	"context"
	"errors"
	"io"

	"github.com/pipelane/pipelaner/source/shared/proto/service"
	"github.com/rs/zerolog"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type PipelanerServer struct {
	service.UnimplementedPipelanerServer
	logger *zerolog.Logger
	buffer chan *service.Message
}

func (s *PipelanerServer) Sink(_ context.Context, message *service.Message) (*emptypb.Empty, error) {
	s.buffer <- message
	return &emptypb.Empty{}, nil
}

func (s *PipelanerServer) SinkStream(stream service.Pipelaner_SinkStreamServer) error {
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
	s := &PipelanerServer{logger: logger, buffer: make(chan *service.Message, bufferSize)}
	return s
}

func (s *PipelanerServer) Recv() <-chan *service.Message {
	return s.buffer
}
