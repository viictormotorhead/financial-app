package request

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type CreateValuationRequest struct {
	CurrentValue float64 `json:"current_value"`
}

func (r CreateValuationRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.CurrentValue,
			validation.Min(0.0).Error("current_value must be greater than or equal to zero"),
		),
	)
}
