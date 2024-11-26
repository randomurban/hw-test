package service

import (
	"context"
	"time"

	"github.com/randomurban/hw-test/hw12_13_14_15_calendar/internal/model"
)

type Event interface {
	Create(ctx context.Context, event model.Event) (model.EventID, error)
	Update(ctx context.Context, id model.EventID, event model.Event) (bool, error)
	Delete(ctx context.Context, id model.EventID) (bool, error)
	GetByID(ctx context.Context, id model.EventID) (*model.Event, error)
	GetDay(ctx context.Context, date time.Time) (*[]model.Event, error)
	GetWeek(ctx context.Context, date time.Time) (*[]model.Event, error)
	GetMonth(ctx context.Context, date time.Time) (*[]model.Event, error)
}
