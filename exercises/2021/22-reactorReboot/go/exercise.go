package exercises

import (
	"fmt"
	"strings"

	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

// Exercise for Advent of Code 2021 day 22.
type Exercise struct {
	common.BaseExercise
}

// cuboid is a closed integer box [x1,x2]×[y1,y2]×[z1,z2] with a sign used for
// inclusion-exclusion accounting.
type cuboid struct {
	x1, x2, y1, y2, z1, z2 int
	sign                   int
}

type step struct {
	on bool
	b  cuboid
}

func parse(instr string) ([]step, error) {
	var steps []step
	for _, line := range strings.Split(strings.TrimSpace(instr), "\n") {
		if line == "" {
			continue
		}
		var state string
		var b cuboid
		_, err := fmt.Sscanf(line, "%s x=%d..%d,y=%d..%d,z=%d..%d",
			&state, &b.x1, &b.x2, &b.y1, &b.y2, &b.z1, &b.z2)
		if err != nil {
			return nil, fmt.Errorf("parsing step %q: %w", line, err)
		}
		steps = append(steps, step{on: state == "on", b: b})
	}
	return steps, nil
}

// volume returns the signed volume of a closed integer cuboid.
func (c cuboid) volume() int {
	return c.sign * (c.x2 - c.x1 + 1) * (c.y2 - c.y1 + 1) * (c.z2 - c.z1 + 1)
}

// intersect returns the overlap of two boxes and whether one exists. The result
// carries the opposite sign of a, used to cancel double-counted regions.
func intersect(a, b cuboid) (cuboid, bool) {
	x1, x2 := max(a.x1, b.x1), min(a.x2, b.x2)
	y1, y2 := max(a.y1, b.y1), min(a.y2, b.y2)
	z1, z2 := max(a.z1, b.z1), min(a.z2, b.z2)
	if x1 > x2 || y1 > y2 || z1 > z2 {
		return cuboid{}, false
	}
	return cuboid{x1, x2, y1, y2, z1, z2, -a.sign}, true
}

// reboot applies the steps using signed inclusion-exclusion and returns the
// number of cubes left on.
func reboot(steps []step) int {
	var boxes []cuboid
	for _, s := range steps {
		var added []cuboid
		for _, existing := range boxes {
			if ov, ok := intersect(existing, s.b); ok {
				added = append(added, ov)
			}
		}
		if s.on {
			nb := s.b
			nb.sign = 1
			added = append(added, nb)
		}
		boxes = append(boxes, added...)
	}

	total := 0
	for _, b := range boxes {
		total += b.volume()
	}
	return total
}

// One counts cubes on after the initialization steps, clipped to -50..50.
func (e Exercise) One(instr string) (any, error) {
	steps, err := parse(instr)
	if err != nil {
		return nil, err
	}

	var clipped []step
	for _, s := range steps {
		b := s.b
		if b.x1 < -50 || b.x2 > 50 || b.y1 < -50 || b.y2 > 50 || b.z1 < -50 || b.z2 > 50 {
			continue
		}
		clipped = append(clipped, s)
	}

	return fmt.Sprintf("%d", reboot(clipped)), nil
}

// Two counts cubes on after the full reboot sequence.
func (e Exercise) Two(instr string) (any, error) {
	steps, err := parse(instr)
	if err != nil {
		return nil, err
	}

	return fmt.Sprintf("%d", reboot(steps)), nil
}
