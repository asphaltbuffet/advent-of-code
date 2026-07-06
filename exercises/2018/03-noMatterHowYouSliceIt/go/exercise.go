package exercises

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"

	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

// Exercise for Advent of Code 2018 day 3.
type Exercise struct {
	common.BaseExercise
}

// claim is a rectangular fabric request: id at (left, top) spanning width x height.
type claim struct {
	id, left, top, width, height int
}

var claimRe = regexp.MustCompile(`-?\d+`)

// parse reads each claim by scanning its five integers (id, left, top, width,
// height), which is robust to whitespace and delivery quirks in the input.
func parse(instr string) ([]claim, error) {
	var claims []claim

	for _, line := range splitLines(instr) {
		nums := claimRe.FindAllString(line, -1)
		if len(nums) != 5 {
			return nil, fmt.Errorf("expected 5 numbers, got %d in %q", len(nums), line)
		}

		vals := make([]int, 5)
		for i, s := range nums {
			n, err := strconv.Atoi(s)
			if err != nil {
				return nil, fmt.Errorf("parsing %q: %w", s, err)
			}
			vals[i] = n
		}

		claims = append(claims, claim{vals[0], vals[1], vals[2], vals[3], vals[4]})
	}

	return claims, nil
}

// splitLines returns the non-empty lines of s.
func splitLines(s string) []string {
	var out []string
	start := -1
	for i, ch := range s {
		if ch == '\n' || ch == '\r' {
			if start >= 0 {
				out = append(out, s[start:i])
				start = -1
			}
			continue
		}
		if start < 0 {
			start = i
		}
	}
	if start >= 0 {
		out = append(out, s[start:])
	}
	return out
}

// coverage returns how many claims cover each fabric cell.
func coverage(claims []claim) map[[2]int]int {
	cover := make(map[[2]int]int)
	for _, c := range claims {
		for x := c.left; x < c.left+c.width; x++ {
			for y := c.top; y < c.top+c.height; y++ {
				cover[[2]int{x, y}]++
			}
		}
	}
	return cover
}

// One returns the answer to the first part of the exercise.
// answer: 116140
func (e Exercise) One(instr string) (any, error) {
	claims, err := parse(instr)
	if err != nil {
		return nil, err
	}

	overlap := 0
	for _, n := range coverage(claims) {
		if n >= 2 {
			overlap++
		}
	}

	return overlap, nil
}

// Two returns the answer to the second part of the exercise.
// answer: 574
func (e Exercise) Two(instr string) (any, error) {
	claims, err := parse(instr)
	if err != nil {
		return nil, err
	}

	cover := coverage(claims)

	for _, c := range claims {
		intact := true
		for x := c.left; x < c.left+c.width && intact; x++ {
			for y := c.top; y < c.top+c.height; y++ {
				if cover[[2]int{x, y}] > 1 {
					intact = false
					break
				}
			}
		}

		if intact {
			return c.id, nil
		}
	}

	return nil, errors.New("no non-overlapping claim found")
}
