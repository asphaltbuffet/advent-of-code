package exercises

import (
	"strconv"
	"strings"

	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

// Exercise for Advent of Code 2015 day 16.
type Exercise struct {
	common.BaseExercise
}

// target is the MFCSAM ticker-tape reading from the problem statement (not the
// input): the exact compound counts of the Aunt Sue we are looking for.
var target = map[string]int{
	"children": 3, "cats": 7, "samoyeds": 2, "pomeranians": 3, "akitas": 0,
	"vizslas": 0, "goldfish": 5, "trees": 3, "cars": 2, "perfumes": 1,
}

// parse turns each "Sue N: compound: x, ..." line into (number, observations).
func parse(instr string) []map[string]int {
	var sues []map[string]int

	for _, line := range strings.Split(strings.TrimSpace(instr), "\n") {
		// Drop the "Sue N:" prefix; the slice index gives the number.
		_, rest, _ := strings.Cut(line, ": ")
		obs := map[string]int{}
		for _, pair := range strings.Split(rest, ", ") {
			k, v, _ := strings.Cut(pair, ": ")
			n, _ := strconv.Atoi(v)
			obs[k] = n
		}
		sues = append(sues, obs)
	}

	return sues
}

// find returns the 1-based number of the first Sue whose observations are all
// consistent with target. ok(compound, got) reports whether the observed count
// for compound is consistent with the target reading.
func find(sues []map[string]int, ok func(compound string, got int) bool) int {
	for i, obs := range sues {
		match := true
		for k, got := range obs {
			if !ok(k, got) {
				match = false
				break
			}
		}
		if match {
			return i + 1
		}
	}
	return -1
}

// One: every observed compound must equal the target exactly.
func (e Exercise) One(instr string) (any, error) {
	return find(parse(instr), func(k string, got int) bool {
		return got == target[k]
	}), nil
}

// Two: the readings are imprecise. cats/trees are "greater than" the target
// (observed must exceed it); pomeranians/goldfish are "fewer than"; all other
// compounds still compare exactly.
func (e Exercise) Two(instr string) (any, error) {
	return find(parse(instr), func(k string, got int) bool {
		switch k {
		case "cats", "trees":
			return got > target[k]
		case "pomeranians", "goldfish":
			return got < target[k]
		default:
			return got == target[k]
		}
	}), nil
}
