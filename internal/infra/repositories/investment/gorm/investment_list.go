package gormrepo

import (
	"context"
	"fmt"
	"strings"

	"github.com/viictormotorhead/financial-app/internal/infra/repositories/investment/entities"
)

func (r *investmentWriteRepository) List(ctx context.Context, tagNames []string) ([]entities.InvestmentEntity, error) {
	var models []entities.Investment
	query := r.db.WithContext(ctx).Model(&entities.Investment{})

	if len(tagNames) > 0 {
		keys := make([]string, len(tagNames))
		for i, name := range tagNames {
			keys[i] = strings.ToLower(strings.TrimSpace(name))
		}

		query = query.Where(`EXISTS (
			SELECT 1 FROM unnest(tags) AS inv_tag
			WHERE LOWER(inv_tag) IN ?
		)`, keys)
	}

	if err := query.Order("name ASC").Find(&models).Error; err != nil {
		return nil, fmt.Errorf("list investments: %w", err)
	}

	earningsByID, err := r.sumEarningsByInvestmentIDs(ctx, models)
	if err != nil {
		return nil, err
	}

	result := make([]entities.InvestmentEntity, 0, len(models))
	for _, model := range models {
		entity := toInvestmentEntity(model)
		entity.EarningsTotal = earningsByID[model.ID]
		result = append(result, entity)
	}

	return result, nil
}

type earningsSumRow struct {
	InvestmentID uint    `gorm:"column:investment_id"`
	Total        float64 `gorm:"column:total"`
}

func (r *investmentWriteRepository) sumEarningsByInvestmentIDs(ctx context.Context, models []entities.Investment) (map[uint]float64, error) {
	earningsByID := make(map[uint]float64, len(models))
	if len(models) == 0 {
		return earningsByID, nil
	}

	ids := make([]uint, len(models))
	for i, model := range models {
		ids[i] = model.ID
		earningsByID[model.ID] = 0
	}

	var rows []earningsSumRow
	err := r.db.WithContext(ctx).
		Model(&entities.InvestmentHistory{}).
		Select("investment_id, COALESCE(SUM(amount), 0) AS total").
		Where("investment_id IN ?", ids).
		Where("movement_type = ?", entities.MovementTypeEarning).
		Group("investment_id").
		Scan(&rows).Error
	if err != nil {
		return nil, fmt.Errorf("sum earning movements: %w", err)
	}

	for _, row := range rows {
		earningsByID[row.InvestmentID] = row.Total
	}

	return earningsByID, nil
}
