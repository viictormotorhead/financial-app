package use_cases

import "math"

func roundToTwoDecimals(value float64) float64 {
	return math.Round(value*100) / 100
}

// Rendimiento por valuaciones: suma de earnings vs capital invertido (balance sin earnings).
func percentageGrowingFromEarnings(balance, earningsTotal float64) float64 {
	investedCapital := balance - earningsTotal
	if investedCapital == 0 {
		return 0
	}
	return (earningsTotal / investedCapital) * 100
}
