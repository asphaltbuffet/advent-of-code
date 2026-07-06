package exercises

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

// Exercise for Advent of Code 2018 day 6.
type Exercise struct {
	common.BaseExercise
}

type point struct{ x, y int }

var intRe = regexp.MustCompile(`-?\d+`)

// parse reads each coordinate by scanning its two integers, tolerant of spacing.
func parse(instr string) ([]point, error) {
	var pts []point

	for line := range strings.SplitSeq(strings.TrimSpace(instr), "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		nums := intRe.FindAllString(line, -1)
		if len(nums) != 2 {
			return nil, fmt.Errorf("expected 2 numbers, got %d in %q", len(nums), line)
		}

		x, _ := strconv.Atoi(nums[0])
		y, _ := strconv.Atoi(nums[1])
		pts = append(pts, point{x, y})
	}

	return pts, nil
}

func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}

func nearestPoint(pts []point, x, y int) (int, bool) {
	best, bestIdx, tie := 1<<30, -1, false
	for i, p := range pts {
		d := abs(p.x-x) + abs(p.y-y)
		switch {
		case d < best:
			best, bestIdx, tie = d, i, false
		case d == best:
			tie = true
		}
	}
	return bestIdx, tie
}

// bounds returns the min/max x and y over the coordinates.
func bounds(pts []point) (int, int, int, int) {
	minX, minY := pts[0].x, pts[0].y
	maxX, maxY := pts[0].x, pts[0].y
	for _, p := range pts {
		if p.x < minX {
			minX = p.x
		}
		if p.x > maxX {
			maxX = p.x
		}
		if p.y < minY {
			minY = p.y
		}
		if p.y > maxY {
			maxY = p.y
		}
	}
	return minX, minY, maxX, maxY
}

// One returns the answer to the first part of the exercise.
// answer: 2906
func (e Exercise) One(instr string) (any, error) {
	pts, err := parse(instr)
	if err != nil {
		return nil, err
	}

	minX, minY, maxX, maxY := bounds(pts)

	area := make([]int, len(pts))
	infinite := make([]bool, len(pts))

	for y := minY; y <= maxY; y++ {
		for x := minX; x <= maxX; x++ {
			bestIdx, tie := nearestPoint(pts, x, y)
			if tie || bestIdx < 0 {
				continue
			}
			area[bestIdx]++
			if x == minX || x == maxX || y == minY || y == maxY {
				infinite[bestIdx] = true
			}
		}
	}

	largest := 0
	for i, a := range area {
		if !infinite[i] && a > largest {
			largest = a
		}
	}

	return largest, nil
}

// Two returns the answer to the second part of the exercise.
//
// The region is every location whose summed Manhattan distance to all
// coordinates is below a threshold. The real puzzle uses 10000; the small
// worked example uses 32. The example is distinguishable by its handful of
// coordinates, so the threshold is chosen from the coordinate count.
//
// A location beyond the bounding box adds at least numPts per step of extra
// margin to the total distance, so the region can extend at most
// threshold/numPts cells outside the box — scanning that padded box captures
// the whole region.
// answer: 50530
func (e Exercise) Two(instr string) (any, error) {
	pts, err := parse(instr)
	if err != nil {
		return nil, err
	}

	threshold := 10000
	if len(pts) <= 10 {
		threshold = 32
	}

	minX, minY, maxX, maxY := bounds(pts)
	pad := threshold/len(pts) + 1

	size := 0
	for y := minY - pad; y <= maxY+pad; y++ {
		for x := minX - pad; x <= maxX+pad; x++ {
			total := 0
			for _, p := range pts {
				total += abs(p.x-x) + abs(p.y-y)
			}
			if total < threshold {
				size++
			}
		}
	}

	return size, nil
}
