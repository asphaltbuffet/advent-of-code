package exercises

import (
	"regexp"
	"sort"
	"strings"

	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

// Exercise for Advent of Code 2016 day 11.
type Exercise struct {
	common.BaseExercise
}

// state holds the elevator floor and, per element, its generator and chip
// floors. Element identity does not matter for the search, only the multiset of
// (gen, chip) floor pairs.
type state struct {
	elevator int
	gens     []int // gens[i], chips[i] are floors (1-4) of element i
	chips    []int
}

var (
	genRe  = regexp.MustCompile(`(\w+) generator`)
	chipRe = regexp.MustCompile(`(\w+)-compatible microchip`)
)

var floorWords = map[string]int{"first": 0, "second": 1, "third": 2, "fourth": 3}

// parse reads the initial state. extra adds that many element pairs on floor 0.
func parse(instr string, extra int) state {
	gens := map[string]int{}
	chips := map[string]int{}
	for _, line := range strings.Split(strings.TrimSpace(instr), "\n") {
		fields := strings.Fields(line)
		floor := floorWords[fields[1]]
		for _, m := range genRe.FindAllStringSubmatch(line, -1) {
			gens[m[1]] = floor
		}
		for _, m := range chipRe.FindAllStringSubmatch(line, -1) {
			chips[m[1]] = floor
		}
	}

	var st state
	for el, gf := range gens {
		st.gens = append(st.gens, gf)
		st.chips = append(st.chips, chips[el])
	}
	for i := 0; i < extra; i++ {
		st.gens = append(st.gens, 0)
		st.chips = append(st.chips, 0)
	}
	return st
}

// valid reports whether no chip is fried on any floor.
func (s state) valid() bool {
	for f := 0; f < 4; f++ {
		hasGen := false
		for _, gf := range s.gens {
			if gf == f {
				hasGen = true
				break
			}
		}
		if !hasGen {
			continue
		}
		for i, cf := range s.chips {
			if cf == f && s.gens[i] != f {
				return false // unshielded chip with a foreign generator present
			}
		}
	}
	return true
}

// done reports whether everything is on the top floor.
func (s state) done() bool {
	for i := range s.gens {
		if s.gens[i] != 3 || s.chips[i] != 3 {
			return false
		}
	}
	return true
}

// key canonicalizes the state: elevator floor plus the sorted multiset of
// (gen, chip) floor pairs (element labels are interchangeable).
func (s state) key() string {
	pairs := make([][2]int, len(s.gens))
	for i := range s.gens {
		pairs[i] = [2]int{s.gens[i], s.chips[i]}
	}
	sort.Slice(pairs, func(a, b int) bool {
		if pairs[a][0] != pairs[b][0] {
			return pairs[a][0] < pairs[b][0]
		}
		return pairs[a][1] < pairs[b][1]
	})
	var b strings.Builder
	b.WriteByte(byte('0' + s.elevator))
	for _, p := range pairs {
		b.WriteByte(byte('0' + p[0]))
		b.WriteByte(byte('0' + p[1]))
	}
	return b.String()
}

// item identifies a carriable item: a generator (gen=true) or chip of element i.
type item struct {
	gen bool
	idx int
}

// itemsOn returns the items present on the elevator's floor.
func (s state) itemsOn() []item {
	var out []item
	for i, gf := range s.gens {
		if gf == s.elevator {
			out = append(out, item{gen: true, idx: i})
		}
	}
	for i, cf := range s.chips {
		if cf == s.elevator {
			out = append(out, item{gen: false, idx: i})
		}
	}
	return out
}

// clone deep-copies the state.
func (s state) clone() state {
	g := make([]int, len(s.gens))
	c := make([]int, len(s.chips))
	copy(g, s.gens)
	copy(c, s.chips)
	return state{elevator: s.elevator, gens: g, chips: c}
}

func (s *state) move(it item, floor int) {
	if it.gen {
		s.gens[it.idx] = floor
	} else {
		s.chips[it.idx] = floor
	}
}

// solve returns the minimum elevator steps to gather everything on floor 4.
func solve(start state) int {
	seen := map[string]bool{start.key(): true}
	type node struct {
		s     state
		steps int
	}
	queue := []node{{start, 0}}
	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]
		if cur.s.done() {
			return cur.steps
		}
		items := cur.s.itemsOn()
		// Candidate loads: each single item, and each pair.
		var loads [][]item
		for i := range items {
			loads = append(loads, []item{items[i]})
			for j := i + 1; j < len(items); j++ {
				loads = append(loads, []item{items[i], items[j]})
			}
		}
		for _, dir := range []int{1, -1} {
			nf := cur.s.elevator + dir
			if nf < 0 || nf > 3 {
				continue
			}
			for _, load := range loads {
				ns := cur.s.clone()
				ns.elevator = nf
				for _, it := range load {
					ns.move(it, nf)
				}
				if !ns.valid() {
					continue
				}
				k := ns.key()
				if seen[k] {
					continue
				}
				seen[k] = true
				queue = append(queue, node{ns, cur.steps + 1})
			}
		}
	}
	return -1
}

// One returns the answer to the first part of the exercise.
func (e Exercise) One(instr string) (any, error) {
	return solve(parse(instr, 0)), nil
}

// Two returns the answer to the second part of the exercise.
func (e Exercise) Two(instr string) (any, error) {
	return solve(parse(instr, 2)), nil
}
