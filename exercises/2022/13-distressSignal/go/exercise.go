package exercises

import (
	"fmt"
	"sort"
	"strings"

	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

// Exercise for Advent of Code 2022 day 13
type Exercise struct {
	common.BaseExercise
}

// One returns the answer to the first part of the exercise.
func (c Exercise) One(instr string) (any, error) {
	data := strings.Split(instr, "\n")
	sum := 0

	for i := 0; i < len(data); i += 3 {
		first, err := ParsePacket(data[i])
		if err != nil {
			return nil, fmt.Errorf("parsing first packet: %w", err)
		}

		second, err := ParsePacket(data[i+1])
		if err != nil {
			return nil, fmt.Errorf("parsing second packet: %w", err)
		}

		if IsOrdered(first, second) {
			sum += (i / 3) + 1
		}
	}

	// fmt.Printf("valid pairs: %v\n", pairResult)

	return sum, nil
}

// Two returns the answer to the second part of the exercise.
func (c Exercise) Two(instr string) (any, error) {
	data := strings.Split(instr, "\n")
	var packets []any
	packets = append(packets, []any{[]any{2.}}, []any{[]any{6.}})

	for i := 0; i < len(data); i += 3 {
		first, err := ParsePacket(data[i])
		if err != nil {
			return nil, fmt.Errorf("error parsing first packet: %w", err)
		}

		second, err := ParsePacket(data[i+1])
		if err != nil {
			return nil, fmt.Errorf("parsing second packet: %w", err)
		}

		packets = append(packets, first, second)
	}

	sort.Slice(packets, func(i, j int) bool {
		return compare(packets[i], packets[j]) < 0
	})

	idx := 1

	for i, p := range packets {
		if fmt.Sprint(p) == "[[2]]" || fmt.Sprint(p) == "[[6]]" {
			idx *= i + 1
		}
	}

	return idx, nil
}
