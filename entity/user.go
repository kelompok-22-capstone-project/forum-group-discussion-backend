package entity

import "time"

type User struct {
	ID        string
	Username  string
	Email     string
	Name      string
	Password  string
	Role      string
	IsActive  bool
	CreatedAt time.Time
	UpdatedAt time.Time
}
