package exercises

import (
	"fmt"
	"image"
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

	seen := map[image.Point]bool{}

	for _, s := range sensors {
		viewRange := s.Dist - abs(s.Location.Y-y)

		// fmt.Printf("Sensor at %v can see %d points\n", s.Location, viewRange*2+1)
		if viewRange > 0 {
			for i := 0; i <= viewRange; i++ {
				seen[image.Point{X: s.Location.X - i, Y: y}] = true
				seen[image.Point{X: s.Location.X + i, Y: y}] = true
			}
		}
	}

	// remove any beacons from the seen coordinates
	for _, s := range sensors {
		delete(seen, s.Beacon)
	}

	return len(seen), nil
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
