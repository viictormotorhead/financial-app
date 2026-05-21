package response

import "time"

type InvestmentResponse struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	Balance   float64   `json:"balance"`
	Tags      []string  `json:"tags"`
	CreatedAt time.Time `json:"created_at"`
}
