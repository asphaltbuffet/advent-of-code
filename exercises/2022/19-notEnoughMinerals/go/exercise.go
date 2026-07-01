package exercises

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"sync"

	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

// Exercise for Advent of Code 2022 day 19.
type Exercise struct {
	common.BaseExercise
}

// One returns the answer to the first part of the exercise.
// answer:
func (c Exercise) One(instr string) (any, error) {
	blueprints, err := parse(instr)
	if err != nil {
		return nil, fmt.Errorf("parsing input: %w", err)
	}

	// how many geodes can be opened in 24 minutes?
	sum := 0

	for _, bp := range blueprints {
		st := newState(*bp)
		best := 0
		geodesMade := st.calcMostGeodes(0, 24, &best)
		sum += st.blueprint.id * geodesMade
	}

	// total quality of all blueprints, quality = id * (# geodes in 24 min)
	return sum, nil
}

// Two returns the answer to the second part of the exercise.
// wrong: 122023936
// answer:
func (c Exercise) Two(instr string) (any, error) {
	blueprints, err := parse(instr)
	if err != nil {
		return nil, fmt.Errorf("parsing input: %w", err)
	}

	if len(blueprints) > 3 {
		blueprints = blueprints[:3]
	}

	numGeodes := 1

	var wg sync.WaitGroup

	resultChan := make(chan int)
	errChan := make(chan error)

	for _, bp := range blueprints {
		wg.Add(1)

		go func(bp *blueprint) {
			defer wg.Done()

			ns := newState(*bp)
			best := 0
			geodesMade := ns.calcMostGeodes(0, 32, &best)
			// fmt.Printf("id=%d, made=%d", bp.id, geodesMade)
			resultChan <- geodesMade
		}(bp)
	}

	go func() {
		wg.Wait()
		close(resultChan)
		close(errChan)
	}()

	for geodesMade := range resultChan {
		numGeodes *= geodesMade
	}

	if err, ok := <-errChan; ok {
		return nil, fmt.Errorf("error calculating geodes: %w", err)
	}

	// total quality of all blueprints, quality = id * (# geodes in 24 min)
	return numGeodes, nil
}

type blueprint struct {
	id                                        int
	oreForOreRobot                            int
	oreForClayRobot                           int
	oreForObsidianRobot, clayForObsidianRobot int
	oreForGeodeRobot, obsidianForGeodeRobot   int
}

type state struct {
	blueprint
	ore, clay, obsidian, geode                         int
	oreRobots, clayRobots, obsidianRobots, geodeRobots int
}

func newState(blueprint blueprint) state {
	return state{
		blueprint: blueprint,
		oreRobots: 1,
	}
}

func (s *state) copy() state {
	return state{
		blueprint:      s.blueprint,
		ore:            s.ore,
		clay:           s.clay,
		obsidian:       s.obsidian,
		geode:          s.geode,
		oreRobots:      s.oreRobots,
		clayRobots:     s.clayRobots,
		obsidianRobots: s.obsidianRobots,
		geodeRobots:    s.geodeRobots,
	}
}

