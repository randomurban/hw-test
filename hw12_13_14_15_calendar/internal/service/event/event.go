package event

import (
	"context"
	"time"

	"github.com/randomurban/hw-test/hw12_13_14_15_calendar/internal/logger"
	"github.com/randomurban/hw-test/hw12_13_14_15_calendar/internal/model"
	"github.com/randomurban/hw-test/hw12_13_14_15_calendar/internal/storage"
)

type Event struct {
	storage storage.EventStorage
}

func New(logger logger.Logger, storage storage.EventStorage) *Event {
	return &Event{}
}

func (a *Event) Create(ctx context.Context, event model.Event) (model.EventID, error) {
	return a.storage.Create(ctx, event)
}

func (a *Event) Update(ctx context.Context, id model.EventID, event model.Event) (bool, error) {
	return a.storage.Update(ctx, id, event)
}

func (a *Event) Delete(ctx context.Context, id model.EventID) (bool, error) {
	return a.storage.Delete(ctx, id)
}

func (a *Event) GetByID(ctx context.Context, id model.EventID) (*model.Event, error) {
	return a.storage.GetByID(ctx, id)
}

func (a *Event) GetDay(ctx context.Context, date time.Time) (*[]model.Event, error) {
	return a.storage.GetDayFromTo(ctx, date, date.AddDate(0, 0, 1))
}

func (a *Event) GetWeek(ctx context.Context, date time.Time) (*[]model.Event, error) {
	return a.storage.GetDayFromTo(ctx, date, date.AddDate(0, 0, 7))
}

func (a *Event) GetMonth(ctx context.Context, date time.Time) (*[]model.Event, error) {
	return a.storage.GetDayFromTo(ctx, date, date.AddDate(0, 1, 0))
}
