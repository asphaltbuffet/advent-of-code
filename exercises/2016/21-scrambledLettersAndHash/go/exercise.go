package exercises

import (
	"strconv"
	"strings"

	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

// Exercise for Advent of Code 2016 day 21.
type Exercise struct {
	common.BaseExercise
}

// seed returns the starting password to scramble: "abcde" for the 5-letter
// example, otherwise "abcdefgh" for the real puzzle. The example never
// references a letter past 'e' or an index past 4, so any such token marks the
// real input.
func seed(ops []string) string {
	for _, op := range ops {
		for _, f := range strings.Fields(op) {
			if len(f) == 1 && f[0] >= 'f' && f[0] <= 'z' {
				return "abcdefgh"
			}
			if n, err := strconv.Atoi(f); err == nil && n > 4 {
				return "abcdefgh"
			}
		}
	}
	return "abcde"
}

// indexOf returns the position of b in s, or -1.
func indexOf(s []byte, b byte) int {
	for i := range s {
		if s[i] == b {
			return i
		}
	}
	return -1
}

// rotate shifts s right by n (negative n rotates left).
func rotate(s []byte, n int) []byte {
	l := len(s)
	n = ((n % l) + l) % l
	return append(append([]byte{}, s[l-n:]...), s[:l-n]...)
}

// scramble applies all operations to pw in order.
func scramble(pw string, ops []string) string {
	s := []byte(pw)
	for _, op := range ops {
		f := strings.Fields(op)
		switch {
		case strings.HasPrefix(op, "swap position"):
			x, _ := strconv.Atoi(f[2])
			y, _ := strconv.Atoi(f[5])
			s[x], s[y] = s[y], s[x]
		case strings.HasPrefix(op, "swap letter"):
			x, y := indexOf(s, f[2][0]), indexOf(s, f[5][0])
			s[x], s[y] = s[y], s[x]
		case strings.HasPrefix(op, "rotate based"):
			i := indexOf(s, f[6][0])
			n := 1 + i
			if i >= 4 {
				n++
			}
			s = rotate(s, n)
		case strings.HasPrefix(op, "rotate left"):
			n, _ := strconv.Atoi(f[2])
			s = rotate(s, -n)
		case strings.HasPrefix(op, "rotate right"):
			n, _ := strconv.Atoi(f[2])
			s = rotate(s, n)
		case strings.HasPrefix(op, "reverse"):
			x, _ := strconv.Atoi(f[2])
			y, _ := strconv.Atoi(f[4])
			for x < y {
				s[x], s[y] = s[y], s[x]
				x++
				y--
			}
		case strings.HasPrefix(op, "move"):
			x, _ := strconv.Atoi(f[2])
			y, _ := strconv.Atoi(f[5])
			b := s[x]
			s = append(s[:x], s[x+1:]...)
			s = append(s[:y], append([]byte{b}, s[y:]...)...)
		}
	}
	return string(s)
}

func parseOps(instr string) []string {
	var ops []string
	for _, line := range strings.Split(strings.TrimSpace(instr), "\n") {
		if line = strings.TrimSpace(line); line != "" {
			ops = append(ops, line)
		}
	}
	return ops
}

// One scrambles the seed password using the listed operations.
func (e Exercise) One(instr string) (any, error) {
	ops := parseOps(instr)
	return scramble(seed(ops), ops), nil
}

// permutations generates every ordering of s, invoking yield for each.
func permutations(s []byte, k int, yield func([]byte)) {
	if k == len(s) {
		yield(s)
		return
	}
	for i := k; i < len(s); i++ {
		s[k], s[i] = s[i], s[k]
		permutations(s, k+1, yield)
		s[k], s[i] = s[i], s[k]
	}
}

// Two finds the password that scrambles to the target hash by brute-forcing
// every permutation of the seed letters.
func (e Exercise) Two(instr string) (any, error) {
	ops := parseOps(instr)
	target := "fbgdceah"
	if seed(ops) == "abcde" {
		target = "decab"
	}

	var found string
	permutations([]byte(seed(ops)), 0, func(s []byte) {
		if found == "" && scramble(string(s), ops) == target {
			found = string(s)
		}
	})
	return found, nil
}
