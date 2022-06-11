package entity

import "time"

type Thread struct {
	ID          string
	Title       string
	Description string
	TotalViewer uint64
	Creator     User
	Category    Category
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
