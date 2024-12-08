package grpcserver

import (
	"context"
	"net"
	"time"

	"github.com/randomurban/hw-test/hw12_13_14_15_calendar/internal/config"
	"github.com/randomurban/hw-test/hw12_13_14_15_calendar/internal/logger"
	"github.com/randomurban/hw-test/hw12_13_14_15_calendar/internal/model"
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

	s.server = grpc.NewServer(grpc.UnaryInterceptor(s.RequestLogInterceptor))
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

func (s *Server) Create(ctx context.Context, req *pb.CreateEventRequest) (*pb.CreateEventResponse, error) {
	id, err := s.service.Create(ctx, *EventFromPB(req.Event))
	if err != nil {
		return nil, err
	}

	return &pb.CreateEventResponse{Id: int32(id)}, nil
}

func (s *Server) Update(ctx context.Context, req *pb.UpdateEventRequest) (*pb.UpdateEventResponse, error) {
	ok, err := s.service.Update(ctx, model.EventID(req.Event.GetId()), *EventFromPB(req.Event))
	if err != nil {
		return nil, err
	}

	return &pb.UpdateEventResponse{Ok: ok}, nil
}

func (s *Server) Delete(ctx context.Context, req *pb.DeleteEventRequest) (*pb.DeleteEventResponse, error) {
	ok, err := s.service.Delete(ctx, model.EventID(req.Id))
	if err != nil {
		return nil, err
	}
	return &pb.DeleteEventResponse{Ok: ok}, nil
}

func (s *Server) GetByID(ctx context.Context, req *pb.GetByIDEventRequest) (*pb.GetByIDEventResponse, error) {
	event, err := s.service.GetByID(ctx, model.EventID(req.GetId()))
	if err != nil {
		return nil, err
	}

	return &pb.GetByIDEventResponse{Event: EventToPB(event)}, nil
}

func (s *Server) GetDay(ctx context.Context, req *pb.GetDayEventRequest) (*pb.GetDayEventResponse, error) {
	events, err := s.service.GetDay(ctx, req.Start.AsTime())
	if err != nil {
		return nil, err
	}
	return &pb.GetDayEventResponse{Events: EventsToPB(events)}, nil
}

func (s *Server) GetWeek(ctx context.Context, req *pb.GetWeekEventRequest) (*pb.GetWeekEventResponse, error) {
	events, err := s.service.GetDay(ctx, req.Start.AsTime())
	if err != nil {
		return nil, err
	}
	return &pb.GetWeekEventResponse{Events: EventsToPB(events)}, nil
}

func (s *Server) GetMonth(ctx context.Context, req *pb.GetMonthEventRequest) (*pb.GetMonthEventResponse, error) {
	events, err := s.service.GetDay(ctx, req.Start.AsTime())
	if err != nil {
		return nil, err
	}
	return &pb.GetMonthEventResponse{Events: EventsToPB(events)}, nil
}
