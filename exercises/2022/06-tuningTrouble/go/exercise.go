package exercises

import (
	"fmt"

	"github.com/asphaltbuffet/advent-of-code/internal/common"
	aoc "github.com/asphaltbuffet/advent-of-code/pkg/utilities"
)

// Exercise for Advent of Code 2022 day 6
type Exercise struct {
	common.BaseExercise
}

// One returns the answer to the first part of the exercise. // answer: 1707
func (c Exercise) One(instr string) (any, error) {
	loc, err := searchForMarker(instr, StartOfPacket)
	if err != nil {
		return nil, fmt.Errorf("processing datastream: %w", err)
	}

	return fmt.Sprintf("%d", loc), nil
}

// Two returns the answer to the second part of the exercise. // answer: 3697
func (c Exercise) Two(instr string) (any, error) {
	loc, err := searchForMarker(instr, StartOfMessage)
	if err != nil {
		return nil, fmt.Errorf("processing datastream: %w", err)
	}

	return fmt.Sprintf("%d", loc), nil
}

// DataStreamMarker is a marker for message parts in the data stream.
type DataStreamMarker int

// DataStreamMarker values.
const (
	StartOfPacket  DataStreamMarker = 4
	StartOfMessage DataStreamMarker = 14
)

func searchForMarker(ds string, marker DataStreamMarker) (int, error) {
	m := int(marker)

	if len(ds) < m {
		return 0, fmt.Errorf("not found")
	}

	// fmt.Printf("checking %s for %c\n", ds[1:m], ds[0])

	d := len(aoc.Unique([]byte(ds[0:m])))
	if d == m {
		return m, nil
	}

	// would be a little more efficient if this started after any other duplicates in substring (like 'dppd' or 'aaab')
	n, err := searchForMarker(ds[1:], marker)

	return n + 1, err
}
