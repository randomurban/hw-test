package grpc_server

import (
	"context"
	"net"
	"time"

	"github.com/randomurban/hw-test/hw12_13_14_15_calendar/internal/config"
	"github.com/randomurban/hw-test/hw12_13_14_15_calendar/internal/logger"
	"github.com/randomurban/hw-test/hw12_13_14_15_calendar/internal/pb"
	"github.com/randomurban/hw-test/hw12_13_14_15_calendar/internal/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type Server struct {
	pb.UnimplementedEventServiceServer
	logger  logger.Logger
	service service.Event
	server  *grpc.Server
	addr    string
}

func NewServer(cfg config.Config, logger logger.Logger, service service.Event) *Server {
	return &Server{
		logger:  logger,
		service: service,
		addr:    cfg.GRPC.Addr,
	}
}

func (s *Server) Start(ctx context.Context) error {
	go func() {
		<-ctx.Done()

		s.logger.Info("calendar is stopping...")
		ctxStop, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		if err := s.Stop(ctxStop); err != nil {
			s.logger.Error("failed to stop grpc: " + err.Error())
		}
	}()

	s.server = grpc.NewServer()
	reflection.Register(s.server)
	pb.RegisterEventServiceServer(s.server, s)
	lsn, err := net.Listen("tcp", s.addr)
	if err != nil {
		return err
	}
	s.logger.Info("GRPC server started ", "addr", s.addr)
	return s.server.Serve(lsn)
}

func (s *Server) Stop(_ context.Context) error {
	s.server.GracefulStop()
	return nil
}
