package outputs

import "time"

type CreateValuationOutputDTO struct {
	HistoryID       uint
	InvestmentID    uint
	PreviousBalance float64
	CurrentValue    float64
	Delta           float64
	Balance         float64
	Date            time.Time
}
