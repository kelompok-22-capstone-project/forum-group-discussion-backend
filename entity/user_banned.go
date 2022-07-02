package entity

import "time"

type UserBanned struct {
	ID        string
	Moderator Moderator
	User      User
	Thread    Thread
	Comment   Comment
	Reason    string
	Status    string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type ReportStatus int

const (
	Review ReportStatus = iota
	Accepted
)
