package exercises

import (
	"fmt"
	"strconv"
	"strings"
)

type Workflows map[string]*Workflow

type Workflow struct {
	Name  string
	Tests []*PartTest
}

type PartTest struct {
	Category   Category
	Comparator rune
	Threshold  int
	Test       TestFunc
	Dest       string
}

type (
	TestFunc func(int) bool
	Result   int
)

const (
	Accepted Result = iota + 1
	Rejected
)

type PartRange struct {
	start, end int
}

func parseWorkflow(in string) *Workflow {
	s := strings.TrimRight(in, "}")
	name, tests, ok := strings.Cut(s, "{")
	if !ok {
		return nil
	}

	tokens := strings.Split(tests, ",")

	w := &Workflow{
		Name:  name,
		Tests: make([]*PartTest, 0, len(tokens)+1),
	}

	for _, token := range tokens {
		var c Category
		var f func(int) bool
		var cmp rune
		var n int
		var err error

		// parse test
		t, d, isTest := strings.Cut(token, ":")
		if !isTest {
			c = Default
			cmp = '*' // any
			n = -1
			f = func(int) bool { return true }
			d = t
		} else {
			c = getCategory(rune(t[0]))
			if c == 0 {
				fmt.Printf("invalid category: %c\n", t[0])
			}

			n, err = strconv.Atoi(t[2:])
			if err != nil {
				fmt.Printf("invalid number: %s\n", t[2:])
			}

			cmp = rune(t[1])
			switch cmp {
			case '<':
				f = LessThan(n)

			case '>':
				f = GreaterThan(n)

			case '=':
				f = Equals(n)

			default:
				fmt.Printf("invalid test: %s\n", t)
				return nil
			}
		}

		w.Tests = append(w.Tests, &PartTest{
			Category:   c,
			Comparator: cmp,
			Threshold:  n,
			Test:       f,
			Dest:       d,
		})
	}

	// fmt.Printf("workflow %q: tests=%#v dests=%v\n", w.Name, w.Tests, w.Dests)
	return w
}

func LessThan(n int) func(int) bool {
	return func(i int) bool {
		return i < n
	}
}

func GreaterThan(n int) func(int) bool {
	return func(i int) bool {
		return i > n
	}
}

func Equals(n int) func(int) bool {
	return func(i int) bool {
		return i == n
	}
}

func countCombinations(workflows Workflows, flow string, ranges []PartRange) int {
	if flow == "R" {
		return 0
	}
	if flow == "A" {
		return partRangesProduct(ranges)
	}

	result := 0
	currentWorkflow := workflows[flow]
	for _, wfTest := range currentWorkflow.Tests {
		newRanges := make([]PartRange, len(ranges))
		copy(newRanges, ranges)
		rangeIndex := map[Category]int{ExtremlyCool: 0, Musical: 1, Aerodynamic: 2, Shiny: 3}[wfTest.Category]

		switch wfTest.Comparator {
		case '<':
			newRanges[rangeIndex].end = wfTest.Threshold - 1
			ranges[rangeIndex].start = wfTest.Threshold
			result += countCombinations(workflows, wfTest.Dest, newRanges)

		case '>':
			newRanges[rangeIndex].start = wfTest.Threshold + 1
			ranges[rangeIndex].end = wfTest.Threshold
			result += countCombinations(workflows, wfTest.Dest, newRanges)

		default:
			result += countCombinations(workflows, wfTest.Dest, ranges)
		}
	}

	return result
}

func partRangesProduct(ranges []PartRange) int {
	result := 1
	for _, r := range ranges {
		result *= r.end - r.start + 1
	}
	return result
}
