package gormrepo

import (
	"context"
	"fmt"

	"github.com/viictormotorhead/financial-app/internal/infra/repositories/investment/entities"
)

func (r *investmentWriteRepository) ListHistoryByInvestmentID(ctx context.Context, investmentID uint) ([]entities.InvestmentHistoryEntity, error) {
	var rows []entities.InvestmentHistory
	if err := r.db.WithContext(ctx).
		Where("investment_id = ?", investmentID).
		Order("date ASC, id ASC").
		Find(&rows).Error; err != nil {
		return nil, fmt.Errorf("list investment history: %w", err)
	}

	result := make([]entities.InvestmentHistoryEntity, 0, len(rows))
	for _, row := range rows {
		result = append(result, toInvestmentHistoryEntity(row))
	}

	return result, nil
}

func (r *investmentWriteRepository) SumEarningsByInvestmentID(ctx context.Context, investmentID uint) (float64, error) {
	var total float64
	err := r.db.WithContext(ctx).
		Model(&entities.InvestmentHistory{}).
		Select("COALESCE(SUM(amount), 0)").
		Where("investment_id = ?", investmentID).
		Where("movement_type = ?", entities.MovementTypeEarning).
		Scan(&total).Error
	if err != nil {
		return 0, fmt.Errorf("sum earning movements: %w", err)
	}

	return total, nil
}

func toInvestmentHistoryEntity(row entities.InvestmentHistory) entities.InvestmentHistoryEntity {
	return entities.InvestmentHistoryEntity{
		ID:           row.ID,
		InvestmentID: row.InvestmentID,
		Date:         row.Date,
		Amount:       row.Amount,
		BalanceAfter: row.BalanceAfter,
		MovementType: row.MovementType,
	}
}
