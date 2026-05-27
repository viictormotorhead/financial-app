package entities

import "time"

type InvestmentHistoryEntity struct {
	ID           uint
	InvestmentID uint
	Date         time.Time
	Amount       float64
	BalanceAfter float64
	MovementType MovementType
}
