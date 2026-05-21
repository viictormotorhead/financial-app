package request

import (
	"errors"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/viictormotorhead/financial-app/internal/infra/repositories/investment/entities"
)

type CreateMovementRequest struct {
	Type   string  `json:"type"`
	Amount float64 `json:"amount"`
}

func (r CreateMovementRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Type,
			validation.Required.Error("type is required"),
			validation.By(validateMovementType),
		),
		validation.Field(&r.Amount,
			validation.Required.Error("amount is required"),
			validation.Min(0.01).Error("amount must be greater than zero"),
		),
	)
}

func (r CreateMovementRequest) MovementType() entities.MovementType {
	return entities.MovementType(r.Type)
}

func validateMovementType(value interface{}) error {
	movementType, _ := value.(string)
	switch entities.MovementType(movementType) {
	case entities.MovementTypeDeposit, entities.MovementTypeWithdrawal:
		return nil
	default:
		return errors.New("type must be deposit or withdrawal")
	}
}
