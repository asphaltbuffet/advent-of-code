package exercises

import (
	"maps"
	"strconv"
	"strings"

	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

// Exercise for Advent of Code 2019 day 23.
type Exercise struct {
	common.BaseExercise
}

// mem is a map-backed unbounded memory for Intcode.
type mem map[int]int

func (m mem) get(addr int) int  { return m[addr] }
func (m mem) set(addr, val int) { m[addr] = val }

// copyMem returns a deep copy of a mem map.
func copyMem(src mem) mem {
	dst := make(mem, len(src))
	maps.Copy(dst, src)
	return dst
}

// parseProgram parses a comma-separated Intcode program into a mem map.
func parseProgram(input string) mem {
	m := make(mem)
	for i, s := range strings.Split(strings.TrimSpace(input), ",") {
		v, _ := strconv.Atoi(s)
		m[i] = v
	}
	return m
}

// nicState holds the state for one Intcode NIC.
type nicState struct {
	m       mem
	ip      int
	relBase int
	inQueue []int
	outBuf  []int
	halted  bool
}

// newNIC creates a NIC initialized with the given program and address.
func newNIC(program mem, addr int) *nicState {
	return &nicState{
		m:       copyMem(program),
		inQueue: []int{addr},
	}
}

// step executes one Intcode instruction on the NIC.
// Returns true if the NIC produced a complete packet (3 output values),
// and false if the NIC is blocked on input (no input available).
// halted is set if opcode 99 is encountered.
//
//nolint:funlen // Intcode VM opcode dispatch is inherently complex
func (n *nicState) step() bool {
	if n.halted {
		return true
	}

	m := n.m
	ip := &n.ip
	relBase := &n.relBase

	param := func(offset, mode int) int {
		raw := m.get(*ip + offset)
		switch mode {
		case 0:
			return m.get(raw)
		case 1:
			return raw
		case 2:
			return m.get(*relBase + raw)
		}
		panic("unknown mode")
	}

	dest := func(offset, mode int) int {
		raw := m.get(*ip + offset)
		switch mode {
		case 0:
			return raw
		case 2:
			return *relBase + raw
		}
		panic("unknown dest mode")
	}

	instr := m.get(*ip)
	op := instr % 100
	m1 := (instr / 100) % 10
	m2 := (instr / 1000) % 10
	m3 := (instr / 10000) % 10

	switch op {
	case 1:
		m.set(dest(3, m3), param(1, m1)+param(2, m2))
		*ip += 4
	case 2:
		m.set(dest(3, m3), param(1, m1)*param(2, m2))
		*ip += 4
	case 3:
		if len(n.inQueue) == 0 {
			return true // blocked on input
		}
		val := n.inQueue[0]
		n.inQueue = n.inQueue[1:]
		m.set(dest(1, m1), val)
		*ip += 2
	case 4:
		n.outBuf = append(n.outBuf, param(1, m1))
		*ip += 2
	case 5:
		if param(1, m1) != 0 {
			*ip = param(2, m2)
		} else {
			*ip += 3
		}
	case 6:
		if param(1, m1) == 0 {
			*ip = param(2, m2)
		} else {
			*ip += 3
		}
	case 7:
		v := 0
		if param(1, m1) < param(2, m2) {
			v = 1
		}
		m.set(dest(3, m3), v)
		*ip += 4
	case 8:
		v := 0
		if param(1, m1) == param(2, m2) {
			v = 1
		}
		m.set(dest(3, m3), v)
		*ip += 4
	case 9:
		*relBase += param(1, m1)
		*ip += 2
	case 99:
		n.halted = true
		return true
	}

	return false
}

// runNetwork runs 50 NICs and returns the first packet sent to address 255.
//
//nolint:gocognit // network simulation with packet routing is inherently multi-branch
func runNetwork(program mem) int {
	const numNICs = 50

	nics := make([]*nicState, numNICs)
	for i := range nics {
		nics[i] = newNIC(program, i)
	}

	for {
		for i := range numNICs {
			nic := nics[i]

			// Run until blocked on input, delivering complete packets along the way.
			for {
				blocked := nic.step()

				// Deliver complete packets from output buffer.
				for len(nic.outBuf) >= 3 {
					dest := nic.outBuf[0]
					x := nic.outBuf[1]
					y := nic.outBuf[2]
					nic.outBuf = nic.outBuf[3:]

					if dest == 255 {
						return y
					}

					if dest >= 0 && dest < numNICs {
						nics[dest].inQueue = append(nics[dest].inQueue, x, y)
					}
				}

				if blocked {
					// Blocked on input with nothing queued — feed -1.
					if len(nic.inQueue) == 0 {
						nic.inQueue = append(nic.inQueue, -1)
					}
					break
				}
			}
		}
	}
}

// runNetworkNAT runs the 50-NIC network with a NAT at address 255.
// It returns the first Y value the NAT delivers to NIC 0 twice in a row.
func drainOutBuf(nic *nicState, nics []*nicState, nat *[2]int, natHasPacket *bool) bool {
	anyOutput := false
	for len(nic.outBuf) >= 3 {
		dest := nic.outBuf[0]
		x := nic.outBuf[1]
		y := nic.outBuf[2]
		nic.outBuf = nic.outBuf[3:]
		anyOutput = true

		if dest == 255 {
			nat[0], nat[1] = x, y
			*natHasPacket = true
		} else if dest >= 0 && dest < len(nics) {
			nics[dest].inQueue = append(nics[dest].inQueue, x, y)
		}
	}
	return anyOutput
}

func stepUntilBlocked(nic *nicState, nics []*nicState, nat *[2]int, natHasPacket *bool) bool {
	anyOutput := false
	for {
		blocked := nic.step()
		if drainOutBuf(nic, nics, nat, natHasPacket) {
			anyOutput = true
		}
		if blocked {
			break
		}
	}
	return anyOutput
}

func runNetworkNAT(program mem) int {
	const numNICs = 50

	nics := make([]*nicState, numNICs)
	for i := range nics {
		nics[i] = newNIC(program, i)
	}

	var nat [2]int
	natHasPacket := false
	lastNATSentY := -1 // sentinel: no Y sent yet

	for {
		anyOutput := false
		allIdle := true // assume idle until proven otherwise

		for i := range numNICs {
			nic := nics[i]

			if stepUntilBlocked(nic, nics, &nat, &natHasPacket) {
				anyOutput = true
			}

			// If this NIC has queued input, the network is not idle.
			if len(nic.inQueue) > 0 {
				allIdle = false
			} else {
				// Queue is empty — feed -1 so the NIC doesn't stall forever.
				nic.inQueue = append(nic.inQueue, -1)
			}
		}

		// The network is idle when no NIC produced output this round AND
		// every NIC had an empty input queue (was blocked on -1).
		if !anyOutput && allIdle && natHasPacket {
			if nat[1] == lastNATSentY {
				return nat[1]
			}
			lastNATSentY = nat[1]
			nics[0].inQueue = append(nics[0].inQueue, nat[0], nat[1])
		}
	}
}

// One returns the answer to the first part of the exercise.
// Finds the Y value of the first packet sent to address 255.
func (e Exercise) One(instr string) (any, error) {
	program := parseProgram(instr)
	return runNetwork(program), nil
}

// Two returns the answer to the second part of the exercise.
func (e Exercise) Two(instr string) (any, error) {
	program := parseProgram(instr)
	return runNetworkNAT(program), nil
}
