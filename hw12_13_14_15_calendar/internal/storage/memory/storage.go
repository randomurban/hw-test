package memorystorage

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/randomurban/hw-test/hw12_13_14_15_calendar/internal/storage"

	"github.com/randomurban/hw-test/hw12_13_14_15_calendar/internal/model"
)

type Storage struct {
	store map[model.EventID]model.Event
	mu    sync.RWMutex
	id    model.EventID
}

func New() *Storage {
	store := make(map[model.EventID]model.Event)
	return &Storage{
		store: store,
		mu:    sync.RWMutex{},
		id:    0,
	}
}

func (s *Storage) newID() model.EventID {
	s.id++
	return s.id
}

var _ storage.EventStorage = (*Storage)(nil)

func (s *Storage) Connect(_ context.Context, _ string) error {
	return nil
}

func (s *Storage) Close(_ context.Context) {
}

func (s *Storage) Create(_ context.Context, event model.Event) (model.EventID, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if event.ID == 0 {
		event.ID = s.newID()
	}
	s.store[event.ID] = event

	return event.ID, nil
}

func (s *Storage) Update(_ context.Context, id model.EventID, event model.Event) (bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, ok := s.store[id]
	if !ok {
		return false, errors.New("event not found")
	}
	s.store[event.ID] = event
	return true, nil
}

func (s *Storage) Delete(_ context.Context, id model.EventID) (bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, ok := s.store[id]
	if !ok {
		return false, errors.New("event not found")
	}
	delete(s.store, id)
	return true, nil
}

func (s *Storage) GetByID(_ context.Context, id model.EventID) (*model.Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	event, ok := s.store[id]
	if !ok {
		return nil, errors.New("event not found")
	}
	return &event, nil
}

func (s *Storage) GetDayFromTo(_ context.Context, from time.Time, to time.Time) (*[]model.Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result := make([]model.Event, 0)
	for _, event := range s.store {
		es := event.Start
		ee := event.End
		StartEnd(&es, &ee)
		if TimeIsBetween(es, from, to) || TimeIsBetween(ee, from, to) ||
			(TimeIsBetween(from, es, ee) && TimeIsBetween(to, es, ee)) {
			result = append(result, event)
		}
	}
	return &result, nil
}

func TimeIsBetween(t, start, end time.Time) bool {
	StartEnd(&start, &end)
	return (t.Equal(start) || t.After(start)) && (t.Equal(end) || t.Before(end))
}

func StartEnd(start *time.Time, end *time.Time) {
	if start.After(*end) {
		*start, *end = *end, *start
	}
}
