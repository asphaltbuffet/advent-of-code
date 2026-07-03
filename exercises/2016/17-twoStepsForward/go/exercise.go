package exercises

import (
	"crypto/md5"
	"encoding/hex"
	"strconv"
	"strings"

	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

// Exercise for Advent of Code 2016 day 17.
type Exercise struct {
	common.BaseExercise
}

// move describes a direction: its letter and the (dx, dy) it applies.
type move struct {
	letter byte
	dx, dy int
}

// Order matters: hash chars 0..3 map to up, down, left, right.
var moves = []move{
	{'U', 0, -1},
	{'D', 0, 1},
	{'L', -1, 0},
	{'R', 1, 0},
}

// openDoors returns which of the four doors are open from the room reached by
// path, given the passcode.
func openDoors(passcode, path string) [4]bool {
	sum := md5.Sum([]byte(passcode + path))
	h := hex.EncodeToString(sum[:])
	var doors [4]bool
	for i := 0; i < 4; i++ {
		doors[i] = h[i] >= 'b' && h[i] <= 'f'
	}
	return doors
}

// search walks every path to the vault, returning the shortest and the length
// of the longest.
func search(passcode string) (shortest string, longest int) {
	type node struct {
		x, y int
		path string
	}
	queue := []node{{0, 0, ""}}
	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]
		if cur.x == 3 && cur.y == 3 {
			if shortest == "" {
				shortest = cur.path
			}
			if len(cur.path) > longest {
				longest = len(cur.path)
			}
			continue // cannot move once at the vault
		}
		doors := openDoors(passcode, cur.path)
		for i, m := range moves {
			if !doors[i] {
				continue
			}
			nx, ny := cur.x+m.dx, cur.y+m.dy
			if nx < 0 || nx > 3 || ny < 0 || ny > 3 {
				continue
			}
			queue = append(queue, node{nx, ny, cur.path + string(m.letter)})
		}
	}
	return shortest, longest
}

// One returns the answer to the first part of the exercise.
func (e Exercise) One(instr string) (any, error) {
	shortest, _ := search(strings.TrimSpace(instr))
	return shortest, nil
}

// Two returns the answer to the second part of the exercise.
func (e Exercise) Two(instr string) (any, error) {
	_, longest := search(strings.TrimSpace(instr))
	return strconv.Itoa(longest), nil
}
