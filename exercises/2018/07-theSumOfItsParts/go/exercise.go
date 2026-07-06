package exercises

import (
	"fmt"
	"regexp"
	"slices"
	"strings"

	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

// Exercise for Advent of Code 2018 day 7.
type Exercise struct {
	common.BaseExercise
}

var stepRe = regexp.MustCompile(`Step (\w) must be finished before step (\w) can begin\.`)

// edge captures a "before must precede after" ordering constraint.
type edge struct{ before, after byte }

// parse reads each dependency by scanning its two step letters.
func parse(instr string) ([]edge, error) {
	var edges []edge

	for line := range strings.SplitSeq(strings.TrimSpace(instr), "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		m := stepRe.FindStringSubmatch(line)
		if m == nil {
			return nil, fmt.Errorf("unrecognized line %q", line)
		}

		edges = append(edges, edge{before: m[1][0], after: m[2][0]})
	}

	return edges, nil
}

// deps builds, for every step, the set of steps that must precede it.
func deps(edges []edge) map[byte]map[byte]bool {
	d := make(map[byte]map[byte]bool)

	for _, e := range edges {
		if d[e.before] == nil {
			d[e.before] = map[byte]bool{}
		}
		if d[e.after] == nil {
			d[e.after] = map[byte]bool{}
		}
		d[e.after][e.before] = true
	}

	return d
}

func readySteps(d map[byte]map[byte]bool, done map[byte]bool, inProgress map[byte]int) []byte {
	var ready []byte
	for step, prereqs := range d {
		if done[step] || inProgress[step] != 0 {
			continue
		}
		blocked := false
		for p := range prereqs {
			if !done[p] {
				blocked = true
				break
			}
		}
		if !blocked {
			ready = append(ready, step)
		}
	}
	slices.Sort(ready)
	return ready
}

// One returns the answer to the first part of the exercise.
// Answer: CQSWKZFJONPBEUMXADLYIGVRHT
func (e Exercise) One(instr string) (any, error) {
	edges, err := parse(instr)
	if err != nil {
		return nil, err
	}

	d := deps(edges)

	done := map[byte]bool{}
	var order strings.Builder

	for len(done) < len(d) {
		// Among steps not yet done whose prerequisites are all satisfied, take
		// the alphabetically first.
		var ready []byte
		for step, prereqs := range d {
			if done[step] {
				continue
			}

			blocked := false
			for p := range prereqs {
				if !done[p] {
					blocked = true
					break
				}
			}

			if !blocked {
				ready = append(ready, step)
			}
		}

		slices.Sort(ready)

		next := ready[0]
		done[next] = true
		order.WriteByte(next)
	}

	return order.String(), nil
}

// Two returns the answer to the second part of the exercise.
// Answer: 914
func (e Exercise) Two(instr string) (any, error) {
	edges, err := parse(instr)
	if err != nil {
		return nil, err
	}

	d := deps(edges)

	workers, base := 5, 60
	if len(d) <= 6 {
		workers, base = 2, 0
	}

	done := map[byte]bool{}
	inProgress := map[byte]int{} // step -> second it finishes

	for t := 0; ; t++ {
		for step, finish := range inProgress {
			if finish <= t {
				done[step] = true
				delete(inProgress, step)
			}
		}

		if len(done) == len(d) {
			return t, nil
		}

		for _, step := range readySteps(d, done, inProgress) {
			if len(inProgress) >= workers {
				break
			}
			inProgress[step] = t + base + int(step-'A'+1)
		}
	}
}
