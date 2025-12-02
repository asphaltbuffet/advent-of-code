package exercises

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

// Exercise for Advent of Code 2025 day 2.
type Exercise struct {
	common.BaseExercise
}

// One returns the answer to the first part of the exercise.
// wrong: 13484414229 (too high)
func (e Exercise) One(instr string) (any, error) {
	// notes:
	// - invalid must have even # of digits
	var sum int

	pp, err := parseInput(instr)
	if err != nil {
		return nil, fmt.Errorf("invalid input: %w", err)
	}

	for _, p := range pp {
		sum += p.sumInvalid()
	}

	return sum, nil
}

// Two returns the answer to the second part of the exercise.
func (e Exercise) Two(instr string) (any, error) {
	return nil, fmt.Errorf("part 2 not implemented")
}

type Pair struct {
	Lower       int
	Upper       int
	LowerString string
	UpperString string
}

func parseInput(s string) ([]Pair, error) {
	pp := strings.Split(s, ",")
	pairs := make([]Pair, len(pp))

	for i, p := range pp {
		raw := strings.Split(p, "-")

		l, err := strconv.Atoi(raw[0])
		if err != nil {
			return nil, fmt.Errorf("pair[%d][0]=%s", i, raw[0])
		}

		u, err := strconv.Atoi(raw[1])
		if err != nil {
			return nil, fmt.Errorf("pair[%d][1]=%s", i, raw[1])
		}

		pairs[i] = Pair{
			Lower:       l,
			Upper:       u,
			LowerString: raw[0],
			UpperString: raw[1],
		}
	}

	return pairs, nil
}

func (p Pair) sumInvalid() int {
	var sum int

	for _, id := range p.findInvalid() {
		sum += id
	}

	return sum
}

func (p Pair) findInvalid() []int {
	//  XXXX => AA,BB
	// XXXXX => 0AA,BBB
	split := len(p.LowerString) / 2
	if len(p.LowerString)%2 != 0 {
		split++
	}

	mult := pow10(split)
	// right := p.Lower % mult

	left := p.Lower / mult
	n := left*mult + left // AA00+AA
	for n < p.Lower {
		left++
		if left >= mult { // we added a digit to left
			mult *= 10
		}
		n = left*mult + left
	}

	invalids := make([]int, 0, 0)

	for n <= p.Upper { // NNNN <= BBBB
		// fmt.Println(n)
		invalids = append(invalids, n) // AAAA

		left++
		if left >= mult { // we added a digit to left
			mult *= 10
		}
		n = left*mult + left // BBBB

	}

	// fmt.Println(invalids)

	return invalids
}

func pow10(m int) int {
	switch m {
	case 0:
		return 1
	case 1:
		return 10
	default:
		result := 10
		for i := 2; i <= m; i++ {
			result *= 10
		}
		return result
	}
}

func NextInvalid(n int) int {
	if n == 0 {
		return 11
	}

	var digits int = int(math.Log10(float64(n)) + 1)
	if digits%2 != 0 {
		n = pow10(digits)
		digits++
	}

	var split = digits/2 + digits%2

	mult := pow10(split)

	left, right := n/mult, n%mult
	if left <= right {
		left++
	}

	if left >= mult { // AAA BB => AAA BBB
		mult *= 10
	}

	return left*mult + left // AA00+AA

}

func IsInvalid(n int) bool {
	s := fmt.Sprintf("%d", n)
	if len(s)%2 != 0 {
		return false
	}

	div := len(s) / 2
	l, r := s[:div], s[div:]
	return l == r
}

func InvalidIds(min, max int) func(yield func(int) bool) {
	return func(yield func(int) bool) {
		for i := min; i <= max; NextInvalid(i) {
			if i == min && !IsInvalid(i) {
				i = NextInvalid(i)
			}

			if !yield(i) {
				return
			}
		}
	}
}
