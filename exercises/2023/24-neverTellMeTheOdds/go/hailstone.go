package exercises

import (
	"fmt"
	"math"
	"strings"
)

type Hailstone struct {
	pos Vector3D
	vel Vector3D
}

func (hs Hailstone) Intersects(hs2 Hailstone) bool {
	const (
		// // Test values
		// min float64 = 7.0
		// max float64 = 27.0

		min float64 = 200000000000000.0
		max float64 = 400000000000000.0
	)

	point, ok := getIntersectionPointWithinBounds(hs, hs2, min, max)
	if !ok {
		return false
	}

	dx := point.X - hs.pos.X
	dy := point.Y - hs.pos.Y

	if (dx > 0) != (hs.vel.X > 0) || (dy > 0) != (hs.vel.Y > 0) {
		return false
	}

	dx = point.X - hs2.pos.X
	dy = point.Y - hs2.pos.Y

	return (dx > 0) == (hs2.vel.X > 0) && (dy > 0) == (hs2.vel.Y > 0)
}

// getIntersectionPointWithinBounds calculates the intersection of two Hailstone objects
// based on their current velocities and positions.
//
// It returns the point of intersection as a Vector2D and a boolean indicating whether an
// intersection will occur within the given bounds.
func getIntersectionPointWithinBounds(a, b Hailstone, min, max float64) (Vector2D, bool) {
	velA := Vector2D{a.vel.X, a.vel.Y}
	velB := Vector2D{b.vel.X, b.vel.Y}

	// get determinant and check if vectors are parallel
	determinant := velA.Cross(velB)
	if isEqualToZero(determinant) {
		return Vector2D{-1, -1}, false
	}

	posDelta := Vector2D{b.pos.X - a.pos.X, b.pos.Y - a.pos.Y}

	scalar := posDelta.Cross(velB) / determinant
	point := Vector2D{a.pos.X + a.vel.X*scalar, a.pos.Y + a.vel.Y*scalar}

	// check if point is within given bounds
	if point.X < min || point.X > max || point.Y < min || point.Y > max {
		return Vector2D{-1, -1}, false
	}

	return point, true
}

// isEqualToZero checks if a float64 is effectively equal to zero within a tolerance.
func isEqualToZero(f float64) bool {
	const epsilon float64 = 1e-9 // Adjust epsilon based on the precision you need.

	return math.Abs(f) < epsilon
}

func parseInput(input string) ([]Hailstone, error) {
	lines := strings.Split(input, "\n")
	hailStones := make([]Hailstone, 0, len(lines))

	for _, line := range lines {
		var px, py, pz, vx, vy, vz float64

		n, err := fmt.Sscanf(line, "%f, %f, %f @ %f, %f, %f", &px, &py, &pz, &vx, &vy, &vz)
		if err != nil || n != 6 {
			return nil, fmt.Errorf("parsing %q: %w", line, err)
		}

		hailStone := Hailstone{Vector3D{px, py, pz}, Vector3D{vx, vy, vz}}
		hailStones = append(hailStones, hailStone)
	}
	return hailStones, nil
}

func calcVelMatch(curPos, targetPos int) []int {
	matchingVels := []int{}

	for vel := -1000; vel < 1000; vel++ { // TODO: this is a hack
		// Check if the velocity is different from the target position and if it would result in
		// reaching the target position from the current position in one time unit.
		if vel != targetPos && curPos%(vel-targetPos) == 0 {
			matchingVels = append(matchingVels, vel)
		}
	}

	return matchingVels
}

// getIntersect finds the intersection of two slices.
func getIntersect(a, b []int) []int {
	intersection := make([]int, 0)

	for _, v := range a {
		for _, w := range b {
			if v == w {
				intersection = append(intersection, v)
				break
			}
		}
	}

	return intersection
}

type Dimension string

const (
	X Dimension = "X"
	Y Dimension = "Y"
	Z Dimension = "Z"
)

func updatePotential(d Dimension, a, b Hailstone, potentialVels map[Dimension][]int) {
	var deltaPos int
	var velA, velB int

	switch d {
	case X:
		deltaPos, velA, velB = int(b.pos.X-a.pos.X), int(a.vel.X), int(b.vel.X)
	case Y:
		deltaPos, velA, velB = int(b.pos.Y-a.pos.Y), int(a.vel.Y), int(b.vel.Y)
	case Z:
		deltaPos, velA, velB = int(b.pos.Z-a.pos.Z), int(a.vel.Z), int(b.vel.Z)
	}

	if velA != velB {
		return
	}

	nextVel := calcVelMatch(deltaPos, velA)

	if len(potentialVels[d]) == 0 {
		potentialVels[d] = nextVel
	} else {
		potentialVels[d] = getIntersect(potentialVels[d], nextVel)
	}
}

// getPosition calculates the 3D position of a rock with a given velocity intersects
// the path between two Hailstone objects in 3D space.
func getPosition(rockVel Vector3D, a, b Hailstone) (Vector3D, error) {
	// calculate slopes of paths of rock to each Hailstone on the X/Y plane
	slopeA := (a.vel.Y - rockVel.Y) / (a.vel.X - rockVel.X)
	slopeB := (b.vel.Y - rockVel.Y) / (b.vel.X - rockVel.X)

	// calculate y-intercepts of paths of rock to each Hailstone on the X/Y plane
	interceptA := a.pos.Y - (slopeA * a.pos.X)
	interceptB := b.pos.Y - (slopeB * b.pos.X)

	// calculate the intersection point of the two paths on the X/Y plane
	xPos := (interceptB - interceptA) / (slopeA - slopeB)
	yPos := slopeA*xPos + interceptA

	// calculate the t it takes for the rock to reach the intersection point
	t := (xPos - a.pos.X) / (a.vel.X - rockVel.X)

	// calculate the Z position of the rock at the intersection point
	zPos := a.pos.Z + (a.vel.Z-rockVel.Z)*t

	return Vector3D{xPos, yPos, zPos}, nil
}
