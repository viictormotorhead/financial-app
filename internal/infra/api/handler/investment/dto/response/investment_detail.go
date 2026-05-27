package response

type InvestmentDetailResponse struct {
	ID            uint                            `json:"id"`
	Name          string                          `json:"name"`
	Tags          []string                        `json:"tags"`
	CurrentValue  float64                         `json:"current_value"`
	GrowthPercent float64                         `json:"growth_percent"`
	GrowthAmount  float64                         `json:"growth_amount"`
	Series        []InvestmentSeriesPointResponse `json:"series"`
	Movements     []InvestmentDetailMovementResponse `json:"movements"`
}

type InvestmentSeriesPointResponse struct {
	Date  string  `json:"date"`
	Value float64 `json:"value"`
}

type InvestmentDetailMovementResponse struct {
	ID           uint     `json:"id"`
	Date         string   `json:"date"`
	Type         string   `json:"type"`
	Amount       *float64 `json:"amount"`
	BalanceAfter float64  `json:"balance_after"`
}
