package entities

type MovementType string

const (
	MovementTypeDeposit    MovementType = "deposit"
	MovementTypeWithdrawal MovementType = "withdrawal"
	MovementTypeEarning    MovementType = "earning"
)

func (m MovementType) String() string {
	return string(m)
}

func (m MovementType) IsValid() bool {
	switch m {
	case MovementTypeDeposit, MovementTypeWithdrawal, MovementTypeEarning:
		return true
	default:
		return false
	}
}
