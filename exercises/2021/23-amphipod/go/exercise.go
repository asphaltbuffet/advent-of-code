package exercises

import (
	"container/heap"
	"fmt"
	"strings"

	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

// Exercise for Advent of Code 2021 day 23.
type Exercise struct {
	common.BaseExercise
}

// The burrow is a hallway of 11 cells (0..10) plus four rooms. Room j (0..3) sits
// under hallway cell roomHall[j]; amphipods may never stop on those four cells.
// Each room is stored as a fixed-length byte string of length `depth`, slot 0
// nearest the hallway, with '.' for an empty slot.
var (
	roomHall = [4]int{2, 4, 6, 8}
	stepCost = map[byte]int{'A': 1, 'B': 10, 'C': 100, 'D': 1000}
)

const empty = '.'

// state is the full burrow: hallway plus room contents. Both use '.' for empty so
// the value is directly comparable and usable as a map key.
type state struct {
	hall  [11]byte
	rooms [4]string
}

func parseBurrow(instr string) state {
	// Gather amphipod letters row by row (top room row first).
	var rows [][]byte
	for _, line := range strings.Split(instr, "\n") {
		var letters []byte
		for i := 0; i < len(line); i++ {
			if line[i] >= 'A' && line[i] <= 'D' {
				letters = append(letters, line[i])
			}
		}
		if len(letters) == 4 {
			rows = append(rows, letters)
		}
	}

	s := state{}
	for i := range s.hall {
		s.hall[i] = empty
	}
	for j := 0; j < 4; j++ {
		b := make([]byte, len(rows))
		for r := 0; r < len(rows); r++ {
			b[r] = rows[r][j]
		}
		s.rooms[j] = string(b)
	}
	return s
}

// done reports whether every room holds only its own amphipod type.
func done(s state) bool {
	for j := 0; j < 4; j++ {
		want := byte('A' + j)
		for i := 0; i < len(s.rooms[j]); i++ {
			if s.rooms[j][i] != want {
				return false
			}
		}
	}
	return true
}

// hallClear reports whether every hallway cell strictly between a and b, plus b
// itself, is empty (a is the amphipod's own cell and is ignored).
func hallClear(s state, a, b int) bool {
	step := 1
	if b < a {
		step = -1
	}
	for c := a + step; ; c += step {
		if s.hall[c] != empty {
			return false
		}
		if c == b {
			return true
		}
	}
}

// withRoomByte returns a copy of a room string with slot i set to v.
func withRoomByte(room string, i int, v byte) string {
	b := []byte(room)
	b[i] = v
	return string(b)
}

type transition struct {
	st   state
	cost int
}

// moves generates all legal transitions from s.
func moves(s state, depth int) []transition {
	var out []transition

	// roomAcceptSlot returns the deepest empty slot of room j if the room holds
	// only its own type (so it can accept), else -1.
	roomAcceptSlot := func(j int) int {
		want := byte('A' + j)
		deepest := -1
		for i := 0; i < depth; i++ {
			switch s.rooms[j][i] {
			case empty:
				deepest = i
			case want:
			default:
				return -1 // a wrong occupant blocks the room
			}
		}
		return deepest
	}

	// topOccupant returns the shallowest occupied slot of room j, or -1.
	topOccupant := func(j int) int {
		for i := 0; i < depth; i++ {
			if s.rooms[j][i] != empty {
				return i
			}
		}
		return -1
	}

	// 1) Hallway -> its destination room.
	for c := 0; c < 11; c++ {
		a := s.hall[c]
		if a == empty {
			continue
		}
		j := int(a - 'A')
		slot := roomAcceptSlot(j)
		if slot == -1 || !hallClear(s, c, roomHall[j]) {
			continue
		}
		dist := abs(c-roomHall[j]) + (slot + 1)
		ns := s
		ns.hall[c] = empty
		ns.rooms[j] = withRoomByte(ns.rooms[j], slot, a)
		out = append(out, transition{ns, dist * stepCost[a]})
	}

	// 2) Room -> hallway.
	for j := 0; j < 4; j++ {
		top := topOccupant(j)
		if top == -1 {
			continue
		}
		// Skip if this room is already all correct from `top` down.
		want := byte('A' + j)
		settled := true
		for i := top; i < depth; i++ {
			if s.rooms[j][i] != want {
				settled = false
				break
			}
		}
		if settled {
			continue
		}
		a := s.rooms[j][top]
		for c := 0; c < 11; c++ {
			if c == 2 || c == 4 || c == 6 || c == 8 {
				continue // cannot stop directly above a room
			}
			if !hallClear(s, roomHall[j], c) {
				continue
			}
			dist := (top + 1) + abs(c-roomHall[j])
			ns := s
			ns.rooms[j] = withRoomByte(ns.rooms[j], top, empty)
			ns.hall[c] = a
			out = append(out, transition{ns, dist * stepCost[a]})
		}
	}

	return out
}

type item struct {
	s    state
	cost int
}

type pq []item

func (p pq) Len() int           { return len(p) }
func (p pq) Less(i, j int) bool { return p[i].cost < p[j].cost }
func (p pq) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p *pq) Push(x any)        { *p = append(*p, x.(item)) }
func (p *pq) Pop() any {
	old := *p
	n := len(old)
	it := old[n-1]
	*p = old[:n-1]
	return it
}

// solve runs Dijkstra over burrow states and returns the least energy to sort.
func solve(start state) int {
	depth := len(start.rooms[0])
	dist := map[state]int{start: 0}
	q := &pq{{start, 0}}
	for q.Len() > 0 {
		cur := heap.Pop(q).(item)
		if done(cur.s) {
			return cur.cost
		}
		if cur.cost > dist[cur.s] {
			continue
		}
		for _, m := range moves(cur.s, depth) {
			nc := cur.cost + m.cost
			if d, ok := dist[m.st]; !ok || nc < d {
				dist[m.st] = nc
				heap.Push(q, item{m.st, nc})
			}
		}
	}
	return -1
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// One sorts the two-deep burrow with the least energy.
func (e Exercise) One(instr string) (any, error) {
	return fmt.Sprintf("%d", solve(parseBurrow(instr))), nil
}

// unfoldRows are the two rows inserted between the folded input's room rows for
// Part Two, top-most first: DCBA over DBAC.
var unfoldRows = [2]string{"DCBA", "DBAC"}

// Two sorts the four-deep burrow, inserting the fixed extra rows if the input is
// still folded (only two room rows).
func (e Exercise) Two(instr string) (any, error) {
	s := parseBurrow(instr)

	if len(s.rooms[0]) == 2 {
		for j := 0; j < 4; j++ {
			r := s.rooms[j]
			// r[0] is the top slot, r[1] the bottom of the folded pair.
			s.rooms[j] = string([]byte{r[0], unfoldRows[0][j], unfoldRows[1][j], r[1]})
		}
	}

	return fmt.Sprintf("%d", solve(s)), nil
}
