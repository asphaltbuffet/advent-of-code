package exercises

import (
	"fmt"
	"strings"

	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

// Exercise for Advent of Code 2020 day 17.
type Exercise struct {
	common.BaseExercise
}

// coord is a 4D point; 3D simulations pin the w axis to 0.
type coord [4]int

// parseActive reads the initial slice into a set of active coordinates on the z=0
// (and w=0) plane.
func parseActive(instr string) map[coord]bool {
	active := map[coord]bool{}
	for y, line := range strings.Split(strings.TrimSpace(instr), "\n") {
		for x, c := range line {
			if c == '#' {
				active[coord{x, y, 0, 0}] = true
			}
		}
	}
	return active
}

// neighborOffsets returns every nonzero offset over `dims` dimensions (the
// remaining axes stay 0), i.e. 26 offsets for 3D and 80 for 4D.
func neighborOffsets(dims int) []coord {
	var offsets []coord
	ranges := [4][2]int{{0, 0}, {0, 0}, {0, 0}, {0, 0}}
	for i := 0; i < dims; i++ {
		ranges[i] = [2]int{-1, 1}
	}
	for dx := ranges[0][0]; dx <= ranges[0][1]; dx++ {
		for dy := ranges[1][0]; dy <= ranges[1][1]; dy++ {
			for dz := ranges[2][0]; dz <= ranges[2][1]; dz++ {
				for dw := ranges[3][0]; dw <= ranges[3][1]; dw++ {
					if dx == 0 && dy == 0 && dz == 0 && dw == 0 {
						continue
					}
					offsets = append(offsets, coord{dx, dy, dz, dw})
				}
			}
		}
	}
	return offsets
}

// simulateTo runs `cycles` steps of the Conway rules over `dims` dimensions and
// returns the final active set.
func simulateTo(active map[coord]bool, dims, cycles int) map[coord]bool {
	offsets := neighborOffsets(dims)

	for step := 0; step < cycles; step++ {
		// Tally active-neighbor counts for every cell adjacent to an active one.
		counts := map[coord]int{}
		for c := range active {
			for _, o := range offsets {
				n := coord{c[0] + o[0], c[1] + o[1], c[2] + o[2], c[3] + o[3]}
				counts[n]++
			}
		}

		next := map[coord]bool{}
		for c, n := range counts {
			if n == 3 || (n == 2 && active[c]) {
				next[c] = true
			}
		}
		active = next
	}

	return active
}

// simulate returns the final active count after running the rules.
func simulate(active map[coord]bool, dims, cycles int) int {
	return len(simulateTo(active, dims, cycles))
}

// One returns the number of active cubes after 6 cycles in 3D.
func (e Exercise) One(instr string) (any, error) {
	return fmt.Sprintf("%d", simulate(parseActive(instr), 3, 6)), nil
}

// Two returns the number of active cells after 6 cycles in 4D.
func (e Exercise) Two(instr string) (any, error) {
	return fmt.Sprintf("%d", simulate(parseActive(instr), 4, 6)), nil
}
