package exercises

import (
	"strings"
)

func CountSplits(input string) int {
	var splits int

	lines := strings.Split(input, "\n")
	beams := make(map[int]bool, len(lines[0])+2)

	// find start
	for i, c := range lines[0] {
		if c == 'S' {
			beams[i] = true
			break
		}
	}

	for _, row := range lines[1:] {
		for x := 0; x < len(row); x++ {
			if !beams[x] {
				continue
			}

			if row[x] == '^' {
				beams[x-1] = true
				beams[x+1] = true
				beams[x] = false
				x++

				splits++
			}
		}
	}

	return splits
}
