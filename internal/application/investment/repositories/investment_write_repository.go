package repositories

import (
	"context"

	"github.com/viictormotorhead/financial-app/internal/infra/repositories/investment/entities"
)

type InvestmentWriteRepositoryIF interface {
	Save(ctx context.Context, investment entities.InvestmentEntity) (entities.InvestmentEntity, error)
	FindByID(ctx context.Context, id uint) (entities.InvestmentEntity, error)
	RecordMovement(ctx context.Context, input entities.RecordMovementInput) (entities.MovementResult, error)
	RecordValuation(ctx context.Context, input entities.RecordValuationInput) (entities.ValuationResult, error)
}
