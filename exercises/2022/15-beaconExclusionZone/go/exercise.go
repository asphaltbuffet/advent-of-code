package exercises

import (
	"fmt"
	"image"
	"sort"
	"strings"

	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

var (
	TargetY = 2000000
	TestY   = 10
)

const (
	minBoundary int = 0
	maxBoundary int = 4000000
)

// Sensor is the X, Y coordinates of a sensor and the closest beacon with its distance.
type Sensor struct {
	Location image.Point
	Beacon   image.Point
	Dist     int
}

// Exercise for Advent of Code 2022 day 15
type Exercise struct {
	common.BaseExercise
}

// One returns the answer to the first part of the exercise.
// wrong: 4807593 (too low)
// answer:
func (c Exercise) One(instr string) (any, error) {
	sensors, err := parse(instr)
	if err != nil {
		return nil, fmt.Errorf("parsing input: %w", err)
	}

	y := TestY
	if len(sensors) == 23 {
		y = TargetY
	}

	// Each sensor covers a contiguous [left,right] span on row y. Collect and
	// merge those spans instead of enumerating individual cells.
	type span struct{ lo, hi int }
	spans := make([]span, 0, len(sensors))
	for _, s := range sensors {
		reach := s.Dist - abs(s.Location.Y-y)
		if reach < 0 {
			continue
		}
		spans = append(spans, span{s.Location.X - reach, s.Location.X + reach})
	}

	sort.Slice(spans, func(i, j int) bool { return spans[i].lo < spans[j].lo })

	// Sum the width of the merged spans (count of covered x positions).
	covered := 0
	curLo, curHi := spans[0].lo, spans[0].hi
	merged := make([]span, 0, len(spans))
	for _, sp := range spans[1:] {
		if sp.lo > curHi+1 {
			covered += curHi - curLo + 1
			merged = append(merged, span{curLo, curHi})
			curLo, curHi = sp.lo, sp.hi
			continue
		}
		if sp.hi > curHi {
			curHi = sp.hi
		}
	}
	covered += curHi - curLo + 1
	merged = append(merged, span{curLo, curHi})

	// A beacon sitting on row y occupies a covered cell but is not "no beacon".
	beacons := map[int]bool{}
	for _, s := range sensors {
		if s.Beacon.Y != y {
			continue
		}
		for _, m := range merged {
			if s.Beacon.X >= m.lo && s.Beacon.X <= m.hi {
				beacons[s.Beacon.X] = true
				break
			}
		}
	}

	return covered - len(beacons), nil
}

// Two returns the answer to the second part of the exercise.
// answer: 11747175442119
func (c Exercise) Two(instr string) (any, error) {
	sensors, err := parse(instr)
	if err != nil {
		return nil, fmt.Errorf("parsing input: %w", err)
	}

	limit := 20

	if len(sensors) == 23 {
		limit = maxBoundary
	}

	for _, s := range sensors {
		// not to far, not too close
		goldilocks := s.Dist + 1

		for row := -goldilocks; row <= goldilocks; row++ {
			curRow := s.Location.Y + row

			if curRow < minBoundary {
				continue
			}

			if curRow > limit {
				break
			}

			offsetX := goldilocks - abs(row)
			leftX := s.Location.X - offsetX
			rightX := s.Location.X + offsetX

			if leftX >= minBoundary && leftX <= limit && !IsReachableLocation(image.Point{X: leftX, Y: curRow}, sensors) {
				return leftX*4000000 + curRow, nil
			}

			if rightX >= minBoundary && rightX <= limit && !IsReachableLocation(image.Point{X: rightX, Y: curRow}, sensors) {
				return rightX*4000000 + curRow, nil
			}
		}
	}

	return nil, fmt.Errorf("no solution found")
}

// Parse sensor and beacon data from input.
func parse(data string) ([]Sensor, error) {
	var sensors []Sensor

	for n, line := range strings.Split(data, "\n") {
		var s Sensor

		_, err := fmt.Sscanf(line, "Sensor at x=%d, y=%d: closest beacon is at x=%d, y=%d",
			&s.Location.X, &s.Location.Y, &s.Beacon.X, &s.Beacon.Y)
		if err != nil {
			return nil, fmt.Errorf("parsing line %d: %q", n, line)
		}

		s.Dist = ManhattanDistance(s.Location, s.Beacon)

		sensors = append(sensors, s)
	}

	return sensors, nil
}

func ManhattanDistance(p1, p2 image.Point) int {
	return abs(p1.X-p2.X) + abs(p1.Y-p2.Y)
}

func IsReachableLocation(p image.Point, sensors []Sensor) bool {
	for _, s := range sensors {
		if s.Dist >= ManhattanDistance(p, s.Location) {
			return true
		}
	}

	return false
}

func abs(n int) int {
	if n < 0 {
		return -n
	}

	return n
}
