package entity

import "time"

type UserBanned struct {
	ID        string
	Moderator Moderator
	User      User
	Reason    string
	Status    string
	CreatedAt time.Time
	UpdatedAt time.Time
}
