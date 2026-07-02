package exercises

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

// Exercise for Advent of Code 2020 day 13.
type Exercise struct {
	common.BaseExercise
}

// bus pairs a bus ID with its position (offset) in the schedule line.
type bus struct {
	id, offset int
}

// parse returns the earliest departure timestamp and the in-service buses.
func parse(instr string) (int, []bus, error) {
	lines := strings.Split(strings.TrimSpace(instr), "\n")
	if len(lines) < 2 {
		return 0, nil, fmt.Errorf("expected two lines")
	}
	earliest, err := strconv.Atoi(strings.TrimSpace(lines[0]))
	if err != nil {
		return 0, nil, fmt.Errorf("parsing timestamp: %w", err)
	}

	var buses []bus
	for i, tok := range strings.Split(strings.TrimSpace(lines[1]), ",") {
		if tok == "x" {
			continue
		}
		id, err := strconv.Atoi(tok)
		if err != nil {
			return 0, nil, fmt.Errorf("parsing bus id %q: %w", tok, err)
		}
		buses = append(buses, bus{id, i})
	}
	return earliest, buses, nil
}

// One finds the earliest bus departing at or after the timestamp and returns its
// ID times the minutes waited.
func (e Exercise) One(instr string) (any, error) {
	earliest, buses, err := parse(instr)
	if err != nil {
		return nil, err
	}

	bestWait, bestID := -1, 0
	for _, b := range buses {
		// Minutes until this bus's next departure at or after `earliest`.
		wait := (b.id - earliest%b.id) % b.id
		if bestWait == -1 || wait < bestWait {
			bestWait, bestID = wait, b.id
		}
	}

	return fmt.Sprintf("%d", bestID*bestWait), nil
}

// Two finds the earliest timestamp t at which every bus id departs at t+offset,
// via an incremental sieve: lock in each bus by stepping t forward, then grow the
// step by that bus's id (the ids are pairwise coprime, so earlier constraints are
// preserved).
func (e Exercise) Two(instr string) (any, error) {
	_, buses, err := parse(instr)
	if err != nil {
		return nil, err
	}

	t, step := 0, 1
	for _, b := range buses {
		for (t+b.offset)%b.id != 0 {
			t += step
		}
		step *= b.id
	}

	return fmt.Sprintf("%d", t), nil
}
