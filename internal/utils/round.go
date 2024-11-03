package utils

import "math"

func Round(value float64, precision int) float64 {
	factor := math.Pow(10, float64(precision))
	return math.Round(value*factor) / factor
}
