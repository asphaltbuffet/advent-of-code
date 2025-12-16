package exercises

import (
	"strings"
	"unicode"
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

type Problem struct {
	Numbers  []int
	Operator string
}

type Operator rune

const (
	OpEmpty Operator = ' '
	OpPlus           = '+'
	OpMult           = '*'
)

func RTLParse(s string) ([]Problem, error) {
	lines := strings.Split(s, "\n")

	opMap := make(map[int]string)
	for i, o := range lines[len(lines)-1] {
		if !unicode.IsSpace(o) {
			opMap[i] = string(o)
		}
	}

	problems := []Problem{}

	nums := []int{}
	// read columns RtL
	for col := len(lines[0]) - 1; col >= 0; col-- {
		n := 0
		// just read the number digits
		for row := 0; row < len(lines)-1; row++ {
			r := rune(lines[row][col])
			if !unicode.IsSpace(r) {
				n = n*10 + int(r-'0')
			}
		}

		nums = append(nums, n)

		// if there's an operator, store current data and move to next problem
		if operator, ok := opMap[col]; ok {
			p := Problem{
				Numbers:  nums,
				Operator: operator,
			}

			problems = append(problems, p)

			// reset nums
			nums = []int{}
			// skip an extra line since next one is empty
			col--
		}
	}

	return problems, nil
}
