package exercises

import (
	"strconv"
	"strings"

	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

// Exercise for Advent of Code 2015 day 13.
type Exercise struct {
	common.BaseExercise
}

// parse builds the happiness map h[a][b] = change in a's happiness when seated
// next to b, and returns the sorted list of distinct people.
func parse(instr string) (map[string]map[string]int, []string) {
	h := map[string]map[string]int{}

	for _, line := range strings.Split(strings.TrimSpace(instr), "\n") {
		f := strings.Fields(strings.TrimSuffix(line, "."))
		// f: <A> would <gain|lose> <N> happiness units by sitting next to <B>
		a, b := f[0], f[10]
		n, _ := strconv.Atoi(f[3])
		if f[2] == "lose" {
			n = -n
		}
		if h[a] == nil {
			h[a] = map[string]int{}
		}
		h[a][b] = n
	}

	people := make([]string, 0, len(h))
	for p := range h {
		people = append(people, p)
	}

	return h, people
}

// bestHappiness returns the maximum total happiness over all circular seatings.
// The first person is fixed to factor out rotations (n-1)! arrangements.
func bestHappiness(h map[string]map[string]int, people []string) int {
	best := 0
	first := true

	rest := append([]string(nil), people[1:]...)
	permute(rest, 0, func(order []string) {
		seating := append([]string{people[0]}, order...)
		total := 0
		n := len(seating)
		for i := 0; i < n; i++ {
			a := seating[i]
			b := seating[(i+1)%n]
			total += h[a][b] + h[b][a]
		}
		if first || total > best {
			best = total
			first = false
		}
	})

	return best
}

// permute generates all permutations of s in place, calling fn for each.
func permute(s []string, k int, fn func([]string)) {
	if k == len(s) {
		fn(s)
		return
	}
	for i := k; i < len(s); i++ {
		s[k], s[i] = s[i], s[k]
		permute(s, k+1, fn)
		s[k], s[i] = s[i], s[k]
	}
}

// One returns the answer to the first part of the exercise.
func (e Exercise) One(instr string) (any, error) {
	h, people := parse(instr)
	return bestHappiness(h, people), nil
}

// Two returns the answer to the second part of the exercise: add yourself, with
// zero happiness in both directions to everyone, then re-optimize.
func (e Exercise) Two(instr string) (any, error) {
	h, people := parse(instr)

	h["me"] = map[string]int{}
	for _, p := range people {
		h["me"][p] = 0
		h[p]["me"] = 0
	}
	people = append(people, "me")

	return bestHappiness(h, people), nil
}
