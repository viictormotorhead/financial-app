package use_cases

import (
	"context"
	"fmt"
	"time"

	"github.com/viictormotorhead/financial-app/internal/application/investment/dto/outputs"
	"github.com/viictormotorhead/financial-app/internal/application/investment/repositories"
	"github.com/viictormotorhead/financial-app/internal/infra/repositories/investment/entities"
)

type CreateInvestmentCommand struct {
	Name    string
	Balance float64
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
	if cmd.Name == "" {
		return outputs.CreateInvestmentOutputDTO{}, fmt.Errorf("name is required")
	}
	if cmd.Balance <= 0 {
		return outputs.CreateInvestmentOutputDTO{}, fmt.Errorf("balance must be greater than zero")
	}

	entity := entities.InvestmentEntity{
		Name:      cmd.Name,
		Balance:   cmd.Balance,
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
		CreatedAt: saved.CreatedAt,
	}, nil
}
