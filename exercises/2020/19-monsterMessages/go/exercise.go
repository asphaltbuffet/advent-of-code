package exercises

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

// Exercise for Advent of Code 2020 day 19.
type Exercise struct {
	common.BaseExercise
}

// rule is one grammar rule: either a literal character, or a set of alternative
// sequences of sub-rule IDs.
type rule struct {
	literal byte
	isLit   bool
	alts    [][]int
}

type grammar map[int]rule

func parse(instr string) (grammar, []string, error) {
	blocks := strings.SplitN(strings.TrimRight(instr, "\n"), "\n\n", 2)
	if len(blocks) != 2 {
		return nil, nil, fmt.Errorf("expected two blocks")
	}

	g := grammar{}
	for _, line := range strings.Split(blocks[0], "\n") {
		idStr, body, ok := strings.Cut(line, ": ")
		if !ok {
			continue
		}
		id, err := strconv.Atoi(idStr)
		if err != nil {
			return nil, nil, fmt.Errorf("rule id %q: %w", idStr, err)
		}
		if strings.HasPrefix(body, `"`) {
			g[id] = rule{literal: body[1], isLit: true}
			continue
		}
		var r rule
		for _, alt := range strings.Split(body, " | ") {
			var seq []int
			for _, f := range strings.Fields(alt) {
				n, err := strconv.Atoi(f)
				if err != nil {
					return nil, nil, fmt.Errorf("rule ref %q: %w", f, err)
				}
				seq = append(seq, n)
			}
			r.alts = append(r.alts, seq)
		}
		g[id] = r
	}

	messages := strings.Split(strings.TrimSpace(blocks[1]), "\n")
	return g, messages, nil
}

// match returns the set of end positions reachable by matching rule id against
// s starting at pos. A message matches rule 0 iff len(s) is reachable from pos 0.
func match(g grammar, id int, s string, pos int) []int {
	r := g[id]
	if r.isLit {
		if pos < len(s) && s[pos] == r.literal {
			return []int{pos + 1}
		}
		return nil
	}

	var ends []int
	for _, seq := range r.alts {
		// Positions reachable after consuming this sequence, threaded left to right.
		cur := []int{pos}
		for _, sub := range seq {
			var next []int
			for _, p := range cur {
				next = append(next, match(g, sub, s, p)...)
			}
			cur = next
			if len(cur) == 0 {
				break
			}
		}
		ends = append(ends, cur...)
	}
	return ends
}

// matches reports whether s fully matches rule 0.
func matches(g grammar, s string) bool {
	for _, end := range match(g, 0, s, 0) {
		if end == len(s) {
			return true
		}
	}
	return false
}

// One counts messages that completely match rule 0.
func (e Exercise) One(instr string) (any, error) {
	g, messages, err := parse(instr)
	if err != nil {
		return nil, err
	}

	count := 0
	for _, m := range messages {
		if matches(g, m) {
			count++
		}
	}
	return fmt.Sprintf("%d", count), nil
}

// Two overrides rules 8 and 11 with their recursive forms, then counts matches.
// The set-of-positions matcher handles the recursion directly. If the grammar has
// no rule 8/11 (the small example), the overrides are harmless no-ops relative to
// what those ids already were, and the answer equals part one.
func (e Exercise) Two(instr string) (any, error) {
	g, messages, err := parse(instr)
	if err != nil {
		return nil, err
	}

	// Only rewrite if these rules exist (the real input and large example).
	if _, ok := g[8]; ok {
		g[8] = rule{alts: [][]int{{42}, {42, 8}}}
	}
	if _, ok := g[11]; ok {
		g[11] = rule{alts: [][]int{{42, 31}, {42, 11, 31}}}
	}

	count := 0
	for _, m := range messages {
		if matches(g, m) {
			count++
		}
	}
	return fmt.Sprintf("%d", count), nil
}
