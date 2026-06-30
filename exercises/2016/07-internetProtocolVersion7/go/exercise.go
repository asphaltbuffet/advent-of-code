package exercises

import (
	"strings"

	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

// Exercise for Advent of Code 2016 day 7.
type Exercise struct {
	common.BaseExercise
}

// segments splits an IPv7 address into supernet (outside brackets) and hypernet
// (inside brackets) sequences. Splitting on '[' and ']' alternates the two.
func segments(ip string) (supernets, hypernets []string) {
	depth := 0
	var cur strings.Builder
	flush := func() {
		if cur.Len() > 0 {
			if depth == 0 {
				supernets = append(supernets, cur.String())
			} else {
				hypernets = append(hypernets, cur.String())
			}
			cur.Reset()
		}
	}
	for i := 0; i < len(ip); i++ {
		switch ip[i] {
		case '[':
			flush()
			depth = 1
		case ']':
			flush()
			depth = 0
		default:
			cur.WriteByte(ip[i])
		}
	}
	flush()
	return supernets, hypernets
}

// hasABBA reports whether s contains a four-char ABBA pattern.
func hasABBA(s string) bool {
	for i := 0; i+3 < len(s); i++ {
		if s[i] == s[i+3] && s[i+1] == s[i+2] && s[i] != s[i+1] {
			return true
		}
	}
	return false
}

// abas returns every ABA pattern (xyx, x != y) found in s.
func abas(s string) []string {
	var out []string
	for i := 0; i+2 < len(s); i++ {
		if s[i] == s[i+2] && s[i] != s[i+1] {
			out = append(out, s[i:i+3])
		}
	}
	return out
}

// One returns the answer to the first part of the exercise.
func (e Exercise) One(instr string) (any, error) {
	count := 0
	for _, ip := range strings.Fields(instr) {
		super, hyper := segments(ip)
		tls := false
		for _, s := range super {
			if hasABBA(s) {
				tls = true
			}
		}
		for _, h := range hyper {
			if hasABBA(h) {
				tls = false
				break
			}
		}
		if tls {
			count++
		}
	}
	return count, nil
}

// Two returns the answer to the second part of the exercise.
func (e Exercise) Two(instr string) (any, error) {
	count := 0
	for _, ip := range strings.Fields(instr) {
		super, hyper := segments(ip)
		ssl := false
		for _, s := range super {
			for _, aba := range abas(s) {
				bab := string([]byte{aba[1], aba[0], aba[1]})
				for _, h := range hyper {
					if strings.Contains(h, bab) {
						ssl = true
					}
				}
			}
		}
		if ssl {
			count++
		}
	}
	return count, nil
}
