package exercises

import "strings"

type Pattern struct {
	Col []string
	Row []string
}

func parsePattern(s string) *Pattern {
	lines := strings.Split(s, "\n")

	rows := make([]string, len(lines))
	cols := make([]string, len(lines[0]))

	for y, line := range lines {
		rows[y] = line

		for x, char := range line {
			cols[x] += string(char)
		}
	}

	return &Pattern{
		Col: cols,
		Row: rows,
	}
}

func getPatterns(input string) []*Pattern {
	sections := strings.Split(input, "\n\n")
	patterns := make([]*Pattern, len(sections))

	for i, s := range sections {
		patterns[i] = parsePattern(s)
	}

	return patterns
}

// returns the horizontal plane of the pattern; 1-indexed
func findMirror(a []string, hasSmudge bool) (int, int) {
	var n int

	for r := 0; r < len(a)-1; r++ {
		cmp := expandingCompare(a, r, hasSmudge)
		if cmp == 1 && hasSmudge {
			return r + 1, cmp
		}

		// clean match doesn't immediately return in case we find a smudge later
		// TODO: maybe return sooner if we find clean but there are no smudges
		if cmp == 0 {
			n = r + 1
		}
	}

	if n > 0 {
		return n, 0
	}

	// no match found
	return 0, -1
}

// -1=invalid; 0=found clean; 1=found smudge
func expandingCompare(ss []string, n int, allowMismatch bool) int {
	if len(ss) == 0 {
		return -1
	}

	var smudgeFound bool

	for i, j := n, n+1; i >= 0 && j < len(ss); i, j = i-1, j+1 {
		found := compareWithSmudge(ss[i], ss[j])

		if found == -1 {
			return -1
		}

		if found == 1 && (!allowMismatch || smudgeFound) {
			return -1
		}

		if found == 1 {
			smudgeFound = true
		}

		// no check for found == 0 because we want to keep going
	}

	if smudgeFound {
		return 1
	}

	return 0
}

// -1 too many invalid, 0 no mismatch, 1 single mismatch
func compareWithSmudge(a, b string) int {
	mismatchCount := 0

	for q := 0; q < len(a); q++ {
		if a[q] != b[q] {
			mismatchCount++

			if mismatchCount > 1 {
				return -1
			}
		}
	}

	return mismatchCount
}
