package exercises

import (
	"strconv"
	"strings"

	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

// Exercise for Advent of Code 2019 day 11.
type Exercise struct {
	common.BaseExercise
}

// mem is a map-backed unbounded memory for Intcode.
type mem map[int]int

func (m mem) get(addr int) int  { return m[addr] }
func (m mem) set(addr, val int) { m[addr] = val }

// parseProgram parses a comma-separated Intcode program into a mem map.
func parseProgram(input string) mem {
	m := make(mem)
	for i, s := range strings.Split(strings.TrimSpace(input), ",") {
		v, _ := strconv.Atoi(s)
		m[i] = v
	}
	return m
}

// intcode runs an Intcode computer in a goroutine, communicating via channels.
//
//nolint:gocognit,funlen // Intcode VM opcode dispatch is inherently complex
func intcode(m mem, in <-chan int, out chan<- int) {
	ip := 0
	relBase := 0

	param := func(offset, mode int) int {
		raw := m.get(ip + offset)
		switch mode {
		case 0:
			return m.get(raw)
		case 1:
			return raw
		case 2:
			return m.get(relBase + raw)
		}
		panic("unknown mode")
	}

	dest := func(offset, mode int) int {
		raw := m.get(ip + offset)
		switch mode {
		case 0:
			return raw
		case 2:
			return relBase + raw
		}
		panic("unknown dest mode")
	}

	for {
		instr := m.get(ip)
		op := instr % 100
		m1 := (instr / 100) % 10
		m2 := (instr / 1000) % 10
		m3 := (instr / 10000) % 10

		switch op {
		case 1:
			m.set(dest(3, m3), param(1, m1)+param(2, m2))
			ip += 4
		case 2:
			m.set(dest(3, m3), param(1, m1)*param(2, m2))
			ip += 4
		case 3:
			m.set(dest(1, m1), <-in)
			ip += 2
		case 4:
			out <- param(1, m1)
			ip += 2
		case 5:
			if param(1, m1) != 0 {
				ip = param(2, m2)
			} else {
				ip += 3
			}
		case 6:
			if param(1, m1) == 0 {
				ip = param(2, m2)
			} else {
				ip += 3
			}
		case 7:
			v := 0
			if param(1, m1) < param(2, m2) {
				v = 1
			}
			m.set(dest(3, m3), v)
			ip += 4
		case 8:
			v := 0
			if param(1, m1) == param(2, m2) {
				v = 1
			}
			m.set(dest(3, m3), v)
			ip += 4
		case 9:
			relBase += param(1, m1)
			ip += 2
		case 99:
			close(out)
			return
		}
	}
}

// pos is a grid position.
type pos [2]int

// runRobot runs the hull-painting robot and returns the map of painted panels.
// startColor is the initial color of panel (0,0): 0=black, 1=white.
func runRobot(program mem, startColor int) map[pos]int {
	in := make(chan int, 1)
	out := make(chan int)
	go intcode(program, in, out)

	painted := make(map[pos]int)
	grid := make(map[pos]int)
	grid[pos{0, 0}] = startColor

	// Directions: up, right, down, left (dx, dy) — y increases upward.
	dirs := [][2]int{{0, 1}, {1, 0}, {0, -1}, {-1, 0}}
	dir := 0 // start facing up
	cur := pos{0, 0}

	for {
		// Send current panel color as input.
		in <- grid[cur]

		// Read paint color.
		color, ok := <-out
		if !ok {
			break
		}
		// Read turn direction.
		turn, ok := <-out
		if !ok {
			break
		}

		// Paint.
		grid[cur] = color
		painted[cur] = color

		// Turn: 0=left(-1), 1=right(+1).
		if turn == 0 {
			dir = (dir + 3) % 4
		} else {
			dir = (dir + 1) % 4
		}

		// Move forward.
		cur[0] += dirs[dir][0]
		cur[1] += dirs[dir][1]
	}

	return painted
}

// One returns the answer to the first part of the exercise.
// Counts how many panels the robot paints at least once, starting on a black panel.
func (e Exercise) One(instr string) (any, error) {
	painted := runRobot(parseProgram(instr), 0)
	return len(painted), nil
}

// Two returns the answer to the second part of the exercise.
// Runs the robot starting on a white panel and renders the registration identifier.
//
//nolint:gocognit,nestif // robot painting requires tracking bounding box in a first-pass branch
func (e Exercise) Two(instr string) (any, error) {
	painted := runRobot(parseProgram(instr), 1)

	// Find bounding box of all painted panels.
	minX, maxX, minY, maxY := 0, 0, 0, 0
	first := true
	for p := range painted {
		if first {
			minX, maxX, minY, maxY = p[0], p[0], p[1], p[1]
			first = false
		} else {
			if p[0] < minX {
				minX = p[0]
			}
			if p[0] > maxX {
				maxX = p[0]
			}
			if p[1] < minY {
				minY = p[1]
			}
			if p[1] > maxY {
				maxY = p[1]
			}
		}
	}

	// Render: y increases upward, so iterate from maxY down to minY.
	var sb strings.Builder
	for y := maxY; y >= minY; y-- {
		if y < maxY {
			sb.WriteByte('\n')
		}
		for x := minX; x <= maxX; x++ {
			if painted[pos{x, y}] == 1 {
				sb.WriteString("█")
			} else {
				sb.WriteString("░")
			}
		}
	}

	return sb.String(), nil
}
