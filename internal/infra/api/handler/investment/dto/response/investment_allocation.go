package response

type InvestmentAllocationResponse struct {
	ID                uint     `json:"id"`
	Investment        string   `json:"investment"`
	Amount            float64  `json:"amount"`
	Percentage        float64  `json:"percentage"`
	PercentageGrowing float64  `json:"percentage-growing"`
	Tags              []string `json:"tags"`
}
