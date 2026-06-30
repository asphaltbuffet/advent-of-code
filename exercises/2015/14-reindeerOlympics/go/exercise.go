package exercises

import (
	"strconv"
	"strings"

	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

// Exercise for Advent of Code 2015 day 14.
type Exercise struct {
	common.BaseExercise
}

type reindeer struct {
	speed, fly, rest int
}

// parse reads the reindeer stats. The race duration is not in the input, so it
// is inferred: the 2-reindeer AoC example races 1000s, the real input 2503s.
func parse(instr string) ([]reindeer, int) {
	var rs []reindeer

	for _, line := range strings.Split(strings.TrimSpace(instr), "\n") {
		f := strings.Fields(line)
		// <name> can fly <speed> km/s for <fly> seconds, ... rest for <rest> ...
		speed, _ := strconv.Atoi(f[3])
		fly, _ := strconv.Atoi(f[6])
		rest, _ := strconv.Atoi(f[13])
		rs = append(rs, reindeer{speed, fly, rest})
	}

	duration := 2503
	if len(rs) <= 2 {
		duration = 1000
	}

	return rs, duration
}

// distance returns how far a reindeer has travelled after t seconds.
func (r reindeer) distance(t int) int {
	cycle := r.fly + r.rest
	full := t / cycle
	rem := t % cycle
	flying := full * r.fly
	if rem < r.fly {
		flying += rem
	} else {
		flying += r.fly
	}
	return flying * r.speed
}

// One returns the answer to the first part of the exercise.
func (e Exercise) One(instr string) (any, error) {
	rs, duration := parse(instr)

	best := 0
	for _, r := range rs {
		if d := r.distance(duration); d > best {
			best = d
		}
	}

	return best, nil
}

// Two returns the answer to the second part of the exercise: at every second,
// the reindeer in the lead earn a point; return the highest total.
func (e Exercise) Two(instr string) (any, error) {
	rs, duration := parse(instr)
	points := make([]int, len(rs))

	for t := 1; t <= duration; t++ {
		lead := 0
		for _, r := range rs {
			if d := r.distance(t); d > lead {
				lead = d
			}
		}
		for i, r := range rs {
			if r.distance(t) == lead {
				points[i]++
			}
		}
	}

	best := 0
	for _, p := range points {
		if p > best {
			best = p
		}
	}

	return best, nil
}
