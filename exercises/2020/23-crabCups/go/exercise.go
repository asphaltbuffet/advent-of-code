package exercises

import (
	"fmt"
	"strings"

	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

// Exercise for Advent of Code 2020 day 23.
type Exercise struct {
	common.BaseExercise
}

func parseCups(instr string) []int {
	s := strings.TrimSpace(instr)
	cups := make([]int, len(s))
	for i, c := range s {
		cups[i] = int(c - '0')
	}
	return cups
}

// play runs the crab game over `moves` rounds on `total` cups. cups holds the
// starting labels; if total exceeds len(cups) the cups are extended with the
// increasing labels up to total. Returns the "next" array: next[label] is the cup
// clockwise after `label`.
func play(cups []int, total, moves int) []int {
	// next[label] -> following label. Index 0 is unused.
	next := make([]int, total+1)
	for i := 0; i < len(cups)-1; i++ {
		next[cups[i]] = cups[i+1]
	}
	if total > len(cups) {
		// chain the last input cup into 10..total, then wrap to the first cup.
		next[cups[len(cups)-1]] = len(cups) + 1
		for c := len(cups) + 1; c < total; c++ {
			next[c] = c + 1
		}
		next[total] = cups[0]
	} else {
		next[cups[len(cups)-1]] = cups[0]
	}

	current := cups[0]
	for m := 0; m < moves; m++ {
		// Pick up the three cups after current.
		a := next[current]
		b := next[a]
		c := next[b]

		// Destination: current-1, wrapping, skipping picked-up.
		dest := current - 1
		if dest < 1 {
			dest = total
		}
		for dest == a || dest == b || dest == c {
			dest--
			if dest < 1 {
				dest = total
			}
		}

		// Splice the three out and in after dest.
		next[current] = next[c]
		next[c] = next[dest]
		next[dest] = a

		current = next[current]
	}

	return next
}

// One returns the labels after cup 1 following 100 moves.
func (e Exercise) One(instr string) (any, error) {
	cups := parseCups(instr)
	next := play(cups, len(cups), 100)

	var b strings.Builder
	for c := next[1]; c != 1; c = next[c] {
		fmt.Fprintf(&b, "%d", c)
	}
	return b.String(), nil
}

// Two extends to one million cups, runs ten million moves, and returns the
// product of the two cups immediately clockwise of cup 1.
func (e Exercise) Two(instr string) (any, error) {
	cups := parseCups(instr)
	next := play(cups, 1_000_000, 10_000_000)

	a := next[1]
	b := next[a]
	return fmt.Sprintf("%d", a*b), nil
}
