package exercises

import (
	"crypto/md5"
	"encoding/hex"
	"strconv"
	"strings"

	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

// Exercise for Advent of Code 2016 day 14.
type Exercise struct {
	common.BaseExercise
}

// hasher produces (and caches) the hex hash for a given index, optionally with
// key stretching applied.
type hasher struct {
	salt    string
	stretch int
	cache   map[int]string
}

func newHasher(salt string, stretch int) *hasher {
	return &hasher{salt: salt, stretch: stretch, cache: map[int]string{}}
}

func (h *hasher) at(i int) string {
	if v, ok := h.cache[i]; ok {
		return v
	}
	sum := md5.Sum([]byte(h.salt + strconv.Itoa(i)))
	s := hex.EncodeToString(sum[:])
	for k := 0; k < h.stretch; k++ {
		sum = md5.Sum([]byte(s))
		s = hex.EncodeToString(sum[:])
	}
	h.cache[i] = s
	return s
}

// tripletChar returns the character of the first three-in-a-row run, or 0.
func tripletChar(s string) byte {
	for i := 0; i+2 < len(s); i++ {
		if s[i] == s[i+1] && s[i] == s[i+2] {
			return s[i]
		}
	}
	return 0
}

// hasFive reports whether s contains five of c in a row.
func hasFive(s string, c byte) bool {
	return strings.Contains(s, strings.Repeat(string(c), 5))
}

// sixtyFourthKey returns the index that produces the 64th valid key.
func sixtyFourthKey(salt string, stretch int) int {
	h := newHasher(salt, stretch)
	found := 0
	for i := 0; ; i++ {
		c := tripletChar(h.at(i))
		if c == 0 {
			continue
		}
		for j := i + 1; j <= i+1000; j++ {
			if hasFive(h.at(j), c) {
				found++
				if found == 64 {
					return i
				}
				break
			}
		}
	}
}

// One returns the answer to the first part of the exercise.
func (e Exercise) One(instr string) (any, error) {
	return sixtyFourthKey(strings.TrimSpace(instr), 0), nil
}

// Two returns the answer to the second part of the exercise.
func (e Exercise) Two(instr string) (any, error) {
	return sixtyFourthKey(strings.TrimSpace(instr), 2016), nil
}
