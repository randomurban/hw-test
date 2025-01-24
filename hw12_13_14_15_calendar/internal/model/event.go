package model

import (
	"time"
)

type (
	EventID int32
	UserID  int32
)

type Event struct {
	ID          EventID
	Title       string
	Start       time.Time
	End         time.Time
	Owner       UserID        // User ID
	Description string        // opt
	NoticeTime  time.Duration // opt
}

// Unique ID if EventID string
// func NewEventID() EventID {
// 	return EventID(xid.New().String())
// }
