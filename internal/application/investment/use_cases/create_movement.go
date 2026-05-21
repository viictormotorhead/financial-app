package use_cases

import (
	"context"
	"errors"
	"fmt"

	appinvestment "github.com/viictormotorhead/financial-app/internal/application/investment"
	"github.com/viictormotorhead/financial-app/internal/application/investment/dto/outputs"
	"github.com/viictormotorhead/financial-app/internal/application/investment/repositories"
	"github.com/viictormotorhead/financial-app/internal/infra/repositories/investment/entities"
)

type CreateMovementCommand struct {
	InvestmentID uint
	Type         entities.MovementType
	Amount       float64
}

type CreateMovementUseCaseIF interface {
	Create(ctx context.Context, cmd CreateMovementCommand) (outputs.CreateMovementOutputDTO, error)
}

type CreateMovementUseCaseImpl struct {
	repository repositories.InvestmentWriteRepositoryIF
}

func NewCreateMovementUseCase(repository repositories.InvestmentWriteRepositoryIF) CreateMovementUseCaseIF {
	return &CreateMovementUseCaseImpl{repository: repository}
}

func (u *CreateMovementUseCaseImpl) Create(ctx context.Context, cmd CreateMovementCommand) (outputs.CreateMovementOutputDTO, error) {
	investment, err := u.repository.FindByID(ctx, cmd.InvestmentID)
	if err != nil {
		if errors.Is(err, appinvestment.ErrInvestmentNotFound) {
			return outputs.CreateMovementOutputDTO{}, err
		}
		return outputs.CreateMovementOutputDTO{}, fmt.Errorf("find investment: %w", err)
	}

	var newBalance float64
	switch cmd.Type {
	case entities.MovementTypeDeposit:
		newBalance = investment.Balance + cmd.Amount
	case entities.MovementTypeWithdrawal:
		if cmd.Amount > investment.Balance {
			return outputs.CreateMovementOutputDTO{}, appinvestment.ErrInsufficientBalance
		}
		newBalance = investment.Balance - cmd.Amount
	default:
		return outputs.CreateMovementOutputDTO{}, fmt.Errorf("unsupported movement type: %s", cmd.Type)
	}

	result, err := u.repository.RecordMovement(ctx, entities.RecordMovementInput{
		InvestmentID: cmd.InvestmentID,
		MovementType: cmd.Type,
		Amount:       cmd.Amount,
		NewBalance:   newBalance,
	})
	if err != nil {
		return outputs.CreateMovementOutputDTO{}, fmt.Errorf("record movement: %w", err)
	}

	return outputs.CreateMovementOutputDTO{
		HistoryID:    result.HistoryID,
		InvestmentID: result.InvestmentID,
		Type:         result.MovementType,
		Amount:       result.Amount,
		Balance:      result.Balance,
		Date:         result.Date,
	}, nil
}
