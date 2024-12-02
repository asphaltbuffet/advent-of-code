package exercises

import (
	"errors"
	"fmt"
	"strings"
)

func parse(s string) ([]int, []int, error) {
	if len(s) == 0 {
		return nil, nil, errors.New("empty string")
	}

	lines := strings.Split(s, "\n")

	a := make([]int, 0, len(lines))
	b := make([]int, 0, len(lines))

	var j, k int

	for i, line := range lines {
		_, err := fmt.Sscanf(line, "%d   %d", &j, &k)
		if err != nil {
			return nil, nil, fmt.Errorf("line %d: %w", i+1, err)
		}

		a = append(a, j)
		b = append(b, k)
	}

	return a, b, nil
}

func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}
