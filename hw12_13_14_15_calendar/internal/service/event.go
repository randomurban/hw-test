package service

import (
	"context"
	"time"

	"github.com/randomurban/hw-test/hw12_13_14_15_calendar/internal/model"
	"github.com/randomurban/hw-test/hw12_13_14_15_calendar/internal/storage"
)

type EventService struct {
	storage storage.EventStorage
}

type Logger interface {
	Debug(msg string)
	Info(msg string)
	Warn(msg string)
	Error(msg string)
}

func New(logger Logger, storage storage.EventStorage) *EventService {
	return &EventService{}
}

func (a *EventService) Create(ctx context.Context, event model.Event) (bool, error) {
	return a.storage.Create(ctx, event)
}

func (a *EventService) Update(ctx context.Context, id model.EventID, event model.Event) (bool, error) {
	return a.storage.Update(ctx, id, event)
}

func (a *EventService) Delete(ctx context.Context, id model.EventID) (bool, error) {
	return a.storage.Delete(ctx, id)
}

func (a *EventService) GetByID(ctx context.Context, id model.EventID) (*model.Event, error) {
	return a.storage.GetByID(ctx, id)
}

func (a *EventService) GetDay(ctx context.Context, date time.Time) (*[]model.Event, error) {
	return a.storage.GetDayFromTo(ctx, date, date.AddDate(0, 0, 1))
}

func (a *EventService) GetWeek(ctx context.Context, date time.Time) (*[]model.Event, error) {
	return a.storage.GetDayFromTo(ctx, date, date.AddDate(0, 0, 7))
}

func (a *EventService) GetMonth(ctx context.Context, date time.Time) (*[]model.Event, error) {
	return a.storage.GetDayFromTo(ctx, date, date.AddDate(0, 1, 0))
}
