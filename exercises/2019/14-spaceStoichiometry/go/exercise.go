package exercises

import (
	"errors"
	"strconv"
	"strings"

	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

// Exercise for Advent of Code 2019 day 14.
type Exercise struct {
	common.BaseExercise
}

type ingredient struct {
	qty  int
	name string
}

type reaction struct {
	qty    int
	inputs []ingredient
}

func parseIngredient(s string) ingredient {
	s = strings.TrimSpace(s)
	parts := strings.SplitN(s, " ", 2)
	qty, _ := strconv.Atoi(parts[0])
	return ingredient{qty: qty, name: strings.TrimSpace(parts[1])}
}

func parseReactions(instr string) map[string]reaction {
	reactions := make(map[string]reaction)
	for line := range strings.SplitSeq(strings.TrimSpace(instr), "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		parts := strings.Split(line, "=>")
		output := parseIngredient(parts[1])
		var inputs []ingredient
		for inp := range strings.SplitSeq(parts[0], ",") {
			inputs = append(inputs, parseIngredient(inp))
		}
		reactions[output.name] = reaction{qty: output.qty, inputs: inputs}
	}
	return reactions
}

func oreForFuel(reactions map[string]reaction, fuelAmt int) int {
	need := map[string]int{"FUEL": fuelAmt}
	surplus := make(map[string]int)

	for {
		// Find a non-ORE chemical we still need
		var chem string
		for k, v := range need {
			if k != "ORE" && v > 0 {
				chem = k
				break
			}
		}
		if chem == "" {
			break
		}

		qty := need[chem]
		delete(need, chem)

		// Use surplus first
		if surplus[chem] >= qty {
			surplus[chem] -= qty
			continue
		}
		qty -= surplus[chem]
		surplus[chem] = 0

		r := reactions[chem]
		// How many times do we need to apply this reaction?
		times := (qty + r.qty - 1) / r.qty
		produced := times * r.qty
		leftover := produced - qty
		surplus[chem] += leftover

		for _, inp := range r.inputs {
			need[inp.name] += times * inp.qty
		}
	}

	return need["ORE"]
}

// One returns the answer to the first part of the exercise.
func (e Exercise) One(instr string) (any, error) {
	reactions := parseReactions(instr)
	return oreForFuel(reactions, 1), nil
}

// Two returns the answer to the second part of the exercise.
func (e Exercise) Two(_ string) (any, error) {
	return nil, errors.New("part 2 not implemented")
}
