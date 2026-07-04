package exercises

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

// Exercise for Advent of Code 2019 day 12.
type Exercise struct {
	common.BaseExercise
}

var moonRe = regexp.MustCompile(`<x=(-?\d+), y=(-?\d+), z=(-?\d+)>`)

type vec3 [3]int

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func sign(a, b int) int {
	if a < b {
		return 1
	} else if a > b {
		return -1
	}
	return 0
}

// simulate runs n steps and returns the total energy.
func simulate(positions []vec3, steps int) int {
	n := len(positions)
	pos := make([]vec3, n)
	vel := make([]vec3, n)
	copy(pos, positions)

	for range steps {
		// Apply gravity.
		for i := range n {
			for j := i + 1; j < n; j++ {
				for axis := range 3 {
					dv := sign(pos[i][axis], pos[j][axis])
					vel[i][axis] += dv
					vel[j][axis] -= dv
				}
			}
		}
		// Apply velocity.
		for i := range n {
			for axis := range 3 {
				pos[i][axis] += vel[i][axis]
			}
		}
	}

	// Calculate total energy.
	total := 0
	for i := range n {
		pot := abs(pos[i][0]) + abs(pos[i][1]) + abs(pos[i][2])
		kin := abs(vel[i][0]) + abs(vel[i][1]) + abs(vel[i][2])
		total += pot * kin
	}
	return total
}

// parseInput extracts step count and moon positions.
// If the first line starts with "Steps:", that value is used; otherwise 1000.
func parseInput(instr string) ([]vec3, int) {
	lines := strings.Split(strings.TrimSpace(instr), "\n")
	steps := 1000
	moonLines := lines

	if rest, ok := strings.CutPrefix(lines[0], "Steps:"); ok {
		n, _ := strconv.Atoi(strings.TrimSpace(rest))
		steps = n
		moonLines = lines[1:]
	}

	var positions []vec3
	for _, line := range moonLines {
		m := moonRe.FindStringSubmatch(line)
		if m == nil {
			continue
		}
		x, _ := strconv.Atoi(m[1])
		y, _ := strconv.Atoi(m[2])
		z, _ := strconv.Atoi(m[3])
		positions = append(positions, vec3{x, y, z})
	}
	return positions, steps
}

// One returns the answer to the first part of the exercise.
// Calculates total energy after N steps (default 1000, or from "Steps: N" header).
func (e Exercise) One(instr string) (any, error) {
	positions, steps := parseInput(instr)
	return simulate(positions, steps), nil
}

func gcd(a, b int64) int64 {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func lcm(a, b int64) int64 {
	return a / gcd(a, b) * b
}

// axisCycle finds how many steps until the given axis returns to its initial state.
func axisCycle(positions []vec3, axis int) int64 {
	n := len(positions)
	pos := make([]int, n)
	vel := make([]int, n)
	for i, p := range positions {
		pos[i] = p[axis]
	}
	// Record initial state.
	initPos := make([]int, n)
	initVel := make([]int, n)
	copy(initPos, pos)
	copy(initVel, vel)

	var step int64
	for {
		// Apply gravity.
		for i := range n {
			for j := i + 1; j < n; j++ {
				dv := sign(pos[i], pos[j])
				vel[i] += dv
				vel[j] -= dv
			}
		}
		// Apply velocity.
		for i := range n {
			pos[i] += vel[i]
		}
		step++
		// Check if back to initial state.
		match := true
		for i := range n {
			if pos[i] != initPos[i] || vel[i] != initVel[i] {
				match = false
				break
			}
		}
		if match {
			return step
		}
	}
}

// Two returns the answer to the second part of the exercise.
// Finds the number of steps until the moon system returns to its initial state.
func (e Exercise) Two(instr string) (any, error) {
	positions, _ := parseInput(instr)
	cx := axisCycle(positions, 0)
	cy := axisCycle(positions, 1)
	cz := axisCycle(positions, 2)
	return strconv.FormatInt(lcm(lcm(cx, cy), cz), 10), nil
}
