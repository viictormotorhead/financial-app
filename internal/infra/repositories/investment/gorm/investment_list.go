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

	result := make([]entities.InvestmentEntity, 0, len(models))
	for _, model := range models {
		result = append(result, toInvestmentEntity(model))
	}

	return result, nil
}
