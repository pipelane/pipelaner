//nolint:revive
package grpc_server

import (
	"fmt"
	"net"
	"os"
	"syscall"

	"github.com/rs/zerolog"
	"google.golang.org/grpc"
)

const (
	unixConnectionType  = "unix"
	http2ConnectionType = "http2"
)

type ServerConfig struct {
	Host           string
	Port           int
	ConnectionType string
	UnixSocketPath *string
	Opts           []grpc.ServerOption
}

type Server struct {
	grpcServer *grpc.Server

	config *ServerConfig
	logger *zerolog.Logger
}

func NewServer(config *ServerConfig, logger *zerolog.Logger) *Server {
	return &Server{
		config: config,
		logger: logger,
	}
}

func (s *Server) Serve(use ...func(grpc *grpc.Server)) error {
	s.grpcServer = grpc.NewServer(s.config.Opts...)

	lis := s.createListener()

	for _, u := range use {
		u(s.grpcServer)
	}

	return s.grpcServer.Serve(lis)
}

func (s *Server) createListener() net.Listener {
	if s.config.ConnectionType == unixConnectionType {
		if err := syscall.Unlink(*s.config.UnixSocketPath); err != nil && !os.IsNotExist(err) {
			s.logger.Fatal().Err(err).Msgf("Failed to unlink Unix socket %s", s.config.UnixSocketPath)
		}

		unixListener, err := net.ListenUnix("unix", &net.UnixAddr{Name: *s.config.UnixSocketPath, Net: unixConnectionType})
		if err != nil {
			s.logger.Fatal().Err(err).Msgf("Failed to listen on Unix socket %s", s.config.UnixSocketPath)
		}
		return unixListener
	}

	tcpListener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", s.config.Host, s.config.Port))
	if err != nil {
		s.logger.Fatal().Err(err).Msgf("Failed to listen on TCP %s:%d", s.config.Host, s.config.Port)
	}
	return tcpListener
}

func (s *Server) Stop() error {
	if s.grpcServer == nil {
		return fmt.Errorf("grpc server not initialized")
	}
	s.grpcServer.GracefulStop()
	return nil
}
