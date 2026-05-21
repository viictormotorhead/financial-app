package response

import "time"

type ValuationResponse struct {
	ID              uint      `json:"id"`
	InvestmentID    uint      `json:"investment_id"`
	PreviousBalance float64   `json:"previous_balance"`
	CurrentValue    float64   `json:"current_value"`
	Delta           float64   `json:"delta"`
	Balance         float64   `json:"balance"`
	Date            time.Time `json:"date"`
}
