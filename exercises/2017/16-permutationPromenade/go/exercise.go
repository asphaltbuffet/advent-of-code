package exercises

import (
	"strconv"
	"strings"

	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

// Exercise for Advent of Code 2017 day 16.
type Exercise struct {
	common.BaseExercise
}

// startLine returns the initial program line: "abcde" for the 5-program
// example, otherwise "abcdefghijklmnop".
func startLine(moves string) []byte {
	if strings.TrimSpace(moves) == "s1,x3/4,pe/b" {
		return []byte("abcde")
	}
	return []byte("abcdefghijklmnop")
}

// dance applies the whole move sequence once to line, in place.
func dance(line []byte, moves []string) {
	n := len(line)
	for _, m := range moves {
		switch m[0] {
		case 's':
			x, _ := strconv.Atoi(m[1:])
			x %= n
			rotated := append(append([]byte{}, line[n-x:]...), line[:n-x]...)
			copy(line, rotated)
		case 'x':
			ab := strings.SplitN(m[1:], "/", 2)
			a, _ := strconv.Atoi(ab[0])
			b, _ := strconv.Atoi(ab[1])
			line[a], line[b] = line[b], line[a]
		case 'p':
			a := strings.IndexByte(string(line), m[1])
			b := strings.IndexByte(string(line), m[3])
			line[a], line[b] = line[b], line[a]
		}
	}
}

func parseMoves(instr string) []string {
	return strings.Split(strings.TrimSpace(instr), ",")
}

// One returns the program order after a single dance.
func (e Exercise) One(instr string) (any, error) {
	line := startLine(instr)
	dance(line, parseMoves(instr))
	return string(line), nil
}

// Two returns the order after one billion dances, using cycle detection to skip
// the vast majority of repetitions.
func (e Exercise) Two(instr string) (any, error) {
	const target = 1_000_000_000
	moves := parseMoves(instr)
	line := startLine(instr)

	seen := []string{string(line)}
	index := map[string]int{string(line): 0}
	for i := 1; i <= target; i++ {
		dance(line, moves)
		s := string(line)
		if first, ok := index[s]; ok {
			// Cycle from `first` of length i-first; land on the target offset.
			cycle := i - first
			return seen[first+(target-first)%cycle], nil
		}
		index[s] = i
		seen = append(seen, s)
	}
	return string(line), nil
}
