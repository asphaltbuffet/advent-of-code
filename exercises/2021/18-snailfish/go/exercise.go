package exercises

import (
	"fmt"
	"strings"

	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

// Exercise for Advent of Code 2021 day 18.
type Exercise struct {
	common.BaseExercise
}

// token is one regular number together with its bracket nesting depth. A whole
// snailfish number is a flat list of tokens in reading order, which makes the
// explode/split left/right-neighbor rules simple index arithmetic.
type token struct {
	value, depth int
}

func parseLine(line string) []token {
	var toks []token
	depth := 0
	i := 0
	for i < len(line) {
		switch c := line[i]; c {
		case '[':
			depth++
			i++
		case ']':
			depth--
			i++
		case ',':
			i++
		default:
			// read a (possibly multi-digit) number
			j := i
			for j < len(line) && line[j] >= '0' && line[j] <= '9' {
				j++
			}
			n := 0
			for _, d := range line[i:j] {
				n = n*10 + int(d-'0')
			}
			toks = append(toks, token{n, depth})
			i = j
		}
	}
	return toks
}

func parse(instr string) [][]token {
	var nums [][]token
	for _, line := range strings.Split(strings.TrimSpace(instr), "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		nums = append(nums, parseLine(line))
	}
	return nums
}

// clone copies a token list so additions do not mutate the parsed inputs.
func clone(t []token) []token {
	out := make([]token, len(t))
	copy(out, t)
	return out
}

// add concatenates two numbers into a pair (all depths +1) and reduces.
func add(a, b []token) []token {
	sum := make([]token, 0, len(a)+len(b))
	sum = append(sum, clone(a)...)
	sum = append(sum, clone(b)...)
	for i := range sum {
		sum[i].depth++
	}
	return reduce(sum)
}

func reduce(t []token) []token {
	for {
		if next, ok := explode(t); ok {
			t = next
			continue
		}
		if next, ok := split(t); ok {
			t = next
			continue
		}
		return t
	}
}

// explode collapses the leftmost pair nested at depth 5 (two adjacent tokens at
// depth 5): its left adds to the previous token, its right to the next, and the
// pair becomes a single 0 at depth 4.
func explode(t []token) ([]token, bool) {
	for i := 0; i+1 < len(t); i++ {
		if t[i].depth >= 5 && t[i].depth == t[i+1].depth {
			if i > 0 {
				t[i-1].value += t[i].value
			}
			if i+2 < len(t) {
				t[i+2].value += t[i+1].value
			}
			out := make([]token, 0, len(t)-1)
			out = append(out, t[:i]...)
			out = append(out, token{0, t[i].depth - 1})
			out = append(out, t[i+2:]...)
			return out, true
		}
	}
	return t, false
}

// split replaces the leftmost value >= 10 with a pair (floor, ceil) one deeper.
func split(t []token) ([]token, bool) {
	for i := range t {
		if t[i].value >= 10 {
			left := t[i].value / 2
			right := t[i].value - left
			d := t[i].depth + 1
			out := make([]token, 0, len(t)+1)
			out = append(out, t[:i]...)
			out = append(out, token{left, d}, token{right, d})
			out = append(out, t[i+1:]...)
			return out, true
		}
	}
	return t, false
}

// magnitude repeatedly combines the deepest adjacent pair as 3*left + 2*right
// until a single value remains.
func magnitude(t []token) int {
	t = clone(t)
	for len(t) > 1 {
		maxDepth, at := 0, -1
		for i := 0; i+1 < len(t); i++ {
			if t[i].depth == t[i+1].depth && t[i].depth > maxDepth {
				maxDepth, at = t[i].depth, i
			}
		}
		combined := token{3*t[at].value + 2*t[at+1].value, t[at].depth - 1}
		out := make([]token, 0, len(t)-1)
		out = append(out, t[:at]...)
		out = append(out, combined)
		out = append(out, t[at+2:]...)
		t = out
	}
	return t[0].value
}

// One returns the answer to the first part of the exercise.
func (e Exercise) One(instr string) (any, error) {
	nums := parse(instr)
	if len(nums) == 0 {
		return nil, fmt.Errorf("no snailfish numbers in input")
	}

	acc := clone(nums[0])
	for _, n := range nums[1:] {
		acc = add(acc, n)
	}

	return fmt.Sprintf("%d", magnitude(acc)), nil
}

// Two returns the answer to the second part of the exercise.
func (e Exercise) Two(instr string) (any, error) {
	nums := parse(instr)
	if len(nums) < 2 {
		return nil, fmt.Errorf("need at least two numbers")
	}

	best := 0
	for i := range nums {
		for j := range nums {
			if i == j {
				continue
			}
			if m := magnitude(add(nums[i], nums[j])); m > best {
				best = m
			}
		}
	}

	return fmt.Sprintf("%d", best), nil
}
