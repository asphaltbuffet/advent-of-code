// Package aoc22_01 contains the solution for day 1 of Advent of Code 2022.
package aoc22_01 //nolint:revive,stylecheck // I don't care about the package name

import (
	"sort"
	"strconv"
)

// D1P1 returns the solution for 2022 day 1 part 1
// answer: 70720
func D1P1(data []string) string {
	sum := 0
	calories := sort.IntSlice{}

	for i := 0; i < len(data); i++ {
		if data[i] == "" {
			calories = append(calories, sum)
			sum = 0
		} else {
			n, _ := strconv.Atoi(data[i])
			sum += n
		}
	}

	sort.Slice(calories, func(i, j int) bool {
		return calories[i] > calories[j]
	})

	return strconv.Itoa(calories[0])
}

// D1P2 returns the solution for 2022 day 1 part 2
// answer: 207148
func D1P2(data []string) string {
	sum := 0
	calories := sort.IntSlice{}

	for i := 0; i < len(data); i++ {
		if data[i] == "" {
			calories = append(calories, sum)
			sum = 0
		} else {
			n, _ := strconv.Atoi(data[i])
			sum += n
		}
	}

	sort.Slice(calories, func(i, j int) bool {
		return calories[i] > calories[j]
	})

	return strconv.Itoa(calories[0] + calories[1] + calories[2])
}
