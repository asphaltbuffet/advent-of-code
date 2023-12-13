package exercises

import (
	"strconv"
	"strings"
)

func lineToIntSlice(line string) []int {
	var out []int

	tokens := strings.Fields(line)
	for _, token := range tokens {
		n, _ := strconv.Atoi(token)
		out = append(out, n)
	}

	return out
}

func reduceToDiffs(in []int) []int {
	var out []int

	for i := 1; i < len(in); i++ {
		out = append(out, in[i]-in[i-1])
	}

	return out
}

func calculateReductions(history []int) int {
	h := history

	var sum int

	for {
		// fmt.Println(h)

		sum += h[len(h)-1]

		h = reduceToDiffs(h)
		if allZero(h) {
			break
		}
	}

	return sum
}

func allZero(x []int) bool {
	for _, i := range x {
		if i != 0 {
			return false
		}
	}

	return true
}
