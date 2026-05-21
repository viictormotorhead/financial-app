package entities

import "time"

type InvestmentHistory struct {
	ID             uint      `gorm:"primaryKey"`
	InvestmentID   uint      `gorm:"not null;index"`
	Date           time.Time `gorm:"not null"`
	Amount         float64   `gorm:"type:double precision;not null"`
	MovementType   string    `gorm:"type:varchar(50);not null;column:movement_type"`
}

func (InvestmentHistory) TableName() string {
	return "investment_history"
}
