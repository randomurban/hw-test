package internalhttp

import (
	"context"
	"fmt"
	"net/http"

	"github.com/randomurban/hw-test/hw12_13_14_15_calendar/internal/config"
	"github.com/randomurban/hw-test/hw12_13_14_15_calendar/internal/logger"
	"github.com/randomurban/hw-test/hw12_13_14_15_calendar/internal/service"
)

type Server struct {
	httpServer *http.Server
	logger     logger.Logger
	service    service.Event
}

func NewServer(cfg config.Config, logger logger.Logger, service service.Event) *Server {
	return &Server{
		logger: logger,
		httpServer: &http.Server{
			Addr:              cfg.HTTP.Addr,
			ReadHeaderTimeout: cfg.HTTP.ReadHeaderTimeout,
		},
		service: service,
	}
}

func (s *Server) Start(_ context.Context) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/", s.helloHandler)
	s.httpServer.Handler = s.loggingMiddleware(mux)
	return s.httpServer.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}

func (s *Server) helloHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Hello %s", r.RemoteAddr)
}
