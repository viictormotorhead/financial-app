package outputs

import "time"

type CreateInvestmentOutputDTO struct {
	ID        uint
	Name      string
	Balance   float64
	CreatedAt time.Time
}
