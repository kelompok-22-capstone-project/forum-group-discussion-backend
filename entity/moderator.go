package entity

import "time"

type Moderator struct {
	ID        string
	User      User
	ThreadID  string
	CreatedAt time.Time
	UpdatedAt time.Time
}
