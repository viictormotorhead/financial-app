package response

type InvestmentAllocationResponse struct {
	Investment string   `json:"investment"`
	Amount     float64  `json:"amount"`
	Percentage float64  `json:"percentage"`
	Tags       []string `json:"tags"`
}
