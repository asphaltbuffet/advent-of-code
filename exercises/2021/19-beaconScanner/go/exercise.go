package exercises

import (
	"fmt"
	"strings"

	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

// Exercise for Advent of Code 2021 day 19.
type Exercise struct {
	common.BaseExercise
}

type vec struct{ x, y, z int }

func (a vec) add(b vec) vec  { return vec{a.x + b.x, a.y + b.y, a.z + b.z} }
func (a vec) sub(b vec) vec  { return vec{a.x - b.x, a.y - b.y, a.z - b.z} }
func (a vec) manhattan() int { return abs(a.x) + abs(a.y) + abs(a.z) }

func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}

// parse groups beacon lines under each "--- scanner N ---" header, which is
// robust whether or not blank lines separate the blocks.
func parse(instr string) [][]vec {
	var scanners [][]vec
	var cur []vec
	started := false

	for _, line := range strings.Split(strings.TrimSpace(instr), "\n") {
		line = strings.TrimSpace(line)
		switch {
		case line == "":
			continue
		case strings.HasPrefix(line, "---"):
			if started {
				scanners = append(scanners, cur)
			}
			cur = nil
			started = true
		default:
			var v vec
			fmt.Sscanf(line, "%d,%d,%d", &v.x, &v.y, &v.z)
			cur = append(cur, v)
		}
	}
	if started {
		scanners = append(scanners, cur)
	}

	return scanners
}

// rotations returns the 24 axis-aligned orientations of a point. Built by
// composing the facing direction (6) with the roll about it (4).
func rotations(v vec) [24]vec {
	// Face transforms: point +x toward each of the 6 axes.
	faces := [6]func(vec) vec{
		func(p vec) vec { return p },                        // +x
		func(p vec) vec { return vec{-p.x, -p.y, p.z} },     // -x
		func(p vec) vec { return vec{p.y, -p.x, p.z} },      // +y
		func(p vec) vec { return vec{-p.y, p.x, p.z} },      // -y
		func(p vec) vec { return vec{p.z, p.y, -p.x} },      // +z
		func(p vec) vec { return vec{-p.z, p.y, p.x} },      // -z
	}
	// Roll about the (now) x axis.
	roll := func(p vec) vec { return vec{p.x, -p.z, p.y} }

	var out [24]vec
	i := 0
	for _, f := range faces {
		p := f(v)
		for r := 0; r < 4; r++ {
			out[i] = p
			i++
			p = roll(p)
		}
	}
	return out
}

// rotateAll applies rotation index r to every beacon in a scanner.
func rotateAll(beacons []vec, r int) []vec {
	out := make([]vec, len(beacons))
	for i, b := range beacons {
		out[i] = rotations(b)[r]
	}
	return out
}

// tryMatch checks whether scanner `unknown` overlaps the already-placed set
// `known` in some orientation. On success it returns the located beacons (in the
// global frame) and the scanner's position.
func tryMatch(known map[vec]bool, unknown []vec) ([]vec, vec, bool) {
	for r := 0; r < 24; r++ {
		rotated := rotateAll(unknown, r)

		// Count candidate translations: known - rotated. A translation shared by
		// >=12 beacon pairs aligns the scanners.
		offsets := make(map[vec]int, len(rotated)*4)
		for k := range known {
			for _, rb := range rotated {
				off := k.sub(rb)
				offsets[off]++
				if offsets[off] >= 12 {
					located := make([]vec, len(rotated))
					for i, rb2 := range rotated {
						located[i] = rb2.add(off)
					}
					return located, off, true
				}
			}
		}
	}
	return nil, vec{}, false
}

// assemble locates every scanner relative to scanner 0 and returns the set of
// all beacons and every scanner position.
func assemble(scanners [][]vec) (map[vec]bool, []vec) {
	beacons := map[vec]bool{}
	for _, b := range scanners[0] {
		beacons[b] = true
	}
	positions := []vec{{0, 0, 0}}

	placed := map[int]bool{0: true}

	// Repeatedly try to place remaining scanners against the growing map.
	for len(placed) < len(scanners) {
		progress := false
		for i := 1; i < len(scanners); i++ {
			if placed[i] {
				continue
			}
			located, pos, ok := tryMatch(beacons, scanners[i])
			if !ok {
				continue
			}
			for _, b := range located {
				beacons[b] = true
			}
			positions = append(positions, pos)
			placed[i] = true
			progress = true
		}
		if !progress {
			break
		}
	}

	return beacons, positions
}

// One returns the answer to the first part of the exercise.
func (e Exercise) One(instr string) (any, error) {
	scanners := parse(instr)
	if len(scanners) == 0 {
		return nil, fmt.Errorf("no scanners in input")
	}

	beacons, _ := assemble(scanners)

	return fmt.Sprintf("%d", len(beacons)), nil
}

// Two returns the answer to the second part of the exercise.
func (e Exercise) Two(instr string) (any, error) {
	scanners := parse(instr)
	if len(scanners) == 0 {
		return nil, fmt.Errorf("no scanners in input")
	}

	_, positions := assemble(scanners)

	best := 0
	for i := range positions {
		for j := i + 1; j < len(positions); j++ {
			if d := positions[i].sub(positions[j]).manhattan(); d > best {
				best = d
			}
		}
	}

	return fmt.Sprintf("%d", best), nil
}
