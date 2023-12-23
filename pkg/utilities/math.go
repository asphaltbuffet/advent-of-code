package utilities

// MeanFloatSlice returns the mean of a slice of floats.
func MeanFloatSlice(arr []float64) float64 {
	var sum float64
	for _, v := range arr {
		sum += v
	}

	return sum / float64(len(arr))
}

// MinFloatSlice returns the minimum value of a slice of floats.
func MinFloatSlice(arr []float64) float64 {
	min := arr[0]
	for _, v := range arr {
		if v < min {
			min = v
		}
	}

	return min
}

// MaxFloatSlice returns the maximum value of a slice of floats.
func MaxFloatSlice(arr []float64) float64 {
	max := arr[0]
	for _, v := range arr {
		if v > max {
			max = v
		}
	}

	return max
}

// greatest common divisor (GCD) via Euclidean algorithm
func GCD(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

// least common multiple (LCM)
func LCM(a, b int) int {
	return a / GCD(a, b) * b
}
