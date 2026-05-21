package request

type CreateInvestmentRequest struct {
	Name    string  `json:"name"`
	Balance float64 `json:"balance"`
}
