// Package aoc22_13 contains the solution for day 13 of Advent of Code 2022.
package aoc22_13 //nolint:revive,stylecheck // I don't care about the package name

import (
	"fmt"
	"sort"
	"strconv"
)

// D13P1 returns the solution for 2022 day 13 part 1.
// https://adventofcode.com/2022/day/13
//
// answer:
func D13P1(data []string) string {
	sum := 0

	for i := 0; i < len(data); i += 3 {
		first, err := ParsePacket(data[i])
		if err != nil {
			return fmt.Sprintf("error parsing first packet: %v", err)
		}

		second, err := ParsePacket(data[i+1])
		if err != nil {
			return fmt.Sprintf("error parsing second packet: %v", err)
		}

		if IsOrdered(first, second) {
			sum += (i / 3) + 1
		}
	}

	// fmt.Printf("valid pairs: %v\n", pairResult)

	return strconv.Itoa(sum)
}

// D13P2 returns the solution for 2022 day 13 part 2.
// answer:
func D13P2(data []string) string {
	var packets []any
	packets = append(packets, []any{[]any{2.}}, []any{[]any{6.}})

	for i := 0; i < len(data); i += 3 {
		first, err := ParsePacket(data[i])
		if err != nil {
			return fmt.Sprintf("error parsing first packet: %v", err)
		}

		second, err := ParsePacket(data[i+1])
		if err != nil {
			return fmt.Sprintf("error parsing second packet: %v", err)
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

	return strconv.Itoa(idx)
}
