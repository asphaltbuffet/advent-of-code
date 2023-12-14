package exercises

import (
	"fmt"
	"strings"
)

type GalaxyMap map[Point]*GalaxyObject

type (
	MapRows map[int][]*GalaxyObject
	MapCols map[int][]*GalaxyObject

	RowExpansions []int
	ColExpansions []int
)

type GalaxyObject struct {
	ID  string
	Pos Point
}

type Point struct {
	X int
	Y int
}

type DataType rune

const (
	Empty  DataType = '.'
	Galaxy DataType = '#'
)

func expandImage(img []string) []string {
	expCols := getEmptyCols(img)
	expRows := getEmptyRows(img)

	expImg := make([]string, 0, len(img)+len(expRows))

	for i, r := range img {
		var rr strings.Builder
		if expRows[i] {
			expImg = append(expImg, "-")
			continue
		}

		for j, c := range r {
			if expCols[j] {
				rr.WriteRune('|')
			} else {
				rr.WriteRune(c)
			}
		}

		expImg = append(expImg, rr.String())

	}

	return expImg
}

func getEmptyCols(img []string) map[int]bool {
	cols := make(map[int]bool)

	for i := 0; i < len(img[0]); i++ {
		empty := true
		for _, r := range img {
			if r[i] != '.' {
				empty = false
				break
			}
		}

		if empty {
			cols[i] = empty
		}
	}

	return cols
}

func getEmptyRows(img []string) map[int]bool {
	rows := make(map[int]bool)

	for i, r := range img {
		empty := true
		for _, c := range r {
			if c != '.' {
				empty = false
				break
			}
		}

		if empty {
			rows[i] = empty
		}
	}

	return rows
}

func manhattanDistance(a, b Point) int {
	return abs(a.X-b.X) + abs(a.Y-b.Y)
}

func abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}

func sumDistances(img []string, factor int) int {
	galaxies := []Point{}
	var sum int
	var yOffset, xOffset int

	for y, line := range img {
		if line == "-" {
			yOffset += factor
			continue
		}

		xOffset = 0

		for x, c := range line {
			switch c {
			case '|':
				xOffset += factor

			case '#':
				p := Point{x + xOffset, y + yOffset}
				for _, g := range galaxies {
					sum += manhattanDistance(g, p)
				}

				galaxies = append(galaxies, p)
			case '-':
				fmt.Println("unexpected vertical expansion")
			default:
				// this should only be a '.'
			}
		}
	}

	return sum
}
