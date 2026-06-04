package request

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type RegisterRequest struct {
	Name     string `json:"name"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func (r RegisterRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Name,
			validation.Required.Error("name is required"),
			validation.Length(1, 255).Error("name must be between 1 and 255 characters"),
		),
		validation.Field(&r.Username,
			validation.Required.Error("username is required"),
			validation.Length(3, 100).Error("username must be between 3 and 100 characters"),
		),
		validation.Field(&r.Password,
			validation.Required.Error("password is required"),
			validation.Length(8, 128).Error("password must be between 8 and 128 characters"),
		),
	)
}
