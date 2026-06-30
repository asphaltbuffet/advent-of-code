package exercises

import (
	"math/rand"
	"strings"

	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

// Exercise for Advent of Code 2015 day 19.
type Exercise struct {
	common.BaseExercise
}

type rule struct{ from, to string }

// parse splits the input into the replacement rules and the target molecule.
func parse(instr string) ([]rule, string) {
	blocks := strings.SplitN(strings.ReplaceAll(instr, "\r\n", "\n"), "\n\n", 2)

	var rules []rule
	for _, line := range strings.Split(strings.TrimSpace(blocks[0]), "\n") {
		parts := strings.Split(line, " => ")
		rules = append(rules, rule{parts[0], parts[1]})
	}

	return rules, strings.TrimSpace(blocks[1])
}

// One counts the distinct molecules produced by applying any single rule once.
func (e Exercise) One(instr string) (any, error) {
	rules, mol := parse(instr)
	seen := map[string]struct{}{}

	for _, r := range rules {
		// Replace each occurrence of r.from independently.
		for i := 0; i+len(r.from) <= len(mol); i++ {
			if mol[i:i+len(r.from)] == r.from {
				seen[mol[:i]+r.to+mol[i+len(r.from):]] = struct{}{}
			}
		}
	}

	return len(seen), nil
}

// Two returns the fewest replacement steps to build the molecule from "e".
// It works backwards greedily: repeatedly collapse any production "to" back to
// its "from" until only "e" remains. Greedy can dead-end, so on getting stuck
// it reshuffles the rule order and retries — input-agnostic for both the small
// example grammar and the real Rn/Ar/Y grammar.
func (e Exercise) Two(instr string) (any, error) {
	rules, mol := parse(instr)
	rng := rand.New(rand.NewSource(1))

	for {
		cur := mol
		steps := 0
		stuck := false

		for cur != "e" {
			applied := false
			for _, r := range rules {
				if idx := strings.Index(cur, r.to); idx >= 0 {
					cur = cur[:idx] + r.from + cur[idx+len(r.to):]
					steps++
					applied = true
					break
				}
			}
			if !applied {
				stuck = true
				break
			}
		}

		if !stuck {
			return steps, nil
		}

		// Dead-ended: shuffle rules and try again.
		rng.Shuffle(len(rules), func(i, j int) { rules[i], rules[j] = rules[j], rules[i] })
	}
}
