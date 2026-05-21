package outputs

import (
	"time"

	"github.com/viictormotorhead/financial-app/internal/infra/repositories/investment/entities"
)

type CreateMovementOutputDTO struct {
	HistoryID    uint
	InvestmentID uint
	Type         entities.MovementType
	Amount       float64
	Balance      float64
	Date         time.Time
}
