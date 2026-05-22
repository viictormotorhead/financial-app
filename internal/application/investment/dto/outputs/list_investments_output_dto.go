package outputs

type InvestmentAllocationDTO struct {
	Investment string
	Amount     float64
	Percentage float64
	Tags       []string
}

type ListInvestmentsOutputDTO struct {
	Items []InvestmentAllocationDTO
}
