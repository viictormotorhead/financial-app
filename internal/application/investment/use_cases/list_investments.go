package use_cases

import (
	"context"
	"fmt"
	"math"
	"sort"

	"github.com/viictormotorhead/financial-app/internal/application/investment/dto/outputs"
	"github.com/viictormotorhead/financial-app/internal/application/investment/repositories"
)

type ListInvestmentsQuery struct {
	Tags []string
}

type ListInvestmentsUseCaseIF interface {
	List(ctx context.Context, query ListInvestmentsQuery) (outputs.ListInvestmentsOutputDTO, error)
}

type ListInvestmentsUseCaseImpl struct {
	repository repositories.InvestmentWriteRepositoryIF
}

func NewListInvestmentsUseCase(repository repositories.InvestmentWriteRepositoryIF) ListInvestmentsUseCaseIF {
	return &ListInvestmentsUseCaseImpl{repository: repository}
}

func (u *ListInvestmentsUseCaseImpl) List(ctx context.Context, query ListInvestmentsQuery) (outputs.ListInvestmentsOutputDTO, error) {
	investments, err := u.repository.List(ctx, normalizeTags(query.Tags))
	if err != nil {
		return outputs.ListInvestmentsOutputDTO{}, fmt.Errorf("list investments: %w", err)
	}

	var total float64
	for _, inv := range investments {
		total += inv.Balance
	}

	items := make([]outputs.InvestmentAllocationDTO, 0, len(investments))
	for _, inv := range investments {
		percentage := 0.0
		if total > 0 {
			percentage = (inv.Balance / total) * 100
		}

		tags := inv.Tags
		if tags == nil {
			tags = []string{}
		}

		items = append(items, outputs.InvestmentAllocationDTO{
			ID:                inv.ID,
			Investment:        inv.Name,
			Amount:            inv.Balance,
			Percentage:        roundToTwoDecimals(percentage),
			PercentageGrowing: roundToTwoDecimals(percentageGrowingFromEarnings(inv.Balance, inv.EarningsTotal)),
			Tags:              tags,
		})
	}

	sort.Slice(items, func(i, j int) bool {
		gi := items[i].PercentageGrowing
		gj := items[j].PercentageGrowing

		// NaN should never happen, but if it does, push it to the bottom.
		if math.IsNaN(gi) {
			gi = math.Inf(-1)
		}
		if math.IsNaN(gj) {
			gj = math.Inf(-1)
		}

		if gi == gj {
			return items[i].Amount > items[j].Amount
		}
		return gi > gj
	})

	return outputs.ListInvestmentsOutputDTO{Items: items}, nil
}

func roundToTwoDecimals(value float64) float64 {
	return math.Round(value*100) / 100
}

// Rendimiento por valuaciones: suma de earnings vs capital invertido (balance sin earnings).
func percentageGrowingFromEarnings(balance, earningsTotal float64) float64 {
	investedCapital := balance - earningsTotal
	if investedCapital == 0 {
		return 0
	}
	return (earningsTotal / investedCapital) * 100
}
