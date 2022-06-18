package entity

import "time"

type UserFollow struct {
	ID        string
	User      User
	Following User
	CreatedAt time.Time
	UpdatedAt time.Time
}
