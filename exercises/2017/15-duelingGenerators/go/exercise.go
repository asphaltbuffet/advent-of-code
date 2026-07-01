package exercises

import (
	"regexp"
	"strconv"

	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

// Exercise for Advent of Code 2017 day 15.
type Exercise struct {
	common.BaseExercise
}

const (
	factorA = 16807
	factorB = 48271
	modulus = 2147483647
)

var intRe = regexp.MustCompile(`\d+`)

// startValues extracts the two generator seeds by scanning integers, so the
// surrounding "Generator A starts with N" text doesn't matter.
func startValues(instr string) (uint64, uint64) {
	nums := intRe.FindAllString(instr, -1)
	a, _ := strconv.Atoi(nums[0])
	b, _ := strconv.Atoi(nums[1])
	return uint64(a), uint64(b)
}

// One counts lowest-16-bit matches over 40 million generated pairs.
func (e Exercise) One(instr string) (any, error) {
	a, b := startValues(instr)
	matches := 0
	for i := 0; i < 40_000_000; i++ {
		a = a * factorA % modulus
		b = b * factorB % modulus
		if a&0xffff == b&0xffff {
			matches++
		}
	}
	return matches, nil
}

// Two counts matches over 5 million pairs where A yields only multiples of 4 and
// B only multiples of 8.
func (e Exercise) Two(instr string) (any, error) {
	a, b := startValues(instr)
	matches := 0
	for i := 0; i < 5_000_000; i++ {
		for {
			a = a * factorA % modulus
			if a&3 == 0 {
				break
			}
		}
		for {
			b = b * factorB % modulus
			if b&7 == 0 {
				break
			}
		}
		if a&0xffff == b&0xffff {
			matches++
		}
	}
	return matches, nil
}
