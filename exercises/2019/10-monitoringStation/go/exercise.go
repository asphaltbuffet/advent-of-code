package exercises

import (
	"math"
	"sort"
	"strings"

	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

// Exercise for Advent of Code 2019 day 10.
type Exercise struct {
	common.BaseExercise
}

func parseAsteroids(instr string) [][2]int {
	var asteroids [][2]int
	for y, line := range strings.Split(strings.TrimSpace(instr), "\n") {
		for x, ch := range line {
			if ch == '#' {
				asteroids = append(asteroids, [2]int{x, y})
			}
		}
	}
	return asteroids
}

func gcd(a, b int) int {
	if a < 0 {
		a = -a
	}
	if b < 0 {
		b = -b
	}
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func normalize(dx, dy int) [2]int {
	if dx == 0 && dy == 0 {
		return [2]int{0, 0}
	}
	g := gcd(dx, dy)
	if g < 0 {
		g = -g
	}
	return [2]int{dx / g, dy / g}
}

func bestVisible(asteroids [][2]int) ([2]int, int) {
	best := 0
	var bestPos [2]int
	dirs := make(map[[2]int]bool, len(asteroids))
	for _, a := range asteroids {
		for k := range dirs {
			delete(dirs, k)
		}
		for _, b := range asteroids {
			if a == b {
				continue
			}
			dx := b[0] - a[0]
			dy := b[1] - a[1]
			dirs[normalize(dx, dy)] = true
		}
		if len(dirs) > best {
			best = len(dirs)
			bestPos = a
		}
	}
	return bestPos, best
}

// One returns the answer to the first part of the exercise.
func (e Exercise) One(instr string) (any, error) {
	asteroids := parseAsteroids(instr)
	_, count := bestVisible(asteroids)
	return count, nil
}

// Two returns the answer to the second part of the exercise.
//
//nolint:gocognit // asteroid vaporization spiral tracking requires many state branches
func (e Exercise) Two(instr string) (any, error) {
	asteroids := parseAsteroids(instr)
	station, _ := bestVisible(asteroids)

	// Group remaining asteroids by normalized direction from the station.
	type asteroid struct{ x, y int }
	dirMap := make(map[[2]int][]asteroid)
	for _, a := range asteroids {
		if a == station {
			continue
		}
		dx := a[0] - station[0]
		dy := a[1] - station[1]
		norm := normalize(dx, dy)
		dirMap[norm] = append(dirMap[norm], asteroid{a[0], a[1]})
	}

	// Sort each direction's asteroids by distance (nearest first).
	abs := func(v int) int {
		if v < 0 {
			return -v
		}
		return v
	}
	for k := range dirMap {
		sort.Slice(dirMap[k], func(i, j int) bool {
			di := abs(dirMap[k][i].x-station[0]) + abs(dirMap[k][i].y-station[1])
			dj := abs(dirMap[k][j].x-station[0]) + abs(dirMap[k][j].y-station[1])
			return di < dj
		})
	}

	// Collect directions sorted by clockwise angle starting from "up" (negative Y).
	// angle = atan2(dx, -dy) gives 0 for straight up, increasing clockwise.
	dirs := make([][2]int, 0, len(dirMap))
	for k := range dirMap {
		dirs = append(dirs, k)
	}
	clockwiseAngle := func(d [2]int) float64 {
		a := math.Atan2(float64(d[0]), -float64(d[1]))
		if a < 0 {
			a += 2 * math.Pi
		}
		return a
	}
	sort.Slice(dirs, func(i, j int) bool {
		return clockwiseAngle(dirs[i]) < clockwiseAngle(dirs[j])
	})

	// Simulate laser sweeps; vaporize in angle order, cycling until 200th.
	count := 0
	for {
		for _, d := range dirs {
			if len(dirMap[d]) == 0 {
				continue
			}
			a := dirMap[d][0]
			dirMap[d] = dirMap[d][1:]
			count++
			if count == 200 {
				return a.x*100 + a.y, nil
			}
		}
		// Check if any asteroids remain.
		remaining := false
		for _, d := range dirs {
			if len(dirMap[d]) > 0 {
				remaining = true
				break
			}
		}
		if !remaining {
			break
		}
	}

	return -1, nil
}
