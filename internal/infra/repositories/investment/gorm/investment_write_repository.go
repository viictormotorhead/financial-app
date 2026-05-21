package gormrepo

import (
	"context"
	"fmt"
	"time"

	"gorm.io/gorm"

	"github.com/viictormotorhead/financial-app/internal/application/investment/repositories"
	"github.com/viictormotorhead/financial-app/internal/infra/repositories/investment/entities"
)

const movementTypeInitial = "initial"

type investmentWriteRepository struct {
	db *gorm.DB
}

func NewInvestmentWriteRepository(db *gorm.DB) repositories.InvestmentWriteRepositoryIF {
	return &investmentWriteRepository{db: db}
}

func (r *investmentWriteRepository) Save(ctx context.Context, investment entities.InvestmentEntity) (entities.InvestmentEntity, error) {
	var saved entities.InvestmentEntity

	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		model := entities.Investment{
			Name:      investment.Name,
			Balance:   investment.Balance,
			CreatedAt: investment.CreatedAt,
		}

		if err := tx.Create(&model).Error; err != nil {
			return fmt.Errorf("insert investment: %w", err)
		}

		history := entities.InvestmentHistory{
			InvestmentID: model.ID,
			Date:         time.Now().UTC(),
			Amount:       investment.Balance,
			MovementType: movementTypeInitial,
		}

		if err := tx.Create(&history).Error; err != nil {
			return fmt.Errorf("insert investment history: %w", err)
		}

		saved = entities.InvestmentEntity{
			ID:        model.ID,
			Name:      model.Name,
			Balance:   model.Balance,
			CreatedAt: model.CreatedAt,
		}

		return nil
	})
	if err != nil {
		return entities.InvestmentEntity{}, err
	}

	return saved, nil
}
