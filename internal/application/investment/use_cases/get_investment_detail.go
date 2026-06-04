package use_cases

import (
	"context"
	"errors"
	"fmt"
	"sort"

	"github.com/viictormotorhead/financial-app/internal/application/auth"
	appinvestment "github.com/viictormotorhead/financial-app/internal/application/investment"
	"github.com/viictormotorhead/financial-app/internal/application/investment/dto/outputs"
	"github.com/viictormotorhead/financial-app/internal/application/investment/repositories"
	"github.com/viictormotorhead/financial-app/internal/infra/repositories/investment/entities"
)

type GetInvestmentDetailQuery struct {
	InvestmentID uint
}

type GetInvestmentDetailUseCaseIF interface {
	Get(ctx context.Context, query GetInvestmentDetailQuery) (outputs.GetInvestmentDetailOutputDTO, error)
}

type GetInvestmentDetailUseCaseImpl struct {
	repository repositories.InvestmentWriteRepositoryIF
}

func NewGetInvestmentDetailUseCase(repository repositories.InvestmentWriteRepositoryIF) GetInvestmentDetailUseCaseIF {
	return &GetInvestmentDetailUseCaseImpl{repository: repository}
}

func (u *GetInvestmentDetailUseCaseImpl) Get(ctx context.Context, query GetInvestmentDetailQuery) (outputs.GetInvestmentDetailOutputDTO, error) {
	userID, err := auth.RequireUserID(ctx)
	if err != nil {
		return outputs.GetInvestmentDetailOutputDTO{}, err
	}

	investment, err := u.repository.FindByID(ctx, userID, query.InvestmentID)
	if err != nil {
		if errors.Is(err, appinvestment.ErrInvestmentNotFound) {
			return outputs.GetInvestmentDetailOutputDTO{}, err
		}
		return outputs.GetInvestmentDetailOutputDTO{}, fmt.Errorf("find investment: %w", err)
	}

	earningsTotal, err := u.repository.SumEarningsByInvestmentID(ctx, query.InvestmentID)
	if err != nil {
		return outputs.GetInvestmentDetailOutputDTO{}, fmt.Errorf("sum earnings: %w", err)
	}

	history, err := u.repository.ListHistoryByInvestmentID(ctx, query.InvestmentID)
	if err != nil {
		return outputs.GetInvestmentDetailOutputDTO{}, fmt.Errorf("list investment history: %w", err)
	}

	tags := investment.Tags
	if tags == nil {
		tags = []string{}
	}

	return outputs.GetInvestmentDetailOutputDTO{
		ID:            investment.ID,
		Name:          investment.Name,
		Tags:          tags,
		CurrentValue:  investment.Balance,
		GrowthPercent: roundToTwoDecimals(percentageGrowingFromEarnings(investment.Balance, earningsTotal)),
		GrowthAmount:  roundToTwoDecimals(earningsTotal),
		Series:        buildSeries(history),
		Movements:     buildMovements(history),
	}, nil
}

func buildSeries(history []entities.InvestmentHistoryEntity) []outputs.InvestmentSeriesPointDTO {
	sorted := append([]entities.InvestmentHistoryEntity(nil), history...)
	sort.Slice(sorted, func(i, j int) bool {
		if sorted[i].Date.Equal(sorted[j].Date) {
			return sorted[i].ID < sorted[j].ID
		}
		return sorted[i].Date.Before(sorted[j].Date)
	})

	series := make([]outputs.InvestmentSeriesPointDTO, 0, len(sorted))
	for _, row := range sorted {
		series = append(series, outputs.InvestmentSeriesPointDTO{
			Date:  row.Date,
			Value: row.BalanceAfter,
		})
	}

	return series
}

func buildMovements(history []entities.InvestmentHistoryEntity) []outputs.InvestmentMovementDTO {
	sorted := append([]entities.InvestmentHistoryEntity(nil), history...)
	sort.Slice(sorted, func(i, j int) bool {
		if sorted[i].Date.Equal(sorted[j].Date) {
			return sorted[i].ID > sorted[j].ID
		}
		return sorted[i].Date.After(sorted[j].Date)
	})

	movements := make([]outputs.InvestmentMovementDTO, 0, len(sorted))
	for _, row := range sorted {
		movements = append(movements, outputs.InvestmentMovementDTO{
			ID:           row.ID,
			Date:         row.Date,
			Type:         movementTypeToAPI(row.MovementType),
			Amount:       movementAmountToAPI(row.MovementType, row.Amount),
			BalanceAfter: row.BalanceAfter,
		})
	}

	return movements
}

func movementTypeToAPI(movementType entities.MovementType) string {
	if movementType == entities.MovementTypeEarning {
		return "value_update"
	}
	return movementType.String()
}

func movementAmountToAPI(movementType entities.MovementType, amount float64) *float64 {
	if movementType == entities.MovementTypeEarning {
		return nil
	}
	return &amount
}
