package exercises

import (
	"strconv"
	"strings"

	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

// Exercise for Advent of Code 2016 day 10.
type Exercise struct {
	common.BaseExercise
}

// dest is a destination for a chip: a bot or an output bin.
type dest struct {
	isOutput bool
	id       int
}

// rule describes where a bot sends its low and high chips.
type rule struct {
	low, high dest
}

// factory is the simulated state after running all instructions.
type factory struct {
	bots    map[int][]int // bot id -> chips held (resolved)
	outputs map[int]int   // output bin -> chip value
	// compareBot is the bot that compared the two target values.
	compareBot int
}

// simulate runs the factory, recording which bot compares loA and loB.
func simulate(instr string, targetA, targetB int) factory {
	rules := map[int]rule{}
	holding := map[int][]int{}
	f := factory{bots: map[int][]int{}, outputs: map[int]int{}, compareBot: -1}

	parseDest := func(kind, id string) dest {
		n, _ := strconv.Atoi(id)
		return dest{isOutput: kind == "output", id: n}
	}

	var ready []int
	give := func(bot, val int) {
		holding[bot] = append(holding[bot], val)
		if len(holding[bot]) == 2 {
			ready = append(ready, bot)
		}
	}

	for line := range strings.SplitSeq(strings.TrimSpace(instr), "\n") {
		w := strings.Fields(line)
		if w[0] == "value" {
			val, _ := strconv.Atoi(w[1])
			give(toInt(w[5]), val)
		} else { // bot B gives low to X x and high to Y y
			b := toInt(w[1])
			rules[b] = rule{
				low:  parseDest(w[5], w[6]),
				high: parseDest(w[10], w[11]),
			}
		}
	}

	for len(ready) > 0 {
		b := ready[0]
		ready = ready[1:]
		chips := holding[b]
		lo, hi := chips[0], chips[1]
		if lo > hi {
			lo, hi = hi, lo
		}
		f.bots[b] = []int{lo, hi}
		if (lo == targetA && hi == targetB) || (lo == targetB && hi == targetA) {
			f.compareBot = b
		}
		r := rules[b]
		send := func(d dest, val int) {
			if d.isOutput {
				f.outputs[d.id] = val
			} else {
				give(d.id, val)
			}
		}
		send(r.low, lo)
		send(r.high, hi)
		holding[b] = nil
	}
	return f
}

func toInt(s string) int {
	n, _ := strconv.Atoi(s)
	return n
}

// targets returns the two chip values to compare: 61 & 17 for the real input,
// 5 & 2 for the small example (which has no value above 5).
func targets(instr string) (int, int) {
	maxVal := 0
	for line := range strings.SplitSeq(instr, "\n") {
		w := strings.Fields(line)
		if len(w) >= 2 && w[0] == "value" {
			if v := toInt(w[1]); v > maxVal {
				maxVal = v
			}
		}
	}
	if maxVal > 5 {
		return 61, 17
	}
	return 5, 2
}

// One returns the answer to the first part of the exercise.
func (e Exercise) One(instr string) (any, error) {
	a, b := targets(instr)
	return simulate(instr, a, b).compareBot, nil
}

// Two returns the answer to the second part of the exercise.
func (e Exercise) Two(instr string) (any, error) {
	a, b := targets(instr)
	f := simulate(instr, a, b)
	return f.outputs[0] * f.outputs[1] * f.outputs[2], nil
}
