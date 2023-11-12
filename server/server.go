/*
 * Copyright (c) 2023 Alexey Khokhlov
 */

package server

import (
	"io"
	"sync"

	"github.com/rs/zerolog"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/pipelane/pipelaner/server/service"
)

type PipelanerServer struct {
	service.UnimplementedPipelanerServer
	mu     sync.Mutex // protects routeNotes
	logger zerolog.Logger
	buffer chan *service.Message
}

func (s *PipelanerServer) Sink(stream service.Pipelaner_SinkServer) error {
	for {
		message, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&emptypb.Empty{})
		}
		if err != nil {
			return err
		}
		s.buffer <- message
	}
}

func NewServer(logger zerolog.Logger, bufferSize int64) *PipelanerServer {
	s := &PipelanerServer{logger: logger, buffer: make(chan *service.Message, bufferSize)}
	return s
}

func (s *PipelanerServer) Recv() <-chan *service.Message {
	return s.buffer
}
