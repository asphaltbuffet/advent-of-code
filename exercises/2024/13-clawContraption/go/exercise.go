package exercises

import (
	"regexp"
	"strconv"

	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

// Exercise for Advent of Code 2024 day 13.
type Exercise struct {
	common.BaseExercise
}

// machine holds the two button vectors and the prize location.
type machine struct {
	ax, ay int // Button A displacement
	bx, by int // Button B displacement
	px, py int // Prize coordinates
}

var numRe = regexp.MustCompile(`-?\d+`)

// parseMachines extracts every group of six integers (ax, ay, bx, by, px, py).
// Reading numbers positionally rather than by fixed line index keeps it robust
// to stray whitespace, CRLF, or a missing trailing newline.
func parseMachines(instr string) []machine {
	nums := numRe.FindAllString(instr, -1)
	var machines []machine
	for i := 0; i+5 < len(nums); i += 6 {
		v := make([]int, 6)
		for j := 0; j < 6; j++ {
			v[j], _ = strconv.Atoi(nums[i+j])
		}
		machines = append(machines, machine{v[0], v[1], v[2], v[3], v[4], v[5]})
	}
	return machines
}

// cost solves the 2x2 system via Cramer's rule and returns the token cost
// (3*na + nb), or 0 if there is no non-negative integer solution.
func (m machine) cost(offset int) int {
	px, py := m.px+offset, m.py+offset
	det := m.ax*m.by - m.ay*m.bx
	if det == 0 {
		return 0
	}
	naNum := px*m.by - py*m.bx
	nbNum := m.ax*py - m.ay*px
	if naNum%det != 0 || nbNum%det != 0 {
		return 0
	}
	na, nb := naNum/det, nbNum/det
	if na < 0 || nb < 0 {
		return 0
	}
	return 3*na + nb
}

func totalCost(instr string, offset int) int {
	total := 0
	for _, m := range parseMachines(instr) {
		total += m.cost(offset)
	}
	return total
}

// One returns the answer to the first part of the exercise.
func (e Exercise) One(instr string) (any, error) {
	return totalCost(instr, 0), nil
}

// Two returns the answer to the second part of the exercise.
func (e Exercise) Two(instr string) (any, error) {
	return totalCost(instr, 10000000000000), nil
}
