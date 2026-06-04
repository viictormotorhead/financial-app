package entities

import "time"

type UserEntity struct {
	UserID       string
	Name         string
	Username     string
	PasswordHash string
	CreatedAt    time.Time
}
