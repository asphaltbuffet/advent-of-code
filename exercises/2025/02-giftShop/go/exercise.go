package exercises

import (
	"fmt"
	"iter"
	"math"
	"slices"
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
		for id := range InvalidIds(p.Lower, p.Upper) {
			sum += id
		}
	}

	return sum, nil
}

// Two returns the answer to the second part of the exercise.
// 20942028300 -> too high
func (e Exercise) Two(instr string) (any, error) {
	var sum int

	pp, err := parseInput(instr)
	if err != nil {
		return nil, fmt.Errorf("invalid input: %w", err)
	}

	for _, p := range pp {
		rep, err := p.Repeated()
		if err != nil {
			return nil, fmt.Errorf("generating repeats for %d-%d: %w", p.Lower, p.Upper, err)
		}
		for _, id := range rep {
			sum += id
		}
	}

	return sum, nil
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

func isRepeatedPattern(n int) bool {
	if n <= 0 {
		return false
	}
	s := strconv.Itoa(n)

	for patLen := 1; patLen <= len(s)/2; patLen++ {
		if len(s)%patLen != 0 {
			continue
		}
		pat := s[:patLen]
		ok := true
		for i := patLen; i < len(s); i += patLen {
			if s[i:i+patLen] != pat {
				ok = false
				break
			}
		}
		if ok {
			return true
		}
	}
	return false
}

func InvalidIds(a, b int) iter.Seq[int] {
	return func(yield func(int) bool) {
		for i := a; i <= b; i = NextInvalid(i) {
			if i == a && !IsInvalid(i) {
				continue
			}

			if !yield(i) {
				return
			}
		}
	}
}

func (p Pair) Repeated() ([]int, error) {
	if p.Lower <= 0 || p.Lower > p.Upper {
		return nil, fmt.Errorf("%d-%d: invalid pair", p.Lower, p.Upper)
	}

	llen, ulen := intLen(p.Lower), intLen(p.Upper)

	// speed up with a lookup table
	powTen := make([]int, ulen+1)
	powTen[0] = 1
	for k := 1; k <= ulen; k++ {
		powTen[k] = powTen[k-1] * 10
	}

	rr := []int{}

	for numLen := max(llen, 2); numLen <= ulen; numLen++ {
		for patLen := 1; patLen <= ulen/2; patLen++ {
			// skip if it won't make any usable patterns
			if numLen%patLen != 0 {
				continue
			}

			start := powTen[patLen-1]
			end := powTen[patLen] - 1

			for j := start; j <= end; j++ {
				n := 0
				repeats := numLen / patLen
				for k := 0; k < repeats; k++ {
					n = n*powTen[patLen] + j
				}
				if n < p.Lower {
					continue
				}

				if n > p.Upper {
					break
				}

				rr = append(rr, n)
			}
		}
	}

	slices.Sort(rr)
	return slices.Compact(rr), nil
}

func intLen(x int) int {
	return int(math.Log10(float64(x)) + 1)
}

func IsRepeating(n int) bool {
	if n < 10 {
		return false
	}

	s := strconv.Itoa(n)
	for i := 1; i <= len(s)/2; i++ {
		r := strings.Repeat(s[:i], len(s)/i)
		if r == s {
			return true
		}
	}
	return false
}
