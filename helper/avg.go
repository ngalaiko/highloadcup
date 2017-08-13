package helper

import "math"

// Avg returns average of given int
func Avg(ii ...uint8) float64 {
	var s uint8
	for _, i := range ii {
		s += i
	}

	return Round(float64(s)/float64(len(ii)), 5)
}

func Round(x float64, prec int) float64 {
	var rounder float64
	pow := math.Pow(10, float64(prec))
	intermed := x * pow
	_, frac := math.Modf(intermed)
	if frac >= 0.5 {
		rounder = math.Ceil(intermed)
	} else {
		rounder = math.Floor(intermed)
	}

	return rounder / pow
}
