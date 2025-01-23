package tests

import (
	"context"
	"testing"
	"time"

	"github.com/randomurban/hw-test/hw12_13_14_15_calendar/internal/api/event"
	"github.com/randomurban/hw-test/hw12_13_14_15_calendar/internal/model"
	"github.com/randomurban/hw-test/hw12_13_14_15_calendar/internal/pb"
	"github.com/randomurban/hw-test/hw12_13_14_15_calendar/internal/service/mocks"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestCreate(t *testing.T) {
	t.Parallel()

	type mock struct {
		args model.Event
		ret  model.EventID
		err  error
	}

	type args struct {
		ctx context.Context
		req *pb.CreateEventRequest
	}
	ctx := context.Background()
	tests := []struct {
		name string
		args args
		want *pb.CreateEventResponse
		mock mock
		err  error
	}{
		{
			name: "success case",
			args: args{
				ctx: ctx,
				req: &pb.CreateEventRequest{
					Event: &pb.Event{
						Id:    1,
						Title: "Test Title 1",
						Start: &timestamppb.Timestamp{
							Seconds: 1733011200,
						},
						End: &timestamppb.Timestamp{
							Seconds: 1733097600,
						},
						Owner:       1,
						Description: "Test Description 1",
						NoticeTime: &durationpb.Duration{
							Seconds: 600,
						},
					},
				},
			},
			mock: mock{
				args: model.Event{
					ID:          1,
					Title:       "Test Title 1",
					Start:       time.Unix(1733011200, 0).UTC(),
					End:         time.Unix(1733097600, 0).UTC(),
					Owner:       1,
					Description: "Test Description 1",
					NoticeTime:  time.Second * 600,
				},
				ret: 1,
				err: nil,
			},
			want: &pb.CreateEventResponse{
				Id: 1,
			},
			err: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			serviceMock := mocks.NewEvent(t)
			serviceMock.EXPECT().Create(tt.args.ctx, tt.mock.args).Return(tt.mock.ret, tt.mock.err)
			api := event.NewImplementation(serviceMock)

			newID, err := api.Create(ctx, tt.args.req)
			assert.Equal(t, tt.err, err)
			assert.Equal(t, tt.want, newID)
		})
	}
}
