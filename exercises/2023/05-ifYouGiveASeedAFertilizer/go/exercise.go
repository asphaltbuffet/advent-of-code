package exercises

import (
	"fmt"
	"math"
	"slices"
	"strings"

	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

// Exercise for Advent of Code 2023 day 5.
type Exercise struct {
	common.BaseExercise
}

// One returns the answer to the first part of the exercise.
func (e Exercise) One(instr string) (any, error) {
	sections := strings.Split(instr, "\n\n")

	seeds := parseSeeds(sections[0])
	fmt.Println("seed count: ", len(seeds))

	maps := parseAllMaps(sections[1:])

	locations := getLocations(maps, seeds)

	min := slices.Min(locations)

	return min, nil
}

// Two returns the answer to the second part of the exercise.
func (e Exercise) Two(instr string) (any, error) {
	sections := strings.Split(instr, "\n\n")

	seedRng := parseSeedRange(sections[0])
	fmt.Println("seed ranges: ", len(seedRng))

	maps := parseAllMaps(sections[1:])

	min := math.MaxInt64

	for i, sr := range seedRng {
		seeds := make([]int, 0, sr.Range)
		for i := sr.Start; i < sr.Start+sr.Range; i++ {
			seeds = append(seeds, i)
		}

		locations := getLocations(maps, seeds)

		subMin := slices.Min(locations)

		fmt.Printf("range %d => min: %d, subMin: %d\n", i, min, subMin)

		if subMin < min {
			min = subMin
		}

	}

	return min, nil
}
