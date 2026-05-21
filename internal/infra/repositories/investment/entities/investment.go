package entities

import "time"

type Investment struct {
	ID        uint      `gorm:"primaryKey"`
	Name      string    `gorm:"type:varchar(255);not null"`
	Balance   float64   `gorm:"type:double precision;not null;default:0"`
	CreatedAt time.Time `gorm:"not null"`

	History []InvestmentHistory `gorm:"foreignKey:InvestmentID"`
}

func (Investment) TableName() string {
	return "investment"
}
