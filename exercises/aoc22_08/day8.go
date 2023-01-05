// Package aoc22_08 contains the solution for day 8 of Advent of Code 2022.
package aoc22_08 //nolint:revive,stylecheck // I don't care about the package name

import (
	"sort"
	"strconv"

	aoc "github.com/asphaltbuffet/advent-of-code/pkg/utilities"
)

var dimX, dimY int

// D8P1 returns the solution for 2022 day 8 part 1.
// answer: 1805
func D8P1(data []string) string {
	dimX = len(data[0])
	dimY = len(data)

	m := GetTreeMap(data)

	n := GetVisibleTreesFromNorth(m)
	s := GetVisibleTreesFromSouth(m)
	e := GetVisibleTreesFromEast(m)
	w := GetVisibleTreesFromWest(m)

	// fmt.Printf("trees: %+v\n", trees)
	// file, _ := os.Create("./mygraph.gv")
	// _ = draw.DOT(tGraph, file)
	out := []int{0, dimX - 1, dimY * (dimX - 1), (dimX * dimY) - 1} // corners
	out = append(out, n...)
	out = append(out, s...)
	out = append(out, e...)
	out = append(out, w...)
	out = aoc.Unique(out)

	sort.Ints(out)
	// fmt.Printf("visible trees: %v\n", out)

	return strconv.Itoa(len(out) - 1) // -1 for the "edge" tree
}

// D8P2 returns the solution for 2022 day 8 part 2.
// answer: 444528
func D8P2(data []string) string {
	dimX = len(data[0])
	dimY = len(data)

	m := GetTreeMap(data)
	maxScenicScore := 0

	for r := 1; r < dimY-1; r++ {
		for c := 1; c < dimX-1; c++ {
			h := m[r][c]
			up := CalculateScoreUp(h, r, c, m)
			down := CalculateScoreDown(h, r, c, m)
			left := CalculateScoreLeft(h, r, c, m)
			right := CalculateScoreRight(h, r, c, m)

			score := up * down * left * right

			// fmt.Printf("[%d,%d]=%d, score: %d * %d * %d * %d = %d\n", r, c, h, up, down, left, right, score)

			if score > maxScenicScore {
				maxScenicScore = score
			}
		}
	}

	return strconv.Itoa(maxScenicScore)
}
