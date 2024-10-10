package model

import "time"

type NoticeID string

type Notice struct {
	ID    NoticeID
	Title string
	Start time.Time
	User  UserID // User ID
}
