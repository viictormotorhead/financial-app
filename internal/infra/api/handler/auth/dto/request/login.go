package request

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (r LoginRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Username,
			validation.Required.Error("username is required"),
		),
		validation.Field(&r.Password,
			validation.Required.Error("password is required"),
		),
	)
}
