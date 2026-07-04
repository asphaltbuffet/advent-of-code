package exercises

import (
	"errors"
	"strconv"
	"strings"

	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

// Exercise for Advent of Code 2019 day 17.
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

// intcode runs an Intcode computer in a goroutine, sending all output to the out channel.
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

// buildGrid runs the Intcode program and decodes its ASCII output into grid rows.
func buildGrid(instr string) []string {
	in := make(chan int, 1)
	out := make(chan int, 4096)

	go intcode(parseProgram(instr), in, out)

	var buf strings.Builder

	for v := range out {
		buf.WriteByte(byte(v))
	}

	// Split on newlines and drop empty trailing lines.
	rows := strings.Split(strings.TrimRight(buf.String(), "\n"), "\n")

	return rows
}

// One returns the answer to the first part of the exercise.
// Runs the Intcode camera program, builds the scaffold map, finds intersections,
// and sums their alignment parameters (x * y).
func (e Exercise) One(instr string) (any, error) {
	grid := buildGrid(instr)

	rows := len(grid)
	if rows == 0 {
		return 0, nil
	}

	sum := 0

	for y := 1; y < rows-1; y++ {
		row := grid[y]
		cols := len(row)

		for x := 1; x < cols-1; x++ {
			if row[x] == '#' &&
				grid[y-1][x] == '#' &&
				grid[y+1][x] == '#' &&
				row[x-1] == '#' &&
				row[x+1] == '#' {
				sum += x * y
			}
		}
	}

	return sum, nil
}

// walkPath builds the full movement path from the scaffold grid.
// Returns a slice of tokens like ["R","8","L","4","R","8",...].
// robotStart finds the robot's initial position and direction in the grid.
func robotStart(grid []string) (int, int, int) {
	var rx, ry, dir int
	for y, row := range grid {
		for x := range len(row) {
			switch row[x] {
			case '^':
				rx, ry, dir = x, y, 0
			case '>':
				rx, ry, dir = x, y, 1
			case 'v':
				rx, ry, dir = x, y, 2
			case '<':
				rx, ry, dir = x, y, 3
			}
		}
	}
	return rx, ry, dir
}

// chooseTurn returns the new direction and token ("L"/"R") if a turn is possible, else ok=false.
func chooseTurn(dir, rx, ry int, dx, dy [4]int, scaffold func(int, int) bool) (int, string, bool) {
	turnLeft := (dir + 3) % 4
	if scaffold(rx+dx[turnLeft], ry+dy[turnLeft]) {
		return turnLeft, "L", true
	}
	turnRight := (dir + 1) % 4
	if scaffold(rx+dx[turnRight], ry+dy[turnRight]) {
		return turnRight, "R", true
	}
	return 0, "", false
}

func walkPath(grid []string) []string {
	rows := len(grid)
	if rows == 0 {
		return nil
	}

	cols := len(grid[0])

	// Directions: 0=up, 1=right, 2=down, 3=left
	dx := [4]int{0, 1, 0, -1}
	dy := [4]int{-1, 0, 1, 0}

	rx, ry, dir := robotStart(grid)

	scaffold := func(x, y int) bool {
		if y < 0 || y >= rows || x < 0 || x >= cols {
			return false
		}
		c := grid[y][x]
		return c == '#' || c == '^' || c == 'v' || c == '<' || c == '>'
	}

	var tokens []string

	for {
		newDir, token, ok := chooseTurn(dir, rx, ry, dx, dy, scaffold)
		if !ok {
			break
		}
		dir = newDir
		tokens = append(tokens, token)

		steps := 0
		for {
			nx, ny := rx+dx[dir], ry+dy[dir]
			if !scaffold(nx, ny) {
				break
			}
			rx, ry = nx, ny
			steps++
		}
		tokens = append(tokens, strconv.Itoa(steps))
	}

	return tokens
}

// compress attempts to split tokens into a main routine + 3 functions (A, B, C),
// each ≤ 20 chars when comma-joined. Returns (main, A, B, C, ok).
// tryBC searches for valid B and C functions given a sequence already partially replaced with A.
func tryBC(afterA []string) ([]string, []string, []string, bool) {
	bStart := firstUnassigned(afterA)
	if bStart == -1 {
		return nil, nil, nil, false
	}

	bPrefix := prefixTokens(afterA, bStart)

	for bLen := 2; bLen <= len(bPrefix); bLen += 2 {
		bFunc := bPrefix[:bLen]
		if !routeFits(bFunc) {
			break
		}

		afterB := replaceLabel(afterA, bFunc, "B")
		mainSeq, cFunc, found := tryC(afterB)
		if found {
			return mainSeq, bFunc, cFunc, true
		}
	}

	return nil, nil, nil, false
}

func tryC(afterB []string) ([]string, []string, bool) {
	cStart := firstUnassigned(afterB)
	if cStart == -1 {
		return nil, nil, false
	}

	cPrefix := prefixTokens(afterB, cStart)

	for cLen := 2; cLen <= len(cPrefix); cLen += 2 {
		cFunc := cPrefix[:cLen]
		if !routeFits(cFunc) {
			break
		}

		afterC := replaceLabel(afterB, cFunc, "C")
		if onlyLabels(afterC) && routeFits(afterC) {
			return afterC, cFunc, true
		}
	}

	return nil, nil, false
}

func routeFits(toks []string) bool {
	return len(strings.Join(toks, ",")) <= 20
}

func replaceLabel(seq, pattern []string, label string) []string {
	var result []string
	i := 0
	for i < len(seq) {
		if i+len(pattern) <= len(seq) {
			match := true
			for j, t := range pattern {
				if seq[i+j] != t {
					match = false
					break
				}
			}
			if match {
				result = append(result, label)
				i += len(pattern)
				continue
			}
		}
		result = append(result, seq[i])
		i++
	}
	return result
}

func firstUnassigned(seq []string) int {
	for i, t := range seq {
		if t != "A" && t != "B" && t != "C" {
			return i
		}
	}
	return -1
}

func prefixTokens(seq []string, idx int) []string {
	var toks []string
	for i := idx; i < len(seq) && seq[i] != "A" && seq[i] != "B" && seq[i] != "C"; i++ {
		toks = append(toks, seq[i])
	}
	return toks
}

func onlyLabels(seq []string) bool {
	for _, t := range seq {
		if t != "A" && t != "B" && t != "C" {
			return false
		}
	}
	return true
}

func compress(tokens []string) ([]string, []string, []string, []string, bool) {
	basePrefix := prefixTokens(tokens, 0)

	for aLen := 2; aLen <= len(basePrefix); aLen += 2 {
		aFunc := basePrefix[:aLen]
		if !routeFits(aFunc) {
			break
		}
		afterA := replaceLabel(tokens, aFunc, "A")
		mainSeq, bFunc, cFunc, found := tryBC(afterA)
		if found {
			return mainSeq, aFunc, bFunc, cFunc, true
		}
	}

	return nil, nil, nil, nil, false
}

// Two returns the answer to the second part of the exercise.
// Wakes the robot, computes the full path, compresses it into 3 functions,
// feeds them as ASCII to the Intcode program, and returns dust collected.
func (e Exercise) Two(instr string) (any, error) {
	// Build grid using mode 1 to find path.
	grid := buildGrid(instr)

	tokens := walkPath(grid)

	mainRoute, aFunc, bFunc, cFunc, ok := compress(tokens)
	if !ok {
		return nil, errors.New("no valid compression found")
	}

	// Build ASCII input.
	join := func(toks []string) string { return strings.Join(toks, ",") }

	ascii := join(mainRoute) + "\n" +
		join(aFunc) + "\n" +
		join(bFunc) + "\n" +
		join(cFunc) + "\n" +
		"n\n"

	// Run Intcode with address 0 = 2.
	prog := parseProgram(instr)
	prog.set(0, 2)

	in := make(chan int, len(ascii))
	out := make(chan int, 4096)

	for _, ch := range ascii {
		in <- int(ch)
	}

	go intcode(prog, in, out)

	var last int

	for v := range out {
		last = v
	}

	return last, nil
}
