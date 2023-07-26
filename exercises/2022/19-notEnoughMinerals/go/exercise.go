package exercises

import (
	"fmt"
	"strings"

	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

// Exercise for Advent of Code 2022 day 19.
type Exercise struct {
	common.BaseExercise
}

// One returns the answer to the first part of the exercise.
// answer:
func (c Exercise) One(instr string) (any, error) {
	blueprints, err := parseInput(instr)
	if err != nil {
		return nil, fmt.Errorf("parsing input: %w", err)
	}

	// how many geodes can be opened in 24 minutes?
	sum := 0
	for _, bp := range blueprints {
		st := newState(bp)
		geodesMade := st.calcMostGeodes(0, map[string]int{}, 24, 24)
		// fmt.Println("ID:", bp.id, geodesMade)
		sum += st.blueprint.id * geodesMade
	}

	// total quality of all blueprints, quality = id * (# geodes in 24 min)
	return sum, nil
}

// Two returns the answer to the second part of the exercise.
// wrong: 122023936
// answer:
func (c Exercise) Two(instr string) (any, error) {
	blueprints, err := parseInput(instr)
	if err != nil {
		return nil, fmt.Errorf("parsing input: %w", err)
	}

	if len(blueprints) > 3 {
		blueprints = blueprints[:3]
	}

	prod := 1

	for _, bp := range blueprints {
		st := newState(bp)
		geodesMade := st.calcMostGeodes(0, map[string]int{}, 32, 32)
		// fmt.Println(bp.id, geodesMade)
		prod *= geodesMade
	}

	// total quality of all blueprints, quality = id * (# geodes in 24 min)
	return prod, nil
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

func (s *state) farm() {
	s.ore += s.oreRobots
	s.clay += s.clayRobots
	s.obsidian += s.obsidianRobots
	s.geode += s.geodeRobots
}

func (s *state) hash(time int) string {
	return fmt.Sprint(time, s.ore, s.clay, s.obsidian,
		s.geode, s.oreRobots, s.clayRobots, s.obsidianRobots, s.geodeRobots)
}

// NOT A POINTER METHOD SO A COPY CAN BE MADE
// this is some cheeky Go struct copying, it'd be easier to read if it was just
// directly recreating all the fields
func (s state) copy() state {
	return s
}

func (s *state) calcMostGeodes(time int, memo map[string]int, totalTime int, earliestGeode int) int {
	if time == totalTime {
		return s.geode
	}

	h := s.hash(time)
	if v, ok := memo[h]; ok {
		return v
	}

	if s.geode == 0 && time > earliestGeode {
		return 0
	}

	// factory can try to make any possible robot, will backtrack if necessary
	mostGeodes := s.geode

	// always make geode robots
	if s.ore >= s.oreForGeodeRobot &&
		s.obsidian >= s.obsidianForGeodeRobot {
		cp := s.copy()

		cp.farm()

		cp.ore -= cp.oreForGeodeRobot
		cp.obsidian -= cp.obsidianForGeodeRobot
		cp.geodeRobots++
		if cp.geodeRobots == 1 {
			earliestGeode = minInt(earliestGeode, time+1)
		}
		mostGeodes = maxInt(mostGeodes, cp.calcMostGeodes(time+1, memo, totalTime, earliestGeode))

		memo[h] = mostGeodes
		return mostGeodes
	}

	if time <= totalTime-16 &&
		s.oreRobots < s.oreForObsidianRobot*2 &&
		s.ore >= s.oreForOreRobot {
		cp := s.copy()
		cp.ore -= cp.oreForOreRobot

		cp.farm()

		cp.oreRobots++
		mostGeodes = maxInt(mostGeodes, cp.calcMostGeodes(time+1, memo, totalTime, earliestGeode))
	}
	if time <= totalTime-8 &&
		s.clayRobots < s.clayForObsidianRobot &&
		s.ore >= s.oreForClayRobot {
		cp := s.copy()
		cp.ore -= cp.oreForClayRobot

		cp.farm()

		cp.clayRobots++
		mostGeodes = maxInt(mostGeodes, cp.calcMostGeodes(time+1, memo, totalTime, earliestGeode))
	}
	if time <= totalTime-4 &&
		s.obsidianRobots < s.obsidianForGeodeRobot &&
		s.ore >= s.oreForObsidianRobot && s.clay >= s.clayForObsidianRobot {

		cp := s.copy()
		cp.ore -= cp.oreForObsidianRobot
		cp.clay -= cp.clayForObsidianRobot
		cp.farm()

		cp.obsidianRobots++
		mostGeodes = maxInt(mostGeodes, cp.calcMostGeodes(time+1, memo, totalTime, earliestGeode))
	}

	// or no factory production this minute
	cp := s.copy()
	cp.ore += cp.oreRobots
	cp.clay += cp.clayRobots
	cp.obsidian += cp.obsidianRobots
	cp.geode += cp.geodeRobots
	mostGeodes = maxInt(mostGeodes, cp.calcMostGeodes(time+1, memo, totalTime, earliestGeode))

	memo[h] = mostGeodes
	return mostGeodes
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func minInt(a, b int) int {
	if a > b {
		return b
	}
	return a
}

func parseInput(input string) ([]blueprint, error) {
	// Blueprint 1: Each ore robot costs 3 ore. Each clay robot costs 3 ore. Each obsidian robot costs 3 ore and 20 clay. Each geode robot costs 2 ore and 12 obsidian.
	ans := make([]blueprint, strings.Count(input, "\n"))
	for _, line := range strings.Split(input, "\n") {
		bp := blueprint{}
		_, err := fmt.Sscanf(line, "Blueprint %d: Each ore robot costs %d ore. Each clay robot costs %d ore. Each obsidian robot costs %d ore and %d clay. Each geode robot costs %d ore and %d obsidian.",
			&bp.id, &bp.oreForOreRobot, &bp.oreForClayRobot, &bp.oreForObsidianRobot,
			&bp.clayForObsidianRobot, &bp.oreForGeodeRobot, &bp.obsidianForGeodeRobot)
		if err != nil {
			return nil, fmt.Errorf("parsing input: %w", err)
		}
		ans = append(ans, bp)
	}
	return ans, nil
}
