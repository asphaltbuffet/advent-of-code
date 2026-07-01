package exercises

import (
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
	// fmt.Println("seed count: ", len(seeds))

	maps := parseAllMaps(sections[1:])

	locations := getLocations(maps, seeds)

	min := slices.Min(locations)

	return min, nil
}

// Two returns the answer to the second part of the exercise. Seeds are handled
// as intervals mapped through each layer rather than enumerated individually,
// so billion-wide ranges resolve in microseconds.
func (e Exercise) Two(instr string) (any, error) {
	sections := strings.Split(instr, "\n\n")

	seedRng := parseSeedRange(sections[0])
	maps := parseAllMaps(sections[1:])

	seeds := make([]interval, 0, len(seedRng))
	for _, sr := range seedRng {
		seeds = append(seeds, interval{sr.Start, sr.Start + sr.Range})
	}

	locations := mapRangesToLocation(maps, seeds)

	best := math.MaxInt64
	for _, iv := range locations {
		if iv.start < best {
			best = iv.start
		}
	}

	return best, nil
}
