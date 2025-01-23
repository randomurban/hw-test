package converter

import (
	"github.com/randomurban/hw-test/hw12_13_14_15_calendar/internal/model"
	"github.com/randomurban/hw-test/hw12_13_14_15_calendar/internal/pb"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func EventFromPB(event *pb.Event) *model.Event {
	return &model.Event{
		ID:          model.EventID(event.GetId()),
		Title:       event.GetTitle(),
		Start:       event.GetStart().AsTime(),
		End:         event.GetEnd().AsTime(),
		Owner:       model.UserID(event.GetOwner()),
		Description: event.GetDescription(),
		NoticeTime:  event.GetNoticeTime().AsDuration(),
	}
}

func EventToPB(event *model.Event) *pb.Event {
	return &pb.Event{
		Id:          int32(event.ID),
		Title:       event.Title,
		Start:       timestamppb.New(event.Start),
		End:         timestamppb.New(event.End),
		Owner:       int32(event.Owner),
		Description: event.Description,
		NoticeTime:  durationpb.New(event.NoticeTime),
	}
}

func EventsToPB(events *[]model.Event) []*pb.Event {
	pbEvents := make([]*pb.Event, len(*events))
	for i, event := range *events {
		pbEvents[i] = EventToPB(&event)
	}
	return pbEvents
}
