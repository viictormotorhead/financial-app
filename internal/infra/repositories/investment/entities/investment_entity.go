package entities

import "time"

// InvestmentEntity is the persistence DTO exchanged with the application layer.
type InvestmentEntity struct {
	ID        uint
	Name      string
	Balance        float64
	InitialBalance float64
	Tags           []string
	CreatedAt time.Time
}
