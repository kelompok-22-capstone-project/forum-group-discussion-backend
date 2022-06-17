package entity

import "time"

type Like struct {
	ID        string
	User      User
	Thread    Thread
	CreatedAt time.Time
	UpdatedAt time.Time
}
