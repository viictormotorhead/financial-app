package outputs

import "time"

type CreateInvestmentOutputDTO struct {
	ID        uint
	Name      string
	Balance        float64
	InitialBalance float64
	Tags           []string
	CreatedAt time.Time
}
