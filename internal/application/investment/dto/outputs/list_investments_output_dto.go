package outputs

type InvestmentAllocationDTO struct {
	ID                 uint
	Investment         string
	Amount             float64
	Percentage         float64
	PercentageGrowing  float64
	Tags               []string
}

type ListInvestmentsOutputDTO struct {
	Items []InvestmentAllocationDTO
}
