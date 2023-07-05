package utilities

// MeanIntSlice returns the mean of a slice of floats.
func MeanFloatSlice(arr []float64) float64 {
	var sum float64
	for _, v := range arr {
		sum += v
	}

	return sum / float64(len(arr))
}

// MinIntSlice returns the minimum value of a slice of floats.
func MinFloatSlice(arr []float64) float64 {
	min := arr[0]
	for _, v := range arr {
		if v < min {
			min = v
		}
	}

	return min
}

// MaxIntSlice returns the maximum value of a slice of floats.
func MaxFloatSlice(arr []float64) float64 {
	max := arr[0]
	for _, v := range arr {
		if v > max {
			max = v
		}
	}

	return max
}
