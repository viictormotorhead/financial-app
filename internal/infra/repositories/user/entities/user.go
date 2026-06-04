package entities

import "time"

type User struct {
	UserID       string    `gorm:"column:user_id;type:varchar(21);primaryKey"`
	Name         string    `gorm:"type:varchar(255);not null"`
	Username     string    `gorm:"type:varchar(100);not null;uniqueIndex"`
	PasswordHash string    `gorm:"column:password_hash;type:varchar(255);not null"`
	CreatedAt    time.Time `gorm:"not null"`
}

func (User) TableName() string {
	return "users"
}
