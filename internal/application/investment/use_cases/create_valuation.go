package use_cases

import (
	"context"
	"errors"
	"fmt"

	"github.com/viictormotorhead/financial-app/internal/application/auth"
	appinvestment "github.com/viictormotorhead/financial-app/internal/application/investment"
	"github.com/viictormotorhead/financial-app/internal/application/investment/dto/outputs"
	"github.com/viictormotorhead/financial-app/internal/application/investment/repositories"
	"github.com/viictormotorhead/financial-app/internal/infra/repositories/investment/entities"
)

type CreateValuationCommand struct {
	InvestmentID uint
	CurrentValue float64
}

type CreateValuationUseCaseIF interface {
	Create(ctx context.Context, cmd CreateValuationCommand) (outputs.CreateValuationOutputDTO, error)
}

type CreateValuationUseCaseImpl struct {
	repository repositories.InvestmentWriteRepositoryIF
}

func NewCreateValuationUseCase(repository repositories.InvestmentWriteRepositoryIF) CreateValuationUseCaseIF {
	return &CreateValuationUseCaseImpl{repository: repository}
}

func (u *CreateValuationUseCaseImpl) Create(ctx context.Context, cmd CreateValuationCommand) (outputs.CreateValuationOutputDTO, error) {
	userID, err := auth.RequireUserID(ctx)
	if err != nil {
		return outputs.CreateValuationOutputDTO{}, err
	}

	investment, err := u.repository.FindByID(ctx, userID, cmd.InvestmentID)
	if err != nil {
		if errors.Is(err, appinvestment.ErrInvestmentNotFound) {
			return outputs.CreateValuationOutputDTO{}, err
		}
		return outputs.CreateValuationOutputDTO{}, fmt.Errorf("find investment: %w", err)
	}

	delta := cmd.CurrentValue - investment.Balance
	if delta == 0 {
		return outputs.CreateValuationOutputDTO{}, appinvestment.ErrValueUnchanged
	}

	result, err := u.repository.RecordValuation(ctx, entities.RecordValuationInput{
		InvestmentID: cmd.InvestmentID,
		Delta:        delta,
		NewBalance:   cmd.CurrentValue,
	})
	if err != nil {
		return outputs.CreateValuationOutputDTO{}, fmt.Errorf("record valuation: %w", err)
	}

	return outputs.CreateValuationOutputDTO{
		HistoryID:       result.HistoryID,
		InvestmentID:    result.InvestmentID,
		PreviousBalance: result.PreviousBalance,
		CurrentValue:    result.CurrentValue,
		Delta:           result.Delta,
		Balance:         result.Balance,
		Date:            result.Date,
	}, nil
}
