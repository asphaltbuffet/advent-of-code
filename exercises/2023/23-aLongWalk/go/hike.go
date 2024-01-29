package exercises

import "strings"

type Point struct {
	X int
	Y int
}

type Tile rune

type TrailMap struct {
	Tiles     map[Point]Tile
	Start     Point
	End       Point
	Junctions map[Point]bool
}

const (
	Path       Tile = '.'
	Forest     Tile = '#'
	SlopeUp    Tile = '^'
	SlopeDown  Tile = 'v'
	SlopeLeft  Tile = '<'
	SlopeRight Tile = '>'
)

func parseInput(s string) *TrailMap {
	lines := strings.Split(s, "\n")

	tm := &TrailMap{
		Tiles: make(map[Point]Tile, len(lines)*len(lines[0])),
		Start: Point{1, 0},
		End:   Point{len(lines[0]) - 2, len(lines) - 1},
	}

	for row, line := range lines {
		for col, c := range line {
			// ignore forest tiles completely
			if t := Tile(c); t != Forest {
				tm.Tiles[Point{col, row}] = Tile(c)
			}
		}
	}

	return tm
}

func (tm *TrailMap) getJunctions() {
	tm.Junctions = map[Point]bool{}

	for p := range tm.Tiles {
		neighbours := 0

		for _, dir := range [4]Point{{1, 0}, {-1, 0}, {0, 1}, {0, -1}} {
			next := Point{p.X + dir.X, p.Y + dir.Y}

			if _, ok := tm.Tiles[next]; ok {
				neighbours++

				if neighbours > 2 {
					tm.Junctions[p] = true
					break
				}
			}
		}
	}

	tm.Junctions[tm.Start] = true
	tm.Junctions[tm.End] = true
}

type PathTo struct {
	end           Point
	length, index int
}

func (tm TrailMap) getPaths() map[Point][]PathTo {
	paths := map[Point][]PathTo{}
	junctionIndex := 0

	for junctionPoint := range tm.Junctions {
		blocked := -1

		for i, startDir := range [4]Point{{0, -1}, {-1, 0}, {0, 1}, {1, 0}} {
			currentPoint := Point{junctionPoint.X + startDir.X, junctionPoint.Y + startDir.Y}

			if _, ok := tm.Tiles[currentPoint]; ok {
				path := tm.getPath(currentPoint, startDir, 1)

				path.index = junctionIndex
				paths[junctionPoint] = append(paths[junctionPoint], path)
			} else {
				blocked = i
			}
		}

		if blocked != -1 && junctionPoint != tm.Start && junctionPoint != tm.End {
			removeIndex := 0

			if blocked == 2 {
				// note: if up is blocked, then left is at index 0, otherwise at 1
				removeIndex = 1
			}

			paths[junctionPoint] = append(paths[junctionPoint][:removeIndex], paths[junctionPoint][removeIndex+1:]...)
		}

		junctionIndex++
	}

	return paths
}

func ToLeft(p Point) Point {
	return Point{p.Y, -p.X}
}

func ToRight(p Point) Point {
	return Point{-p.Y, p.X}
}

func (tm TrailMap) getPath(curPt, curDir Point, plen int) PathTo {
	for _, dir := range [3]Point{curDir, ToLeft(curDir), ToRight(curDir)} {
		next := Point{curPt.X + dir.X, curPt.Y + dir.Y}

		if _, ok := tm.Tiles[next]; !ok {
			continue
		}

		if _, found := tm.Junctions[next]; found {
			return PathTo{next, plen + 1, 0}
		}

		return tm.getPath(next, dir, plen+1)
	}

	// should never happen
	return PathTo{Point{-1, -1}, 0, 0}
}

func (tm TrailMap) getLongestPath(paths map[Point][]PathTo, start, end Point, step int, visited []bool) int {
	maxStep := 0

	for _, path := range paths[start] {
		if path.end == end {
			return step + path.length
		}

		index := paths[path.end][0].index

		if !visited[index] {
			visited[index] = true
			maxStep = max(maxStep, tm.getLongestPath(paths, path.end, end, step+path.length, visited))
			visited[index] = false
		}
	}

	return maxStep
}

func (tm TrailMap) walk(start, currentDir Point, visited map[Point]int) {
	current := start
	currentStep := visited[current]

	for _, dir := range [3]Point{currentDir, ToLeft(currentDir), ToRight(currentDir)} {
		next := Point{current.X + dir.X, current.Y + dir.Y}
		if nt, ok := tm.Tiles[next]; ok {
			slope := nt

			oppositeSlope := map[Point]Tile{{1, 0}: SlopeLeft, {-1, 0}: SlopeRight, {0, 1}: SlopeUp, {0, -1}: SlopeDown}
			if oppositeSlope[dir] == slope {
				continue
			}

			if val, found := visited[next]; !found || val < currentStep+1 {
				visited[next] = currentStep + 1

				tm.walk(next, dir, visited)
			}
		}
	}
}
