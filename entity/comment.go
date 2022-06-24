package entity

import "time"

type Comment struct {
	ID        string
	User      User
	Thread    Thread
	Comment   string
	CreatedAt time.Time
	UpdatedAt time.Time
}
