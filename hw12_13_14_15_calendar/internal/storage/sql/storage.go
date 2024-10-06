package sqlstorage

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/randomurban/hw-test/hw12_13_14_15_calendar/internal/storage"

	"github.com/randomurban/hw-test/hw12_13_14_15_calendar/internal/model"
)

type Storage struct {
	pool *pgxpool.Pool
}

var _ storage.EventStorage = (*Storage)(nil)

func New() *Storage {
	return &Storage{}
}

func (s *Storage) Connect(ctx context.Context, dsn string) error {
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return fmt.Errorf("connecting to database: %w", err)
	}
	s.pool = pool
	return pool.Ping(ctx)
}

func (s *Storage) Close(_ context.Context) {
	s.pool.Close()
}

func (s *Storage) Create(ctx context.Context, event model.Event) (model.EventID, error) {
	query := `insert into event (title, start_time, end_time, owner, description, notice_time) 
		VALUES ($1, $2, $3, $4, $5, $6) returning id`
	err := s.pool.QueryRow(ctx, query, event.Title, event.Start.Local(), event.End.Local(),
		event.Owner, event.Description, event.Start.Add(-event.NoticeTime)).Scan(&event.ID)
	if err != nil {
		return 0, fmt.Errorf("creating event: %w", err)
	}
	return event.ID, nil
}

func (s *Storage) Update(ctx context.Context, id model.EventID, event model.Event) (bool, error) {
	query := `update event set title=$1, start_time=$2, end_time=$3, owner=$4, description=$5, notice_time=$6 
             where id=$7`
	_, err := s.pool.Exec(ctx, query, event.Title, event.Start, event.End,
		event.Owner, event.Description, event.Start.Add(-event.NoticeTime), id)
	if err != nil {
		return false, fmt.Errorf("updating event: %w", err)
	}
	return true, nil
}

func (s *Storage) Delete(ctx context.Context, id model.EventID) (bool, error) {
	query := `delete from event where id=$1`
	_, err := s.pool.Exec(ctx, query, id)
	if err != nil {
		return false, fmt.Errorf("deleting event: %w", err)
	}
	return true, nil
}

func (s *Storage) GetByID(ctx context.Context, id model.EventID) (*model.Event, error) {
	var event model.Event
	var timestamp time.Time
	query := `select id, title, start_time, end_time, owner, description, notice_time from event where id=$1`
	err := s.pool.QueryRow(ctx, query, id).Scan(&event.ID, &event.Title, &event.Start, &event.End,
		&event.Owner, &event.Description, &timestamp)
	if err != nil {
		return nil, fmt.Errorf("getting event: %w", err)
	}
	event.NoticeTime = timestamp.Sub(event.Start)
	return &event, nil
}

func (s *Storage) GetDayFromTo(ctx context.Context, from time.Time, to time.Time) (*[]model.Event, error) {
	query := `select id,title,start_time,end_time,owner,description,notice_time from event 
			where (start_time between $1 and $2) and (end_time between $1 and $2)`
	rows, err := s.pool.Query(ctx, query, from, to)
	if err != nil {
		return nil, fmt.Errorf("getting events: %w", err)
	}
	defer rows.Close()
	var res []model.Event
	for rows.Next() {
		var timestamp time.Time
		var event model.Event
		err = rows.Scan(&event.ID, &event.Title, &event.Start, &event.End,
			&event.Owner, &event.Description, &timestamp)
		if err != nil {
			return nil, fmt.Errorf("getting events: %w", err)
		}
		event.NoticeTime = -timestamp.Sub(event.Start)
		res = append(res, event)
	}
	return &res, nil
}
