package pprof

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/pprof"
	"time"

	pprofCfg "github.com/pipelane/pipelaner/gen/settings/pprof"
)

type Server struct {
	server *http.Server
	cfg    *pprofCfg.Config
}

func NewServer(cfg *pprofCfg.Config) *Server {
	r := http.NewServeMux()
	r.HandleFunc(cfg.Path, pprof.Index)
	r.HandleFunc(cfg.Path+"/cmdline", pprof.Cmdline)
	r.HandleFunc(cfg.Path+"/profile", pprof.Profile)
	r.HandleFunc(cfg.Path+"/symbol", pprof.Symbol)
	r.HandleFunc(cfg.Path+"/trace", pprof.Trace)
	server := &http.Server{
		Addr:              fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		ReadHeaderTimeout: 10 * time.Second,
		Handler:           r,
	}
	return &Server{
		server: server,
		cfg:    cfg,
	}
}

func (m *Server) Serve(_ context.Context) error {
	if err := m.server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}

func (m *Server) Shutdown(ctx context.Context) error {
	return m.server.Shutdown(ctx)
}
