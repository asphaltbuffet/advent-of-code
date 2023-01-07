// Package aoc22_25 contains the solution for day 25 of Advent of Code 2022.
package aoc22_25 //nolint:revive,stylecheck // I don't care about the package name

// D25P1 returns the solution for 2022 day 25 part 1.
//
// https://adventofcode.com/2022/day/25
//
// answer: 2==221=-002=0-02-000
func D25P1(data []string) string {
	sum := 0

	for _, line := range data {
		sum += Decode(line)
	}

	return Encode(sum)
}

// D25P2 returns the solution for 2022 day 25 part 2.
// answer:
func D25P2(data []string) string {
	return ""
}
