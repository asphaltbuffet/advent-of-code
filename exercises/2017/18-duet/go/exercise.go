package exercises

import (
	"strconv"
	"strings"

	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

// Exercise for Advent of Code 2017 day 18.
type Exercise struct {
	common.BaseExercise
}

func parseProgram(instr string) [][]string {
	var prog [][]string
	for _, line := range strings.Split(strings.TrimSpace(instr), "\n") {
		if line = strings.TrimSpace(line); line != "" {
			prog = append(prog, strings.Fields(line))
		}
	}
	return prog
}

// One plays the program as a solo: snd records the last frequency, and the
// first rcv with a non-zero argument recovers (and returns) it.
func (e Exercise) One(instr string) (any, error) {
	prog := parseProgram(instr)
	regs := map[string]int64{}
	val := func(x string) int64 {
		if n, err := strconv.ParseInt(x, 10, 64); err == nil {
			return n
		}
		return regs[x]
	}

	var lastSound int64
	for ip := 0; ip >= 0 && ip < len(prog); {
		op := prog[ip]
		switch op[0] {
		case "snd":
			lastSound = val(op[1])
		case "set":
			regs[op[1]] = val(op[2])
		case "add":
			regs[op[1]] += val(op[2])
		case "mul":
			regs[op[1]] *= val(op[2])
		case "mod":
			regs[op[1]] %= val(op[2])
		case "rcv":
			if val(op[1]) != 0 {
				return lastSound, nil
			}
		case "jgz":
			if val(op[1]) > 0 {
				ip += int(val(op[2]))
				continue
			}
		}
		ip++
	}
	return lastSound, nil
}

// vm is one of the two duet programs in Part Two.
type vm struct {
	prog    [][]string
	regs    map[string]int64
	ip      int
	inbox   []int64
	outbox  *[]int64 // the other program's inbox
	sends   int
	blocked bool
}

func (m *vm) val(x string) int64 {
	if n, err := strconv.ParseInt(x, 10, 64); err == nil {
		return n
	}
	return m.regs[x]
}

// run executes until the program blocks on an empty queue or halts. It returns
// once no more progress can be made without new input.
func (m *vm) run() {
	for m.ip >= 0 && m.ip < len(m.prog) {
		op := m.prog[m.ip]
		switch op[0] {
		case "snd":
			*m.outbox = append(*m.outbox, m.val(op[1]))
			m.sends++
		case "set":
			m.regs[op[1]] = m.val(op[2])
		case "add":
			m.regs[op[1]] += m.val(op[2])
		case "mul":
			m.regs[op[1]] *= m.val(op[2])
		case "mod":
			m.regs[op[1]] %= m.val(op[2])
		case "rcv":
			if len(m.inbox) == 0 {
				m.blocked = true
				return // wait for the other program to send
			}
			m.regs[op[1]] = m.inbox[0]
			m.inbox = m.inbox[1:]
		case "jgz":
			if m.val(op[1]) > 0 {
				m.ip += int(m.val(op[2]))
				continue
			}
		}
		m.ip++
	}
	m.blocked = true // ran off the end
}

// Two runs the two programs concurrently and returns how many values program 1
// sends before the pair deadlocks.
func (e Exercise) Two(instr string) (any, error) {
	prog := parseProgram(instr)

	p0 := &vm{prog: prog, regs: map[string]int64{"p": 0}}
	p1 := &vm{prog: prog, regs: map[string]int64{"p": 1}}
	p0.outbox = &p1.inbox
	p1.outbox = &p0.inbox

	for {
		p0.blocked, p1.blocked = false, false
		p0.run()
		p1.run()
		// Deadlock: both blocked and neither has pending input to consume.
		if p0.blocked && p1.blocked && len(p0.inbox) == 0 && len(p1.inbox) == 0 {
			break
		}
	}
	return p1.sends, nil
}
