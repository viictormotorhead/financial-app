package response

import (
	"time"

	"github.com/viictormotorhead/financial-app/internal/infra/repositories/investment/entities"
)

type MovementResponse struct {
	ID           uint                  `json:"id"`
	InvestmentID uint                  `json:"investment_id"`
	Type         entities.MovementType `json:"type"`
	Amount       float64               `json:"amount"`
	Balance      float64               `json:"balance"`
	Date         time.Time             `json:"date"`
}
