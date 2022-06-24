package entity

import "time"

type User struct {
	ID            string
	Username      string
	Email         string
	Name          string
	Password      string
	Role          string
	IsActive      bool
	TotalThread   uint64
	TotalFollower uint64
	IsFollowed    bool
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type UserStatus bool

const (
	Banned UserStatus = false
	Active UserStatus = true
)

type UserOrderBy int

const (
	RegisteredDate UserOrderBy = iota
	Ranking
)
