package exercises

import (
	"strconv"
	"strings"

	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

// Exercise for Advent of Code 2019 day 13.
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

// One returns the answer to the first part of the exercise.
// Runs the Intcode game program and counts tiles with tile_id == 2 (block tiles).
func (e Exercise) One(instr string) (any, error) {
	in := make(chan int, 1)
	out := make(chan int)

	go intcode(parseProgram(instr), in, out)

	blockCount := 0

	for {
		_, ok := <-out
		if !ok {
			break
		}

		_, ok = <-out
		if !ok {
			break
		}

		tileID, ok := <-out
		if !ok {
			break
		}

		if tileID == 2 {
			blockCount++
		}
	}

	return blockCount, nil
}

// Two returns the answer to the second part of the exercise.
func (e Exercise) Two(instr string) (any, error) { //nolint:gocognit,funlen
	prog := parseProgram(instr)
	prog.set(0, 2)

	ip, relBase := 0, 0
	var ballX, paddleX, score int
	outBuf := make([]int, 0, 3)

	param := func(offset, mode int) int {
		raw := prog.get(ip + offset)
		switch mode {
		case 0:
			return prog.get(raw)
		case 1:
			return raw
		case 2:
			return prog.get(relBase + raw)
		}
		panic("unknown mode")
	}

	dest := func(offset, mode int) int {
		raw := prog.get(ip + offset)
		switch mode {
		case 0:
			return raw
		case 2:
			return relBase + raw
		}
		panic("unknown dest mode")
	}

	for {
		instr := prog.get(ip)
		op := instr % 100
		m1 := (instr / 100) % 10
		m2 := (instr / 1000) % 10
		m3 := (instr / 10000) % 10

		switch op {
		case 1:
			prog.set(dest(3, m3), param(1, m1)+param(2, m2))
			ip += 4
		case 2:
			prog.set(dest(3, m3), param(1, m1)*param(2, m2))
			ip += 4
		case 3:
			// Provide joystick: track ball toward paddle.
			j := 0
			if ballX < paddleX {
				j = -1
			} else if ballX > paddleX {
				j = 1
			}
			prog.set(dest(1, m1), j)
			ip += 2
		case 4:
			outBuf = append(outBuf, param(1, m1))
			if len(outBuf) == 3 {
				x, y, val := outBuf[0], outBuf[1], outBuf[2]
				outBuf = outBuf[:0]
				if x == -1 && y == 0 {
					score = val
				} else {
					switch val {
					case 3:
						paddleX = x
					case 4:
						ballX = x
					}
				}
			}
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
			prog.set(dest(3, m3), v)
			ip += 4
		case 8:
			v := 0
			if param(1, m1) == param(2, m2) {
				v = 1
			}
			prog.set(dest(3, m3), v)
			ip += 4
		case 9:
			relBase += param(1, m1)
			ip += 2
		case 99:
			return score, nil
		}
	}
}
