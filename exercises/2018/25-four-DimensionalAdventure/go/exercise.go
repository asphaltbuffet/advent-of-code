package exercises

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

// Exercise for Advent of Code 2018 day 25.
type Exercise struct {
	common.BaseExercise
}

type point [4]int

func manhattan(a, b point) int {
	d := 0
	for i := 0; i < 4; i++ {
		x := a[i] - b[i]
		if x < 0 {
			x = -x
		}
		d += x
	}
	return d
}

// One returns the number of constellations: connected components of points where
// any two within Manhattan distance 3 join the same group.
func (e Exercise) One(instr string) (any, error) {
	pts := parse(instr)

	parent := make([]int, len(pts))
	for i := range parent {
		parent[i] = i
	}
	var find func(int) int
	find = func(x int) int {
		for parent[x] != x {
			parent[x] = parent[parent[x]] // path compression
			x = parent[x]
		}
		return x
	}
	union := func(a, b int) { parent[find(a)] = find(b) }

	for i := range pts {
		for j := i + 1; j < len(pts); j++ {
			if manhattan(pts[i], pts[j]) <= 3 {
				union(i, j)
			}
		}
	}

	roots := map[int]struct{}{}
	for i := range pts {
		roots[find(i)] = struct{}{}
	}
	return len(roots), nil
}

// Two is the day 25 finale: there is no second puzzle, only the closing message.
func (e Exercise) Two(instr string) (any, error) {
	return "Merry Christmas!", nil
}

// parse reads the 4D points, four comma-separated integers per line.
func parse(instr string) []point {
	re := regexp.MustCompile(`-?\d+`)
	var pts []point
	for _, line := range strings.Split(strings.TrimSpace(instr), "\n") {
		nums := re.FindAllString(line, -1)
		if len(nums) != 4 {
			continue
		}
		var p point
		for i, s := range nums {
			p[i], _ = strconv.Atoi(s)
		}
		pts = append(pts, p)
	}
	return pts
}
