package gormrepo

import (
	"context"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	appinvestment "github.com/viictormotorhead/financial-app/internal/application/investment"
	"github.com/viictormotorhead/financial-app/internal/infra/repositories/investment/entities"
	"github.com/viictormotorhead/financial-app/internal/infra/repositories/scopes"
)

func (r *investmentWriteRepository) FindByID(ctx context.Context, userID string, id uint) (entities.InvestmentEntity, error) {
	var model entities.Investment
	if err := scopes.UserID(r.db.WithContext(ctx), userID).First(&model, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entities.InvestmentEntity{}, appinvestment.ErrInvestmentNotFound
		}
		return entities.InvestmentEntity{}, fmt.Errorf("find investment: %w", err)
	}

	return toInvestmentEntity(model), nil
}

func (r *investmentWriteRepository) RecordMovement(ctx context.Context, input entities.RecordMovementInput) (entities.MovementResult, error) {
	var result entities.MovementResult
	now := time.Now().UTC()

	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var model entities.Investment
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&model, input.InvestmentID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return appinvestment.ErrInvestmentNotFound
			}
			return fmt.Errorf("find investment: %w", err)
		}

		if err := tx.Model(&model).Update("balance", input.NewBalance).Error; err != nil {
			return fmt.Errorf("update investment balance: %w", err)
		}

		history := entities.InvestmentHistory{
			InvestmentID: model.ID,
			Date:         now,
			Amount:       input.Amount,
			BalanceAfter: input.NewBalance,
			MovementType: input.MovementType,
		}

		if err := tx.Create(&history).Error; err != nil {
			return fmt.Errorf("insert investment history: %w", err)
		}

		result = entities.MovementResult{
			HistoryID:    history.ID,
			InvestmentID: model.ID,
			MovementType: input.MovementType,
			Amount:       input.Amount,
			Balance:      input.NewBalance,
			Date:         now,
		}

		return nil
	})
	if err != nil {
		return entities.MovementResult{}, err
	}

	return result, nil
}

func (r *investmentWriteRepository) RecordValuation(ctx context.Context, input entities.RecordValuationInput) (entities.ValuationResult, error) {
	var result entities.ValuationResult
	now := time.Now().UTC()

	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var model entities.Investment
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&model, input.InvestmentID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return appinvestment.ErrInvestmentNotFound
			}
			return fmt.Errorf("find investment: %w", err)
		}

		previousBalance := model.Balance

		if err := tx.Model(&model).Update("balance", input.NewBalance).Error; err != nil {
			return fmt.Errorf("update investment balance: %w", err)
		}

		history := entities.InvestmentHistory{
			InvestmentID: model.ID,
			Date:         now,
			Amount:       input.Delta,
			BalanceAfter: input.NewBalance,
			MovementType: entities.MovementTypeEarning,
		}

		if err := tx.Create(&history).Error; err != nil {
			return fmt.Errorf("insert investment history: %w", err)
		}

		result = entities.ValuationResult{
			HistoryID:       history.ID,
			InvestmentID:    model.ID,
			PreviousBalance: previousBalance,
			CurrentValue:    input.NewBalance,
			Delta:           input.Delta,
			Balance:         input.NewBalance,
			Date:            now,
		}

		return nil
	})
	if err != nil {
		return entities.ValuationResult{}, err
	}

	return result, nil
}

func toInvestmentEntity(model entities.Investment) entities.InvestmentEntity {
	return entities.InvestmentEntity{
		ID:             model.ID,
		UserID:         model.UserID,
		Name:           model.Name,
		Balance:        model.Balance,
		InitialBalance: model.InitialBalance,
		Tags:           model.Tags.Strings(),
		CreatedAt:      model.CreatedAt,
	}
}
