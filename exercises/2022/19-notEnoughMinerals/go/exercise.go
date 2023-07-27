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
		geodesMade := st.calcMostGeodes(0, map[string]int{}, 24, 24)
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
			geodesMade := ns.calcMostGeodes(0, map[string]int{}, 32, 32)
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

func (s *state) hash(time int) string {
	return fmt.Sprintf("%d%d%d%d%d%d%d%d%d",
		time,
		s.ore,
		s.clay,
		s.obsidian,
		s.geode,
		s.oreRobots,
		s.clayRobots,
		s.obsidianRobots,
		s.geodeRobots)
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

func (s *state) calcMostGeodes(time int, seenStates map[string]int, totalTime int, earliestGeode int) int {
	if time == totalTime {
		return s.geode
	}

	h := s.hash(time)
	if v, ok := seenStates[h]; ok {
		return v
	}

	if s.geode == 0 && time > earliestGeode {
		return 0
	}

	mostGeodes := s.geode

	if s.ore >= s.oreForGeodeRobot && s.obsidian >= s.obsidianForGeodeRobot {
		// make geode robots
		cp := s.copy()

		cp.farm()
		cp.makeGeodeRobot()

		if cp.geodeRobots == 1 {
			if time+1 < earliestGeode {
				earliestGeode = time + 1
			}
		}

		curMostGeodes := cp.calcMostGeodes(time+1, seenStates, totalTime, earliestGeode)
		if curMostGeodes > mostGeodes {
			mostGeodes = curMostGeodes
		}

		seenStates[h] = mostGeodes

		return mostGeodes
	}

	if time <= totalTime-16 && s.oreRobots < s.oreForObsidianRobot*2 && s.ore >= s.oreForOreRobot {
		// make ore robot
		cp := s.copy()
		cp.farm()
		cp.makeOreRobot()

		curMostGeodes := cp.calcMostGeodes(time+1, seenStates, totalTime, earliestGeode)
		if curMostGeodes > mostGeodes {
			mostGeodes = curMostGeodes
		}
	}

	if time <= totalTime-8 && s.clayRobots < s.clayForObsidianRobot && s.ore >= s.oreForClayRobot {
		// make clay robot
		cp := s.copy()
		cp.farm()
		cp.makeClayRobot()

		curMostGeodes := cp.calcMostGeodes(time+1, seenStates, totalTime, earliestGeode)
		if curMostGeodes > mostGeodes {
			mostGeodes = curMostGeodes
		}
	}

	if time <= totalTime-4 && s.obsidianRobots < s.obsidianForGeodeRobot && s.ore >= s.oreForObsidianRobot && s.clay >= s.clayForObsidianRobot {
		// make obsidian robot
		cp := s.copy()
		cp.farm()
		cp.makeObsidianRobot()

		curMostGeodes := cp.calcMostGeodes(time+1, seenStates, totalTime, earliestGeode)
		if curMostGeodes > mostGeodes {
			mostGeodes = curMostGeodes
		}
	}

	// or no factory production this minute
	cp := s.copy()
	cp.farm()

	curMostGeodes := cp.calcMostGeodes(time+1, seenStates, totalTime, earliestGeode)
	if curMostGeodes > mostGeodes {
		mostGeodes = curMostGeodes
	}

	seenStates[h] = mostGeodes

	return mostGeodes
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
