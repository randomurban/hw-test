package sqlstorage

import (
	"context"
	"time"

	"github.com/randomurban/hw-test/hw12_13_14_15_calendar/internal/model"
)

type Storage struct { // TODO
}

func New() *Storage {
	return &Storage{}
}

func (s *Storage) Connect(ctx context.Context) error {
	// TODO
	return nil
}

func (s *Storage) Close(ctx context.Context) error {
	// TODO
	return nil
}

func (s *Storage) Create(ctx context.Context, event model.Event) error {
	return nil
}

func (s *Storage) Update(ctx context.Context, id string, event model.Event) error {
	// TODO
	return nil
}

func (s *Storage) Delete(ctx context.Context, id string) error {
	// TODO
	return nil
}

func (s *Storage) GetByID(ctx context.Context, id string) error {
	// TODO
	return nil
}

func (s *Storage) GetDay(ctx context.Context, date time.Time) error {
	// TODO
	return nil
}

func (s *Storage) GetWeek(ctx context.Context, date time.Time) error {
	// TODO
	return nil
}

func (s *Storage) GetMonth(ctx context.Context, date time.Time) error {
	// TODO
	return nil
}

// TODO
