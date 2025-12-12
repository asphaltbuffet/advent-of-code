package exercises

import (
	"strings"
)

var op = map[string]func(nums ...int) int{
	"+": func(nums ...int) int {
		sum := 0
		for _, n := range nums {
			sum += n
		}
		return sum
	},
	"*": func(nums ...int) int {
		prod := 1
		for _, n := range nums {
			prod *= n
		}

		return prod
	},
}

func LoadHomework(s string) ([]string, []string, int) {
	tokens := strings.Fields(s)
	lCount := strings.Count(s, "\n") + 1

	width := len(tokens) / lCount
	opIdx := width * (lCount - 1)

	return tokens[:opIdx], tokens[opIdx:], width
}
