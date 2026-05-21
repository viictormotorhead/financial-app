package gormrepo

import (
	"context"
	"fmt"

	"gorm.io/gorm"

	"github.com/viictormotorhead/financial-app/internal/application/investment/repositories"
	"github.com/viictormotorhead/financial-app/internal/infra/repositories/investment/entities"
)

type investmentWriteRepository struct {
	db *gorm.DB
}

func NewInvestmentWriteRepository(db *gorm.DB) repositories.InvestmentWriteRepositoryIF {
	return &investmentWriteRepository{db: db}
}

func (r *investmentWriteRepository) Save(ctx context.Context, investment entities.InvestmentEntity) (entities.InvestmentEntity, error) {
	var saved entities.InvestmentEntity

	tags := investment.Tags
	if tags == nil {
		tags = []string{}
	}

	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		model := entities.Investment{
			Name:           investment.Name,
			Balance:        investment.Balance,
			InitialBalance: investment.InitialBalance,
			Tags:           entities.StringList(tags),
			CreatedAt:      investment.CreatedAt,
		}

		if err := tx.Create(&model).Error; err != nil {
			return fmt.Errorf("insert investment: %w", err)
		}

		history := entities.InvestmentHistory{
			InvestmentID: model.ID,
			Date:         investment.CreatedAt,
			Amount:       investment.Balance,
			BalanceAfter: investment.Balance,
			MovementType: entities.MovementTypeDeposit,
		}

		if err := tx.Create(&history).Error; err != nil {
			return fmt.Errorf("insert investment history: %w", err)
		}

		saved = entities.InvestmentEntity{
			ID:             model.ID,
			Name:           model.Name,
			Balance:        model.Balance,
			InitialBalance: model.InitialBalance,
			Tags:           model.Tags.Strings(),
			CreatedAt:      model.CreatedAt,
		}

		return nil
	})
	if err != nil {
		return entities.InvestmentEntity{}, err
	}

	return saved, nil
}
