package entity

import "time"

type Thread struct {
	ID            string
	Title         string
	Description   string
	TotalViewer   uint64
	TotalLike     uint64
	TotalFollower uint64
	TotalComment  uint64
	Creator       User
	Category      Category
	IsLiked       bool
	IsFollowed    bool
	Moderators    []Moderator
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
