package exercises

// CalculateScoreUp calculates the score for the tree at the given location.
func CalculateScoreUp(h, r, c int, m [][]int) int {
	switch {
	case r == 0:
		return 0
	case h > m[r-1][c]:
		return 1 + CalculateScoreUp(h, r-1, c, m)
	default:
		return 1
	}
}

// CalculateScoreDown calculates the score for the tree at the given location.
func CalculateScoreDown(h, r, c int, m [][]int) int {
	switch {
	case r == dimY-1:
		return 0
	case h > m[r+1][c]:
		return 1 + CalculateScoreDown(h, r+1, c, m)
	default:
		return 1
	}
}

// CalculateScoreLeft calculates the score for the tree at the given location.
func CalculateScoreLeft(h, r, c int, m [][]int) int {
	switch {
	case c == 0:
		return 0
	case h > m[r][c-1]:
		return 1 + CalculateScoreLeft(h, r, c-1, m)
	default:
		return 1
	}
}

// CalculateScoreRight calculates the score for the tree at the given location.
func CalculateScoreRight(h, r, c int, m [][]int) int {
	switch {
	case c == dimX-1:
		return 0
	case h > m[r][c+1]:
		return 1 + CalculateScoreRight(h, r, c+1, m)
	default:
		return 1
	}
}
