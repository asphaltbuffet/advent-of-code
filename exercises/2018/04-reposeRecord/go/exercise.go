package exercises

import (
	"fmt"
	"sort"
	"strings"

	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

// Exercise for Advent of Code 2018 day 4.
type Exercise struct {
	common.BaseExercise
}

// sleepByGuard sorts the shuffled log lines (the timestamp format sorts
// lexically) and tallies, per guard, how many days they were asleep at each
// minute 0..59. asleep[guardID][minute] = number of days asleep at that minute.
func sleepByGuard(instr string) (map[int][60]int, error) {
	lines := strings.Split(strings.TrimSpace(instr), "\n")
	for i := range lines {
		lines[i] = strings.TrimSpace(lines[i])
	}
	sort.Strings(lines)

	asleep := make(map[int][60]int)
	guard, start := 0, 0

	for _, line := range lines {
		if line == "" {
			continue
		}

		switch {
		case strings.Contains(line, "Guard"):
			// "[...] Guard #NN begins shift" — the ID is the only # number.
			hash := strings.IndexByte(line, '#')
			if hash < 0 {
				return nil, fmt.Errorf("no guard id in %q", line)
			}
			id := 0
			for j := hash + 1; j < len(line) && line[j] >= '0' && line[j] <= '9'; j++ {
				id = id*10 + int(line[j]-'0')
			}
			guard = id

		case strings.Contains(line, "falls asleep"):
			start = minuteOf(line)

		case strings.Contains(line, "wakes up"):
			end := minuteOf(line)
			row := asleep[guard]
			for m := start; m < end; m++ {
				row[m]++
			}
			asleep[guard] = row
		}
	}

	return asleep, nil
}

// minuteOf extracts the two-digit minute from "[YYYY-MM-DD HH:MM] ...".
func minuteOf(line string) int {
	colon := strings.IndexByte(line, ':')
	if colon < 0 || colon+2 >= len(line) {
		return 0
	}
	return int(line[colon+1]-'0')*10 + int(line[colon+2]-'0')
}

// One returns the answer to the first part of the exercise.
// answer: 125444
func (e Exercise) One(instr string) (any, error) {
	asleep, err := sleepByGuard(instr)
	if err != nil {
		return nil, err
	}

	// Strategy 1: guard with the most total minutes asleep.
	bestGuard, bestTotal := 0, -1
	for id, mins := range asleep {
		total := 0
		for _, n := range mins {
			total += n
		}
		if total > bestTotal {
			bestGuard, bestTotal = id, total
		}
	}

	// The minute that guard was most often asleep.
	bestMin, bestCount := 0, -1
	for m, n := range asleep[bestGuard] {
		if n > bestCount {
			bestMin, bestCount = m, n
		}
	}

	return bestGuard * bestMin, nil
}

// Two returns the answer to the second part of the exercise.
// answer: 18325
func (e Exercise) Two(instr string) (any, error) {
	asleep, err := sleepByGuard(instr)
	if err != nil {
		return nil, err
	}

	// Strategy 2: the single (guard, minute) with the highest sleep count.
	bestGuard, bestMin, bestCount := 0, 0, -1
	for id, mins := range asleep {
		for m, n := range mins {
			if n > bestCount {
				bestGuard, bestMin, bestCount = id, m, n
			}
		}
	}

	return bestGuard * bestMin, nil
}
