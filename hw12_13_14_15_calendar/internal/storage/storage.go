package storage

import (
	"context"
	"time"

	"github.com/randomurban/hw-test/hw12_13_14_15_calendar/internal/model"
)

type EventStorage interface {
	Create(ctx context.Context, event model.Event) (model.EventID, error)
	Update(ctx context.Context, id model.EventID, event model.Event) (bool, error)
	Delete(ctx context.Context, id model.EventID) (bool, error)
	GetByID(ctx context.Context, id model.EventID) (*model.Event, error)
	GetDayFromTo(ctx context.Context, from time.Time, to time.Time) (*[]model.Event, error)
}
