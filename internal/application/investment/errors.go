package investment

import "errors"

var (
	ErrInvestmentNotFound  = errors.New("investment not found")
	ErrInsufficientBalance = errors.New("insufficient balance")
	ErrValueUnchanged      = errors.New("current value is the same as the investment balance")
)
