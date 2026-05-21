package use_cases

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/viictormotorhead/financial-app/internal/application/investment/dto/outputs"
	"github.com/viictormotorhead/financial-app/internal/application/investment/repositories"
	"github.com/viictormotorhead/financial-app/internal/infra/repositories/investment/entities"
)

type CreateInvestmentCommand struct {
	Name    string
	Balance float64
	Tags    []string
}

type CreateInvestmentUseCaseIF interface {
	Create(ctx context.Context, cmd CreateInvestmentCommand) (outputs.CreateInvestmentOutputDTO, error)
}

type CreateInvestmentUseCaseImpl struct {
	repository repositories.InvestmentWriteRepositoryIF
}

func NewCreateInvestmentUseCase(repository repositories.InvestmentWriteRepositoryIF) CreateInvestmentUseCaseIF {
	return &CreateInvestmentUseCaseImpl{repository: repository}
}

func (u *CreateInvestmentUseCaseImpl) Create(ctx context.Context, cmd CreateInvestmentCommand) (outputs.CreateInvestmentOutputDTO, error) {
	entity := entities.InvestmentEntity{
		Name:      cmd.Name,
		Balance:   cmd.Balance,
		Tags:      normalizeTags(cmd.Tags),
		CreatedAt: time.Now().UTC(),
	}

	saved, err := u.repository.Save(ctx, entity)
	if err != nil {
		return outputs.CreateInvestmentOutputDTO{}, fmt.Errorf("create investment: %w", err)
	}

	return outputs.CreateInvestmentOutputDTO{
		ID:        saved.ID,
		Name:      saved.Name,
		Balance:   saved.Balance,
		Tags:      saved.Tags,
		CreatedAt: saved.CreatedAt,
	}, nil
}

func normalizeTags(tags []string) []string {
	if len(tags) == 0 {
		return []string{}
	}

	seen := make(map[string]struct{}, len(tags))
	result := make([]string, 0, len(tags))
	for _, tag := range tags {
		trimmed := strings.TrimSpace(tag)
		if trimmed == "" {
			continue
		}
		key := strings.ToLower(trimmed)
		if _, exists := seen[key]; exists {
			continue
		}
		seen[key] = struct{}{}
		result = append(result, trimmed)
	}

	return result
}
