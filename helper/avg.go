package helper

// Avg returns average of given int
func Avg(ii ...uint8) float32 {
	var s uint8
	for _, i := range ii {
		s += i
	}

	return float32(s) / float32(len(ii))
}
