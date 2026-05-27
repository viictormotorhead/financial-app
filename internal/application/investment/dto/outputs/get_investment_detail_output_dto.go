package outputs

import "time"

type GetInvestmentDetailOutputDTO struct {
	ID             uint
	Name           string
	Tags           []string
	CurrentValue   float64
	GrowthPercent  float64
	GrowthAmount   float64
	Series         []InvestmentSeriesPointDTO
	Movements      []InvestmentMovementDTO
}

type InvestmentSeriesPointDTO struct {
	Date  time.Time
	Value float64
}

type InvestmentMovementDTO struct {
	ID           uint
	Date         time.Time
	Type         string
	Amount       *float64
	BalanceAfter float64
}
