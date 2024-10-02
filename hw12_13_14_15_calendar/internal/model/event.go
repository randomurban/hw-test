package model

import (
	"time"

	"github.com/rs/xid"
)

type (
	EventID string
	UserId  string
)

type Event struct {
	ID          EventID
	Title       string
	Start       time.Time
	End         time.Time
	Owner       UserId        // User ID
	Description string        // opt
	NoticeTime  time.Duration // opt
}

func NewEventID() EventID {
	return EventID(xid.New().String())
}
