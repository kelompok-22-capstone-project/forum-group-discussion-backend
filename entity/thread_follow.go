package entity

import "time"

type ThreadFollow struct {
	ID        string
	User      User
	Thread    Thread
	CreatedAt time.Time
	UpdatedAt time.Time
}
