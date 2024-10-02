package memorystorage

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/randomurban/hw-test/hw12_13_14_15_calendar/internal/model"
)

func TestStorage(t *testing.T) {
	t.Run("memory store", func(t *testing.T) {
	})
}

func TestStartEnd(t *testing.T) {
	type args struct {
		start *time.Time
		end   *time.Time
	}

	t1 := time.Date(2010, 10, 20, 15, 10, 0, 0, time.UTC)
	t2 := time.Date(2020, 10, 20, 15, 10, 0, 0, time.UTC)

	tests := []struct {
		name string
		args args
		want time.Time
	}{
		{
			name: "max min",
			args: args{
				&t2, &t1,
			},
			want: t1,
		},
		{
			name: "min max",
			args: args{
				&t1, &t2,
			},
			want: t1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			StartEnd(tt.args.start, tt.args.end)
			// fmt.Printf("%v, %v - %v", tt.want, tt.args.start, tt.args.end)
			require.True(t, tt.args.start.Equal(tt.want))
		})
	}
}

func TestStorage_Create(t *testing.T) {
	t1 := time.Date(2010, 10, 20, 15, 10, 0, 0, time.UTC)
	t2 := time.Date(2020, 10, 20, 15, 10, 0, 0, time.UTC)
	event1 := model.Event{
		ID:          "event1",
		Title:       "test",
		Start:       t1,
		End:         t2,
		Owner:       "user1",
		Description: "test event1",
		NoticeTime:  time.Hour,
	}
	tests := []struct {
		name    string
		event   model.Event
		want    model.EventID
		wantErr bool
	}{
		{name: "simple", event: event1, want: "event1", wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := New()
			got, err := s.Create(context.Background(), tt.event)
			if (err != nil) != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Create() got = %v, want %v", got, tt.want)
			}
			require.Equal(t, 1, len(s.store))
		})
	}
}

func TestStorage_Delete(t *testing.T) {
	t1 := time.Date(2010, 10, 20, 15, 10, 0, 0, time.UTC)
	t2 := time.Date(2020, 10, 20, 15, 10, 0, 0, time.UTC)
	event1 := model.Event{
		ID:          "",
		Title:       "test",
		Start:       t1,
		End:         t2,
		Owner:       "user1",
		Description: "test event1",
		NoticeTime:  time.Hour,
	}
	tests := []struct {
		name    string
		event   model.Event
		want    bool
		wantErr bool
	}{
		{name: "simple", event: event1, want: true, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := New()
			gotEvent, err := s.Create(context.Background(), tt.event)
			if (err != nil) != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			require.Equal(t, 1, len(s.store))
			var got bool
			got, err = s.Delete(context.Background(), gotEvent)
			if (err != nil) != tt.wantErr {
				t.Errorf("Delete() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Delete() got = %v, want %v", got, tt.want)
			}
			require.Equal(t, 0, len(s.store))
		})
	}
}

func TestStorage_GetByID(t *testing.T) {
	t1 := time.Date(2010, 10, 20, 15, 10, 0, 0, time.UTC)
	t2 := time.Date(2020, 10, 20, 15, 10, 0, 0, time.UTC)
	event1 := model.Event{
		ID:          "event1",
		Title:       "test",
		Start:       t1,
		End:         t2,
		Owner:       "user1",
		Description: "test event1",
		NoticeTime:  time.Hour,
	}
	tests := []struct {
		name    string
		event   model.Event
		want    model.EventID
		wantErr bool
	}{
		{name: "simple", event: event1, want: "event1", wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := New()
			gotEvent, err := s.Create(context.Background(), tt.event)
			if (err != nil) != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			require.Equal(t, 1, len(s.store))
			var got *model.Event
			got, err = s.GetByID(context.Background(), gotEvent)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(*got, tt.event) {
				// if *got != tt.event {
				t.Errorf("GetByID() got = %v, want %v", *got, tt.event)
			}
			require.Equal(t, 1, len(s.store))
		})
	}
}

func TestStorage_GetDayFromTo(t *testing.T) {
	t1 := time.Date(2010, 10, 20, 15, 10, 0, 0, time.UTC)
	t2 := time.Date(2020, 10, 20, 15, 10, 0, 0, time.UTC)
	event1 := model.Event{
		ID:          "",
		Title:       "test",
		Start:       t1,
		End:         t2,
		Owner:       "user1",
		Description: "test event1",
		NoticeTime:  time.Hour,
	}
	tests := []struct {
		name    string
		event   model.Event
		from    time.Time
		to      time.Time
		want    int
		wantErr bool
	}{
		{name: "simple", event: event1, from: t1, to: t2, want: 1, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := New()
			_, err := s.Create(context.Background(), tt.event)
			if (err != nil) != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			require.Equal(t, 1, len(s.store))
			var got *[]model.Event
			got, err = s.GetDayFromTo(context.Background(), tt.from, tt.to)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetDayFromTo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			require.Equal(t, tt.want, len(*got))
			require.Equal(t, 1, len(s.store))
		})
	}
}

func TestStorage_Update(t *testing.T) {
	t1 := time.Date(2010, 10, 20, 15, 10, 0, 0, time.UTC)
	t2 := time.Date(2020, 10, 20, 15, 10, 0, 0, time.UTC)
	event1 := model.Event{
		ID:          "",
		Title:       "test",
		Start:       t1,
		End:         t2,
		Owner:       "user1",
		Description: "test event1",
		NoticeTime:  time.Hour,
	}
	tests := []struct {
		name    string
		event   model.Event
		eventID model.EventID
		want    bool
		wantErr bool
	}{
		{name: "simple", event: event1, want: true, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := New()
			gotEventID, err := s.Create(context.Background(), tt.event)
			if (err != nil) != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			require.Equal(t, 1, len(s.store))

			var event2 *model.Event
			event2, err = s.GetByID(context.Background(), gotEventID)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			event2.Title = "changed"

			var got bool
			got, err = s.Update(context.Background(), gotEventID, *event2)
			if (err != nil) != tt.wantErr {
				t.Errorf("Update() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Update() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(*event2, s.store[gotEventID]) {
				// if got != tt.want {
				t.Errorf("Update() event2 = %v, s.store[gotEventID] %v", event2, s.store[gotEventID])
			}
			require.Equal(t, 1, len(s.store))
		})
	}
}

func TestTimeIsBetween(t *testing.T) {
	t1 := time.Date(2010, 10, 20, 15, 10, 0, 0, time.UTC)
	t2 := time.Date(2015, 10, 20, 15, 10, 0, 0, time.UTC)
	t3 := time.Date(2020, 10, 20, 15, 10, 0, 0, time.UTC)
	t4 := time.Date(2025, 10, 20, 15, 10, 0, 0, time.UTC)
	type args struct {
		t     time.Time
		start time.Time
		end   time.Time
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "simple",
			args: args{
				t:     t1,
				start: t2,
				end:   t3,
			},
			want: false,
		},
		{
			name: "simple2",
			args: args{
				t:     t2,
				start: t4,
				end:   t1,
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := TimeIsBetween(tt.args.t, tt.args.start, tt.args.end); got != tt.want {
				t.Errorf("TimeIsBetween() = %v, want %v", got, tt.want)
			}
		})
	}
}
