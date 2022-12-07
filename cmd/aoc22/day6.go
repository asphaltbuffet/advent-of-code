package aoc22

import (
	"fmt"

	"github.com/asphaltbuffet/advent-of-code/internal/common"
	aoc "github.com/asphaltbuffet/advent-of-code/pkg/utilities"
)

// DataStreamMarker is a marker for message parts in the data stream.
type DataStreamMarker int

// DataStreamMarker values.
const (
	StartOfPacket  DataStreamMarker = 4
	StartOfMessage DataStreamMarker = 14
)

func init() { //nolint:gochecknoinits // init needed to register command
	common.NewDayCommand(2022, 6, D6P1, D6P2, Get2022Command())
}

// Day6 is the struct for day 6.
type Day6 struct{}

// D6P1 returns the solution for 2022 day 6 part 1.
// answer: 1707
func D6P1(data []string) string {
	m := data[0]

	loc, err := searchForMarker(m, StartOfPacket)
	if err != nil {
		return fmt.Sprintf("processing datastream: %v", err)
	}

	return fmt.Sprintf("%d", loc)
}

// D6P2 returns the solution for 2022 day 6 part 2.
// answer: 3697
func D6P2(data []string) string {
	m := data[0]

	loc, err := searchForMarker(m, StartOfMessage)
	if err != nil {
		return fmt.Sprintf("processing datastream: %v", err)
	}

	return fmt.Sprintf("%d", loc)
}

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
