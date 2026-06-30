package exercises

import (
	"strconv"
	"strings"

	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

// Exercise for Advent of Code 2024 day 21.
type Exercise struct {
	common.BaseExercise
}

type point struct{ r, c int }

// Keypad layouts as button -> position. The gap is excluded.
var numPad = map[byte]point{
	'7': {0, 0}, '8': {0, 1}, '9': {0, 2},
	'4': {1, 0}, '5': {1, 1}, '6': {1, 2},
	'1': {2, 0}, '2': {2, 1}, '3': {2, 2},
	'0': {3, 1}, 'A': {3, 2},
}
var numGap = point{3, 0}

var dirPad = map[byte]point{
	'^': {0, 1}, 'A': {0, 2},
	'<': {1, 0}, 'v': {1, 1}, '>': {1, 2},
}
var dirGap = point{0, 0}

// moveOptions returns the candidate move strings (each ending in 'A') that take
// the arm from `from` to `to` on the given pad without crossing the gap.
func moveOptions(from, to, gap point) []string {
	dr := to.r - from.r
	dc := to.c - from.c

	var vert, horiz strings.Builder
	for i := 0; i < abs(dr); i++ {
		if dr > 0 {
			vert.WriteByte('v')
		} else {
			vert.WriteByte('^')
		}
	}
	for i := 0; i < abs(dc); i++ {
		if dc > 0 {
			horiz.WriteByte('>')
		} else {
			horiz.WriteByte('<')
		}
	}
	v, h := vert.String(), horiz.String()

	// Two orderings: horizontal-first and vertical-first. Keep those that don't
	// pass over the gap.
	var opts []string
	// horizontal then vertical: passes through (from.r, to.c)
	if !(point{from.r, to.c} == gap) {
		opts = append(opts, h+v+"A")
	}
	// vertical then horizontal: passes through (to.r, from.c)
	if !(point{to.r, from.c} == gap) {
		opts = append(opts, v+h+"A")
	}
	// If both endpoints share a row or column, the two are identical; dedup.
	if len(opts) == 2 && opts[0] == opts[1] {
		opts = opts[:1]
	}
	return opts
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

type cacheKey struct {
	from, to byte
	depth    int
	numeric  bool
}

// solver memoizes the minimum number of outermost presses to move from->to and
// press it, with `depth` directional pads remaining above the current pad.
type solver struct {
	cache map[cacheKey]int
}

// cost returns the min outermost presses to type `to` after `from` on a pad.
// numeric selects the numeric vs directional pad layout; depth is how many
// directional pads sit between this pad and the human.
func (s *solver) cost(from, to byte, depth int, numeric bool) int {
	key := cacheKey{from, to, depth, numeric}
	if v, ok := s.cache[key]; ok {
		return v
	}

	pad, gap := dirPad, dirGap
	if numeric {
		pad, gap = numPad, numGap
	}
	opts := moveOptions(pad[from], pad[to], gap)

	best := -1
	for _, seq := range opts {
		var c int
		if depth == 0 {
			c = len(seq) // human presses these directly
		} else {
			// Type `seq` on the next directional pad up, starting from 'A'.
			prev := byte('A')
			for i := 0; i < len(seq); i++ {
				c += s.cost(prev, seq[i], depth-1, false)
				prev = seq[i]
			}
		}
		if best == -1 || c < best {
			best = c
		}
	}
	s.cache[key] = best
	return best
}

// typeCost returns the total outermost presses to type the whole code on the
// numeric pad through `dirPads` directional pads.
func (s *solver) typeCost(code string, dirPads int) int {
	total := 0
	prev := byte('A')
	for i := 0; i < len(code); i++ {
		total += s.cost(prev, code[i], dirPads, true)
		prev = code[i]
	}
	return total
}

func sumComplexities(instr string, dirPads int) int {
	s := &solver{cache: map[cacheKey]int{}}
	sum := 0
	for _, code := range strings.Fields(instr) {
		length := s.typeCost(code, dirPads)
		numeric, _ := strconv.Atoi(strings.TrimRight(strings.TrimLeft(code, "0"), "A"))
		sum += length * numeric
	}
	return sum
}

// One returns the answer to the first part of the exercise.
func (e Exercise) One(instr string) (any, error) {
	return sumComplexities(instr, 2), nil
}

// Two returns the answer to the second part of the exercise.
func (e Exercise) Two(instr string) (any, error) {
	return sumComplexities(instr, 25), nil
}
