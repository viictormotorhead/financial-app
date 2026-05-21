package repositories

import (
	"context"

	"github.com/viictormotorhead/financial-app/internal/infra/repositories/investment/entities"
)

type InvestmentWriteRepositoryIF interface {
	Save(ctx context.Context, investment entities.InvestmentEntity) (entities.InvestmentEntity, error)
}
