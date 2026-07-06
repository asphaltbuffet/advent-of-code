package exercises

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

// Exercise for Advent of Code 2018 day 10.
type Exercise struct {
	common.BaseExercise
}

var intRe = regexp.MustCompile(`-?\d+`)

// star is a point of light with a position and constant velocity.
type star struct {
	x, y, vx, vy int
}

// parse reads each star's four integers.
func parse(instr string) ([]star, error) {
	var stars []star

	for line := range strings.SplitSeq(strings.TrimSpace(instr), "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		nums := intRe.FindAllString(line, -1)
		if len(nums) != 4 {
			return nil, fmt.Errorf("expected 4 numbers, got %d in %q", len(nums), line)
		}

		var s star
		s.x, _ = strconv.Atoi(nums[0])
		s.y, _ = strconv.Atoi(nums[1])
		s.vx, _ = strconv.Atoi(nums[2])
		s.vy, _ = strconv.Atoi(nums[3])
		stars = append(stars, s)
	}

	return stars, nil
}

// extent returns the combined width and height of the stars' bounding box at time
// t. The message appears when the stars are tightest, and this extent reaches its
// minimum there.
func extent(stars []star, t int) int {
	minX, minY := 1<<62, 1<<62
	maxX, maxY := -(1 << 62), -(1 << 62)
	for _, s := range stars {
		x, y := s.x+s.vx*t, s.y+s.vy*t
		if x < minX {
			minX = x
		}
		if x > maxX {
			maxX = x
		}
		if y < minY {
			minY = y
		}
		if y > maxY {
			maxY = y
		}
	}
	return (maxX - minX) + (maxY - minY)
}

// converge returns the second at which the stars form the message: the time that
// minimizes their bounding-box extent. The extent shrinks to that minimum then
// grows, so we step forward until it stops shrinking.
func converge(stars []star) int {
	t := 0
	for extent(stars, t+1) < extent(stars, t) {
		t++
	}
	return t
}

// render draws the stars at time t as a grid of '#' and ' '.
func render(stars []star, t int) string {
	minX, minY := 1<<62, 1<<62
	maxX, maxY := -(1 << 62), -(1 << 62)
	for _, s := range stars {
		x, y := s.x+s.vx*t, s.y+s.vy*t
		if x < minX {
			minX = x
		}
		if x > maxX {
			maxX = x
		}
		if y < minY {
			minY = y
		}
		if y > maxY {
			maxY = y
		}
	}

	// Solid blocks render far more legibly than '#' for reading the letters.
	const lit, dark = '█', ' '

	w, h := maxX-minX+1, maxY-minY+1
	grid := make([][]rune, h)
	for i := range grid {
		grid[i] = make([]rune, w)
		for j := range grid[i] {
			grid[i][j] = dark
		}
	}
	for _, s := range stars {
		grid[s.y+s.vy*t-minY][s.x+s.vx*t-minX] = lit
	}

	var b strings.Builder
	for i, row := range grid {
		if i > 0 {
			b.WriteByte('\n')
		}
		b.WriteString(string(row))
	}
	return b.String()
}

// One returns the answer to the first part of the exercise.
// The rendered message reads: BFFZCNXE
func (e Exercise) One(instr string) (any, error) {
	stars, err := parse(instr)
	if err != nil {
		return nil, err
	}

	return render(stars, converge(stars)), nil
}

// Two returns the answer to the second part of the exercise.
// Answer: 10391
func (e Exercise) Two(instr string) (any, error) {
	stars, err := parse(instr)
	if err != nil {
		return nil, err
	}

	// The number of seconds until the message appears is the convergence time.
	return converge(stars), nil
}