// calcMostGeodes runs a branch-and-bound DFS over "which robot to build next".
// Rather than branching every minute, each branch fast-forwards to the minute a
// chosen robot becomes affordable and builds it — far fewer nodes. best is the
// running maximum shared across the search so an optimistic bound can prune.
func (s *state) calcMostGeodes(time, totalTime int, best *int) int {
	remaining := totalTime - time

	// If we build nothing more, current geode robots run out the clock.
	finalGeodes := s.geode + s.geodeRobots*remaining
	if finalGeodes > *best {
		*best = finalGeodes
	}

	// Optimistic bound: if we could build a geode robot every remaining minute,
	// how many geodes could we finish with? If that can't beat best, abandon.
	upper := s.geode + s.geodeRobots*remaining + remaining*(remaining-1)/2
	if upper <= *best {
		return *best
	}

	// Building more of a robot than the max per-minute demand for its output
	// never helps, so cap each type at the highest consumer's cost.
	maxOre := maxInt(maxInt(s.oreForOreRobot, s.oreForClayRobot), maxInt(s.oreForObsidianRobot, s.oreForGeodeRobot))

	mostGeodes := finalGeodes

	try := func(build func(*state), need func(*state) bool, cost func(*state) (int, int, int)) {
		if !need(s) {
			return
		}
		// Minutes to wait until this robot is affordable given current robots.
		oreC, clayC, obsC := cost(s)
		wait := 0
		wait = maxInt(wait, ceilDiv(oreC-s.ore, s.oreRobots))
		if clayC > 0 {
			if s.clayRobots == 0 {
				return
			}
			wait = maxInt(wait, ceilDiv(clayC-s.clay, s.clayRobots))
		}
		if obsC > 0 {
			if s.obsidianRobots == 0 {
				return
			}
			wait = maxInt(wait, ceilDiv(obsC-s.obsidian, s.obsidianRobots))
		}
		build2 := wait + 1 // minutes until robot is producing
		if build2 >= remaining {
			return // no time left to benefit from this robot
		}
		cp := s.copy()
		for i := 0; i < build2; i++ {
			cp.farm()
		}
		build(&cp)
		mostGeodes = maxInt(mostGeodes, cp.calcMostGeodes(time+build2, totalTime, best))
	}

	// Geode first so best tightens quickly for the bound above.
	try((*state).makeGeodeRobot,
		func(*state) bool { return true },
		func(s *state) (int, int, int) { return s.oreForGeodeRobot, 0, s.obsidianForGeodeRobot })
	try((*state).makeObsidianRobot,
		func(s *state) bool { return s.obsidianRobots < s.obsidianForGeodeRobot },
		func(s *state) (int, int, int) { return s.oreForObsidianRobot, s.clayForObsidianRobot, 0 })
	try((*state).makeClayRobot,
		func(s *state) bool { return s.clayRobots < s.clayForObsidianRobot },
		func(s *state) (int, int, int) { return s.oreForClayRobot, 0, 0 })
	try((*state).makeOreRobot,
		func(s *state) bool { return s.oreRobots < maxOre },
		func(s *state) (int, int, int) { return s.oreForOreRobot, 0, 0 })

	return mostGeodes
}

// ceilDiv returns ceil(a/b) for b > 0, clamped to 0 for non-positive a.
func ceilDiv(a, b int) int {
	if a <= 0 {
		return 0
	}
	return (a + b - 1) / b
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func parse(input string) ([]*blueprint, error) {
	ans := make([]*blueprint, strings.Count(input, "\n")+1)

	for i, line := range strings.Split(input, "\n") {
		bp := blueprint{}

		_, err := fmt.Sscanf(line, "Blueprint %d: Each ore robot costs %d ore. Each clay robot costs %d ore. Each obsidian robot costs %d ore and %d clay. Each geode robot costs %d ore and %d obsidian.",
			&bp.id, &bp.oreForOreRobot, &bp.oreForClayRobot, &bp.oreForObsidianRobot,
			&bp.clayForObsidianRobot, &bp.oreForGeodeRobot, &bp.obsidianForGeodeRobot)
		if err != nil {
			return nil, fmt.Errorf("parsing input: %w", err)
		}

		ans[i] = &bp
	}

	return ans, nil
}

func parseWithRegex(input string) ([]*blueprint, error) {
	ans := make([]*blueprint, strings.Count(input, "\n")+1)

	for i, line := range strings.Split(input, "\n") {
		bp := blueprint{}

		pattern := `Blueprint (\d+): Each ore robot costs (\d+) ore. Each clay robot costs (\d+) ore. Each obsidian robot costs (\d+) ore and (\d+) clay. Each geode robot costs (\d+) ore and (\d+) obsidian.`
		regex := regexp.MustCompile(pattern)
		match := regex.FindStringSubmatch(line)

		if len(match) != 8 { // We expect 8 items: the full match, and 7 capture groups.
			return nil, fmt.Errorf("parsing input: incorrect format")
		}

		var err error

		bp.id, err = strconv.Atoi(match[1])
		if err != nil {
			return nil, fmt.Errorf("parsing input: %w", err)
		}

		bp.oreForOreRobot, err = strconv.Atoi(match[2])
		if err != nil {
			return nil, fmt.Errorf("parsing input: %w", err)
		}

		bp.oreForClayRobot, err = strconv.Atoi(match[3])
		if err != nil {
			return nil, fmt.Errorf("parsing input: %w", err)
		}

		bp.oreForObsidianRobot, err = strconv.Atoi(match[4])
		if err != nil {
			return nil, fmt.Errorf("parsing input: %w", err)
		}

		bp.clayForObsidianRobot, err = strconv.Atoi(match[5])
		if err != nil {
			return nil, fmt.Errorf("parsing input: %w", err)
		}

		bp.oreForGeodeRobot, err = strconv.Atoi(match[6])
		if err != nil {
			return nil, fmt.Errorf("parsing input: %w", err)
		}

		bp.obsidianForGeodeRobot, err = strconv.Atoi(match[7])
		if err != nil {
			return nil, fmt.Errorf("parsing input: %w", err)
		}

		ans[i] = &bp
	}

	return ans, nil
}
