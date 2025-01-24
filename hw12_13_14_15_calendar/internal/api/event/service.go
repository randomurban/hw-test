package event

import (
	"context"

	"github.com/randomurban/hw-test/hw12_13_14_15_calendar/internal/converter"
	"github.com/randomurban/hw-test/hw12_13_14_15_calendar/internal/model"
	"github.com/randomurban/hw-test/hw12_13_14_15_calendar/internal/pb"
	"github.com/randomurban/hw-test/hw12_13_14_15_calendar/internal/service"
)

type Implementation struct {
	pb.UnimplementedEventServiceServer
	service service.Event
}

func NewImplementation(eventService service.Event) *Implementation {
	return &Implementation{
		service: eventService,
	}
}

func (i *Implementation) Create(ctx context.Context, req *pb.CreateEventRequest) (*pb.CreateEventResponse, error) {
	id, err := i.service.Create(ctx, *converter.EventFromPB(req.Event))
	if err != nil {
		return nil, err
	}

	return &pb.CreateEventResponse{Id: int32(id)}, nil
}

func (i *Implementation) Update(ctx context.Context, req *pb.UpdateEventRequest) (*pb.UpdateEventResponse, error) {
	ok, err := i.service.Update(ctx, model.EventID(req.Event.GetId()), *converter.EventFromPB(req.Event))
	if err != nil {
		return nil, err
	}

	return &pb.UpdateEventResponse{Ok: ok}, nil
}

func (i *Implementation) Delete(ctx context.Context, req *pb.DeleteEventRequest) (*pb.DeleteEventResponse, error) {
	ok, err := i.service.Delete(ctx, model.EventID(req.Id))
	if err != nil {
		return nil, err
	}
	return &pb.DeleteEventResponse{Ok: ok}, nil
}

func (i *Implementation) GetByID(ctx context.Context, req *pb.GetByIDEventRequest) (*pb.GetByIDEventResponse, error) {
	event, err := i.service.GetByID(ctx, model.EventID(req.GetId()))
	if err != nil {
		return nil, err
	}

	return &pb.GetByIDEventResponse{Event: converter.EventToPB(event)}, nil
}

func (i *Implementation) GetDay(ctx context.Context, req *pb.GetDayEventRequest) (*pb.GetDayEventResponse, error) {
	events, err := i.service.GetDay(ctx, req.Start.AsTime())
	if err != nil {
		return nil, err
	}
	return &pb.GetDayEventResponse{Events: converter.EventsToPB(events)}, nil
}

func (i *Implementation) GetWeek(ctx context.Context, req *pb.GetWeekEventRequest) (*pb.GetWeekEventResponse, error) {
	events, err := i.service.GetDay(ctx, req.Start.AsTime())
	if err != nil {
		return nil, err
	}
	return &pb.GetWeekEventResponse{Events: converter.EventsToPB(events)}, nil
}

func (i *Implementation) GetMonth(
	ctx context.Context, req *pb.GetMonthEventRequest,
) (*pb.GetMonthEventResponse, error) {
	events, err := i.service.GetDay(ctx, req.Start.AsTime())
	if err != nil {
		return nil, err
	}
	return &pb.GetMonthEventResponse{Events: converter.EventsToPB(events)}, nil
}
