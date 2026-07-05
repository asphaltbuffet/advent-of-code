package exercises

import (
	"regexp"
	"strconv"

	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

// Exercise for Advent of Code 2017 day 25.
type Exercise struct {
	common.BaseExercise
}

// action is what a state does for a given current tape value.
type action struct {
	write int
	move  int // +1 right, -1 left
	next  string
}

// rule holds the two actions (for tape value 0 and 1) of a state.
type rule struct{ on [2]action }

var (
	reBegin = regexp.MustCompile(`Begin in state (\w+)`)
	reSteps = regexp.MustCompile(`after (\d+) steps`)
	reState = regexp.MustCompile(`In state (\w+):`)
	reWrite = regexp.MustCompile(`Write the value (\d)`)
	reMove  = regexp.MustCompile(`Move one slot to the (right|left)`)
	reNext  = regexp.MustCompile(`Continue with state (\w+)`)
)

// parseBlueprint reads the machine's start state, step count, and per-state
// rules from the prose blueprint.
func parseBlueprint(instr string) (string, int, map[string]rule) {
	start := reBegin.FindStringSubmatch(instr)[1]
	steps, _ := strconv.Atoi(reSteps.FindStringSubmatch(instr)[1])

	rules := map[string]rule{}
	// Split into per-state blocks; the text before the first "In state" is the
	// header and is discarded.
	blocks := reState.Split(instr, -1)
	names := reState.FindAllStringSubmatch(instr, -1)
	for i, name := range names {
		block := blocks[i+1]
		writes := reWrite.FindAllStringSubmatch(block, -1)
		moves := reMove.FindAllStringSubmatch(block, -1)
		nexts := reNext.FindAllStringSubmatch(block, -1)

		var r rule
		for k := range 2 {
			w, _ := strconv.Atoi(writes[k][1])
			mv := 1
			if moves[k][1] == "left" {
				mv = -1
			}
			r.on[k] = action{write: w, move: mv, next: nexts[k][1]}
		}
		rules[name[1]] = r
	}
	return start, steps, rules
}

// One runs the Turing machine and returns the diagnostic checksum (count of 1s).
func (e Exercise) One(instr string) (any, error) {
	state, steps, rules := parseBlueprint(instr)

	tape := map[int]int{}
	cursor := 0
	for range steps {
		cur := tape[cursor]
		act := rules[state].on[cur]
		tape[cursor] = act.write
		cursor += act.move
		state = act.next
	}

	ones := 0
	for _, v := range tape {
		ones += v
	}
	return ones, nil
}

// Two has no puzzle: day 25 completes the year once every other star is earned.
func (e Exercise) Two(_ string) (any, error) {
	return "Merry Christmas!", nil
}
