package entities

import "time"

type RecordMovementInput struct {
	InvestmentID uint
	MovementType MovementType
	Amount       float64
	NewBalance   float64
}

type MovementResult struct {
	HistoryID    uint
	InvestmentID uint
	MovementType MovementType
	Amount       float64
	Balance      float64
	Date         time.Time
}

type RecordValuationInput struct {
	InvestmentID uint
	Delta        float64
	NewBalance   float64
}

type ValuationResult struct {
	HistoryID        uint
	InvestmentID     uint
	PreviousBalance  float64
	CurrentValue     float64
	Delta            float64
	Balance          float64
	Date             time.Time
}
