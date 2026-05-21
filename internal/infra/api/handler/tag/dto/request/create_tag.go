package request

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type CreateTagRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (r CreateTagRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Name,
			validation.Required.Error("name is required"),
			validation.Length(1, 255).Error("name must be between 1 and 255 characters"),
		),
		validation.Field(&r.Description,
			validation.Length(0, 500).Error("description must be at most 500 characters"),
		),
	)
}
