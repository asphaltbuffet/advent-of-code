package exercises

import (
	"fmt"
	"math/bits"
	"strings"

	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

// Exercise for Advent of Code 2021 day 8.
type Exercise struct {
	common.BaseExercise
}

type entry struct {
	patterns []uint8 // 10 signal patterns as 7-bit segment masks
	output   []uint8 // 4 output digits as 7-bit segment masks
}

// segMask turns a wire string like "acf" into a bitmask over segments a..g.
func segMask(s string) uint8 {
	var m uint8
	for _, c := range s {
		m |= 1 << (c - 'a')
	}
	return m
}

func parse(instr string) ([]entry, error) {
	var entries []entry

	for _, line := range strings.Split(strings.TrimSpace(instr), "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		left, right, ok := strings.Cut(line, " | ")
		if !ok {
			return nil, fmt.Errorf("line missing ` | ` separator: %q", line)
		}

		e := entry{}
		for _, p := range strings.Fields(left) {
			e.patterns = append(e.patterns, segMask(p))
		}
		for _, p := range strings.Fields(right) {
			e.output = append(e.output, segMask(p))
		}

		entries = append(entries, e)
	}

	return entries, nil
}

// One returns the answer to the first part of the exercise.
func (e Exercise) One(instr string) (any, error) {
	entries, err := parse(instr)
	if err != nil {
		return nil, fmt.Errorf("parsing input: %w", err)
	}

	// 1, 4, 7, 8 have unique segment counts (2, 4, 3, 7).
	count := 0
	for _, ent := range entries {
		for _, o := range ent.output {
			switch bits.OnesCount8(o) {
			case 2, 3, 4, 7:
				count++
			}
		}
	}

	return fmt.Sprintf("%d", count), nil
}

// decode works out the mask for each digit 0..9 on one entry, then reads its
// four output digits as a number.
func decode(ent entry) int {
	var digit [10]uint8

	// Anchors by unique length.
	var fives, sixes []uint8
	for _, p := range ent.patterns {
		switch bits.OnesCount8(p) {
		case 2:
			digit[1] = p
		case 3:
			digit[7] = p
		case 4:
			digit[4] = p
		case 7:
			digit[8] = p
		case 5:
			fives = append(fives, p)
		case 6:
			sixes = append(sixes, p)
		}
	}

	// 6-segment digits: 9 contains all of 4; 0 contains all of 1 (but isn't 9);
	// 6 is the remaining one.
	for _, p := range sixes {
		switch {
		case p&digit[4] == digit[4]:
			digit[9] = p
		case p&digit[1] == digit[1]:
			digit[0] = p
		default:
			digit[6] = p
		}
	}

	// 5-segment digits: 3 contains all of 1; 5 is a subset of 6; 2 is the rest.
	for _, p := range fives {
		switch {
		case p&digit[1] == digit[1]:
			digit[3] = p
		case p&digit[6] == p:
			digit[5] = p
		default:
			digit[2] = p
		}
	}

	// Reverse map mask -> value for output lookup.
	value := map[uint8]int{}
	for v, m := range digit {
		value[m] = v
	}

	n := 0
	for _, o := range ent.output {
		n = n*10 + value[o]
	}

	return n
}

// Two returns the answer to the second part of the exercise.
func (e Exercise) Two(instr string) (any, error) {
	entries, err := parse(instr)
	if err != nil {
		return nil, fmt.Errorf("parsing input: %w", err)
	}

	sum := 0
	for _, ent := range entries {
		sum += decode(ent)
	}

	return fmt.Sprintf("%d", sum), nil
}
