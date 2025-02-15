package httpserver

import (
	"context"
	"net/http"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/randomurban/hw-test/hw12_13_14_15_calendar/internal/config"
	"github.com/randomurban/hw-test/hw12_13_14_15_calendar/internal/logger"
	"github.com/randomurban/hw-test/hw12_13_14_15_calendar/internal/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Server struct {
	httpServer *http.Server
	logger     logger.Logger
	service    pb.EventServiceServer
	cfg        config.Config
}

func NewServer(cfg config.Config, logger logger.Logger, service pb.EventServiceServer) *Server {
	return &Server{
		logger: logger,
		httpServer: &http.Server{
			Addr:              cfg.HTTP.Addr,
			ReadHeaderTimeout: cfg.HTTP.ReadHeaderTimeout,
		},
		service: service,
		cfg:     cfg,
	}
}

func (s *Server) Start(ctx context.Context) error {
	go func() {
		<-ctx.Done()
		s.logger.Info("calendar is stopping...")
		ctxStop, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		if err := s.Stop(ctxStop); err != nil {
			s.logger.Error("failed to stop http: " + err.Error())
		}
	}()

	mux := runtime.NewServeMux()

	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}
	err := pb.RegisterEventServiceHandlerFromEndpoint(ctx, mux, s.cfg.GRPC.Addr, opts)
	if err != nil {
		return err
	}

	s.httpServer.Handler = s.loggingMiddleware(mux)
	s.logger.Info("HTTP server started ", "addr", s.httpServer.Addr)

	return s.httpServer.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
