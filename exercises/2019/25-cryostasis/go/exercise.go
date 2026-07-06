package exercises

import (
	"errors"
	"fmt"
	"maps"
	"strconv"
	"strings"

	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

// Exercise for Advent of Code 2019 day 25.
type Exercise struct {
	common.BaseExercise
}

// mem is a map-backed unbounded memory for Intcode.
type mem map[int]int

func (m mem) get(addr int) int  { return m[addr] }
func (m mem) set(addr, val int) { m[addr] = val }

func copyMem(src mem) mem {
	dst := make(mem, len(src))
	maps.Copy(dst, src)
	return dst
}

func parseProgram(input string) mem {
	m := make(mem)
	for i, s := range strings.Split(strings.TrimSpace(input), ",") {
		v, _ := strconv.Atoi(s)
		m[i] = v
	}
	return m
}

// icvm holds the state of an Intcode computer.
type icvm struct {
	m       mem
	ip      int
	relBase int
	inQueue []int
	halted  bool
}

func newICVM(program mem) *icvm {
	return &icvm{m: copyMem(program)}
}

// sendInput queues ASCII input with a trailing newline.
func (v *icvm) sendInput(line string) {
	for _, c := range line {
		v.inQueue = append(v.inQueue, int(c))
	}
	v.inQueue = append(v.inQueue, '\n')
}

// param returns the value of an instruction parameter given its mode.
func (v *icvm) param(offset, mode int) int {
	raw := v.m.get(v.ip + offset)
	switch mode {
	case 0:
		return v.m.get(raw)
	case 1:
		return raw
	case 2:
		return v.m.get(v.relBase + raw)
	}
	panic("unknown param mode")
}

// dest returns the write address for an instruction parameter given its mode.
func (v *icvm) dest(offset, mode int) int {
	raw := v.m.get(v.ip + offset)
	switch mode {
	case 0:
		return raw
	case 2:
		return v.relBase + raw
	}
	panic("unknown dest mode")
}

// run runs until the VM blocks on empty input or halts, collecting output.
//
//nolint:funlen // Intcode opcode dispatch is inherently verbose
func (v *icvm) run() string {
	var sb strings.Builder

	for !v.halted {
		instr := v.m.get(v.ip)
		op := instr % 100
		m1 := (instr / 100) % 10
		m2 := (instr / 1000) % 10
		m3 := (instr / 10000) % 10

		switch op {
		case 1:
			v.m.set(v.dest(3, m3), v.param(1, m1)+v.param(2, m2))
			v.ip += 4
		case 2:
			v.m.set(v.dest(3, m3), v.param(1, m1)*v.param(2, m2))
			v.ip += 4
		case 3:
			if len(v.inQueue) == 0 {
				return sb.String()
			}
			val := v.inQueue[0]
			v.inQueue = v.inQueue[1:]
			v.m.set(v.dest(1, m1), val)
			v.ip += 2
		case 4:
			sb.WriteRune(rune(v.param(1, m1)))
			v.ip += 2
		case 5:
			if v.param(1, m1) != 0 {
				v.ip = v.param(2, m2)
			} else {
				v.ip += 3
			}
		case 6:
			if v.param(1, m1) == 0 {
				v.ip = v.param(2, m2)
			} else {
				v.ip += 3
			}
		case 7:
			val := 0
			if v.param(1, m1) < v.param(2, m2) {
				val = 1
			}
			v.m.set(v.dest(3, m3), val)
			v.ip += 4
		case 8:
			val := 0
			if v.param(1, m1) == v.param(2, m2) {
				val = 1
			}
			v.m.set(v.dest(3, m3), val)
			v.ip += 4
		case 9:
			v.relBase += v.param(1, m1)
			v.ip += 2
		case 99:
			v.halted = true
			return sb.String()
		default:
			panic(fmt.Sprintf("unknown opcode %d at ip %d", op, v.ip))
		}
	}

	return sb.String()
}

func sendCmd(v *icvm, cmd string) string {
	v.sendInput(cmd)
	return v.run()
}

// dangerousItems cause death or infinite loops — never pick them up.
var dangerousItems = map[string]bool{
	"infinite loop":       true,
	"giant electromagnet": true,
	"escape pod":          true,
	"photons":             true,
	"molten lava":         true,
}

const (
	dirNorth = "north"
	dirSouth = "south"
	dirEast  = "east"
	dirWest  = "west"
)

func opposite(d string) string {
	switch d {
	case dirNorth:
		return dirSouth
	case dirSouth:
		return dirNorth
	case dirEast:
		return dirWest
	case dirWest:
		return dirEast
	}
	return ""
}

// parseOutput extracts room name, cardinal exits, and items from ASCII game output.
func parseOutput(output string) (string, []string, []string) {
	var roomName string
	var exits, items []string
	section := ""

	for line := range strings.SplitSeq(output, "\n") {
		line = strings.TrimSpace(line)
		switch {
		case strings.HasPrefix(line, "== ") && strings.HasSuffix(line, " =="):
			roomName = line[3 : len(line)-3]
			section = ""
		case line == "Doors here lead:":
			section = "exits"
		case line == "Items here:":
			section = "items"
		case line == "Command?":
			section = ""
		case strings.HasPrefix(line, "- "):
			val := line[2:]
			switch section {
			case "exits":
				if val == dirNorth || val == dirSouth || val == dirEast || val == dirWest {
					exits = append(exits, val)
				}
			case "items":
				items = append(items, val)
			}
		}
	}
	return roomName, exits, items
}

// extractPassword pulls the first long digit-string from the success output.
func extractPassword(output string) string {
	for w := range strings.FieldsSeq(output) {
		w = strings.Trim(w, ".,!?;:")
		if len(w) < 6 {
			continue
		}
		allDigits := true
		for _, c := range w {
			if c < '0' || c > '9' {
				allDigits = false
				break
			}
		}
		if allDigits {
			return w
		}
	}
	return ""
}

// checkpointInfo records the security checkpoint location.
type checkpointInfo struct {
	checkDir string   // direction through the pressure-sensitive floor
	pathTo   []string // directions from start to reach Security Checkpoint
}

// dfsFrame is the DFS backtracking stack entry.
type dfsFrame struct{ back string }

// dfsState holds the mutable state threaded through the DFS.
type dfsState struct {
	v         *icvm
	visited   map[string]bool
	inventory []string
	path      []dfsFrame
	cpInfo    *checkpointInfo
}

// visit explores roomName (already confirmed unvisited), collects items, and recurses.
//
//nolint:cyclop // DFS with checkpoint detection and backtracking is inherently branchy
func (s *dfsState) visit(roomName string, exits, items []string) {
	s.visited[roomName] = true

	for _, item := range items {
		if !dangerousItems[item] {
			sendCmd(s.v, "take "+item)
			s.inventory = append(s.inventory, item)
		}
	}

	if roomName == "Security Checkpoint" {
		backDir := ""
		if len(s.path) > 0 {
			backDir = s.path[len(s.path)-1].back
		}
		var fwd string
		for _, ex := range exits {
			if ex != backDir {
				fwd = ex
				break
			}
		}
		var pathDirs []string
		for _, f := range s.path {
			pathDirs = append(pathDirs, opposite(f.back))
		}
		s.cpInfo = &checkpointInfo{checkDir: fwd, pathTo: pathDirs}
		return
	}

	for _, dir := range exits {
		if len(s.path) > 0 && dir == s.path[len(s.path)-1].back {
			continue
		}
		out := sendCmd(s.v, dir)
		nextRoom, nextExits, nextItems := parseOutput(out)
		if s.visited[nextRoom] {
			sendCmd(s.v, opposite(dir))
			continue
		}
		s.path = append(s.path, dfsFrame{back: opposite(dir)})
		s.visit(nextRoom, nextExits, nextItems)
		s.path = s.path[:len(s.path)-1]
		sendCmd(s.v, opposite(dir))
	}
}

// crackCheckpoint tries all 2^N item subsets to pass the weight sensor.
func crackCheckpoint(v *icvm, checkDir string, inventory []string) (string, error) {
	n := len(inventory)
	for _, item := range inventory {
		sendCmd(v, "drop "+item)
	}
	for mask := range 1 << n {
		var held []string
		for i, item := range inventory {
			if mask&(1<<i) != 0 {
				held = append(held, item)
				sendCmd(v, "take "+item)
			}
		}
		out := sendCmd(v, checkDir)
		if !strings.Contains(out, "lighter") && !strings.Contains(out, "heavier") {
			if pass := extractPassword(out); pass != "" {
				return pass, nil
			}
		}
		for _, item := range held {
			sendCmd(v, "drop "+item)
		}
	}
	return "", errors.New("no valid item combination found")
}

// explorer navigates the ship, collects items, and cracks the checkpoint.
func explorer(program mem) (string, error) {
	v := newICVM(program)
	initial := v.run()

	s := &dfsState{
		v:       v,
		visited: make(map[string]bool),
	}

	roomName, exits, items := parseOutput(initial)
	s.visit(roomName, exits, items)

	if s.cpInfo == nil {
		return "", errors.New("never reached Security Checkpoint")
	}

	// Navigate back to Security Checkpoint (DFS ends at start).
	for _, dir := range s.cpInfo.pathTo {
		sendCmd(v, dir)
	}

	return crackCheckpoint(v, s.cpInfo.checkDir, s.inventory)
}

// One returns the answer to the first part of the exercise.
func (e Exercise) One(instr string) (any, error) {
	program := parseProgram(instr)
	return explorer(program)
}

// Two returns the answer to the second part of the exercise.
func (e Exercise) Two(_ string) (any, error) {
	return "", nil
}
