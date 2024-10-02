package memorystorage

import (
	"context"
	"reflect"
	"sync"
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
			want: t2,
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
		genID   func() model.EventID
		event   model.Event
		want    model.EventID
		wantErr bool
	}{
		{name: "simple", event: event1, genID: model.NewEventID, want: "event1", wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := New(tt.genID)
			got, err := s.Create(context.Background(), tt.event)
			if (err != nil) != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Create() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStorage_Delete(t *testing.T) {
	type fields struct {
		store map[model.EventID]model.Event
		mu    sync.RWMutex
	}
	type args struct {
		ctx context.Context
		id  model.EventID
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    bool
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Storage{
				store: tt.fields.store,
				mu:    tt.fields.mu,
			}
			got, err := s.Delete(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("Delete() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Delete() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStorage_GetByID(t *testing.T) {
	type fields struct {
		store map[model.EventID]model.Event
		mu    sync.RWMutex
	}
	type args struct {
		ctx context.Context
		id  model.EventID
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *model.Event
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Storage{
				store: tt.fields.store,
				mu:    tt.fields.mu,
			}
			got, err := s.GetByID(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetByID() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStorage_GetDayFromTo(t *testing.T) {
	type fields struct {
		store map[model.EventID]model.Event
		mu    sync.RWMutex
	}
	type args struct {
		ctx  context.Context
		from time.Time
		to   time.Time
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *[]model.Event
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Storage{
				store: tt.fields.store,
				mu:    tt.fields.mu,
			}
			got, err := s.GetDayFromTo(tt.args.ctx, tt.args.from, tt.args.to)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetDayFromTo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetDayFromTo() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStorage_Update(t *testing.T) {
	type fields struct {
		store map[model.EventID]model.Event
		mu    sync.RWMutex
	}
	type args struct {
		ctx   context.Context
		id    model.EventID
		event model.Event
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    bool
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Storage{
				store: tt.fields.store,
				mu:    tt.fields.mu,
			}
			got, err := s.Update(tt.args.ctx, tt.args.id, tt.args.event)
			if (err != nil) != tt.wantErr {
				t.Errorf("Update() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Update() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTimeIsBetween(t *testing.T) {
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
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := TimeIsBetween(tt.args.t, tt.args.start, tt.args.end); got != tt.want {
				t.Errorf("TimeIsBetween() = %v, want %v", got, tt.want)
			}
		})
	}
}
