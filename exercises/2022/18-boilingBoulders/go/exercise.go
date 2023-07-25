package exercises

import (
	"fmt"
	"strings"

	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

type cube struct {
	x, y, z int
}

// Exercise for Advent of Code 2022 day 18.
type Exercise struct {
	common.BaseExercise
}

var minX, maxX, minY, maxY, minZ, maxZ int

var adjacent = []cube{
	{0, 0, -1},
	{0, 0, 1},
	{0, -1, 0},
	{0, 1, 0},
	{-1, 0, 0},
	{1, 0, 0},
}

// One returns the answer to the first part of the exercise.
// answer: 4456
func (c Exercise) One(instr string) (any, error) {
	cubes, err := parse(instr)
	if err != nil {
		return nil, fmt.Errorf("parsing input: %w", err)
	}

	cubeMap := make(map[cube]bool)
	for _, c := range cubes {
		cubeMap[c] = true
	}

	totalExposure := len(cubes) * 6

	// check all cubes for adjacent cubes
	for _, c := range cubes {
		for _, a := range adjacent {
			x := c.x + a.x
			y := c.y + a.y
			z := c.z + a.z

			if cubeMap[cube{x, y, z}] {
				totalExposure--
			}
		}
	}

	return totalExposure, nil
}

// Two returns the answer to the second part of the exercise.
// answer: 2510
func (c Exercise) Two(instr string) (any, error) {
	cubes, err := parse(instr)
	if err != nil {
		return nil, fmt.Errorf("parsing input: %w", err)
	}

	cubeMap := make(map[cube]bool)
	for _, c := range cubes {
		cubeMap[c] = true
	}

	totalExposure := 0

	// check all cube faces for path to edge
	for c := range cubeMap {
		totalExposure += facesThatCanReachEdge(c, cubeMap)
	}

	return totalExposure, nil
}

func parse(instr string) ([]cube, error) {
	var cubes []cube

	for _, line := range strings.Split(instr, "\n") {
		c := cube{}

		_, err := fmt.Sscanf(line, "%d,%d,%d", &c.x, &c.y, &c.z)
		if err != nil {
			return nil, fmt.Errorf("parsing cube from %q: %w", line, err)
		}

		if c.x < minX {
			minX = c.x
		}

		if c.x > maxX {
			maxX = c.x
		}

		if c.y < minY {
			minY = c.y
		}

		if c.y > maxY {
			maxY = c.y
		}

		if c.z < minZ {
			minZ = c.z
		}

		if c.z > maxZ {
			maxZ = c.z
		}

		cubes = append(cubes, c)
	}

	return cubes, nil
}

func facesThatCanReachEdge(cur cube, set map[cube]bool) int {
	count := 0

	for _, c := range adjacent {
		next := addCubes(cur, c)

		if canReachEdge(next, set) {
			count++
		}
	}

	return count
}

func canReachEdge(start cube, set map[cube]bool) bool {
	queue := []cube{start}
	visited := map[cube]bool{}

	for len(queue) > 0 {
		curCube := queue[0]
		queue = queue[1:]

		if visited[curCube] || set[curCube] {
			continue
		}

		visited[curCube] = true

		// edge reached
		if curCube.x <= minX || curCube.x >= maxX ||
			curCube.y <= minY || curCube.y >= maxY ||
			curCube.z <= minZ || curCube.z >= maxZ {
			return true
		}

		for _, a := range adjacent {
			next := addCubes(curCube, a)

			queue = append(queue, next)
		}
	}

	return false
}

func addCubes(c1, c2 cube) cube {
	return cube{
		c1.x + c2.x,
		c1.y + c2.y,
		c1.z + c2.z,
	}
}

// check all cubes using a flood fill algorithm
func floodFill(filled map[cube]bool) map[cube]bool {
	// set space bounds to encompass min/max values for x,y,z
	// use flood fill to find all cubes inside larger space that are not filled
	// return the number of orthagonal faces on the filled cubes that are adjascent to the larger space

	return nil
}
