package exercises

import (
	"fmt"
	"strings"
)

type basin struct {
	start     point
	end       point
	winds     []wind
	totalRows int
	totalCols int
}

var directions = []point{
	{1, 0},
	{0, 1},
	{0, -1},
	{-1, 0},
}

func parseInput(input string) (*basin, error) {
	b := &basin{
		start:     point{},
		end:       point{},
		winds:     []wind{},
		totalRows: 0,
		totalCols: 0,
	}

	lines := strings.Split(input, "\n")

	for c := 0; c < len(lines); c++ {
		if lines[0][c] == '.' {
			b.start = point{c - 1, -1}
			break
		}
	}

	// 0,0 will be top left of box defined by # boundaries
	b.totalRows = len(lines) - 2
	b.totalCols = len(lines[0]) - 2

	for c := 0; c < len(lines[0]); c++ {
		if lines[len(lines)-1][c] == '.' {
			b.end = point{c - 1, b.totalRows}
			break
		}
	}

	for row := 1; row < len(lines)-1; row++ {
		chars := strings.Split(lines[row], "")
		for col := 1; col < len(chars)-1; col++ {
			p := point{col - 1, row - 1}

			switch c := chars[col]; c {
			case ">", "<", "^", "v":
				b.winds = append(b.winds, wind{
					start:     p,
					direction: relative[c],
					totalRows: b.totalRows,
					totalCols: b.totalCols,
					char:      c,
				})
			case ".", "#":
				// do nothing
				continue
			default:
				return nil, fmt.Errorf("unexpected character %q at [%d,%d]", c, col, row)
			}
		}
	}

	return b, nil
}

func calcPath(blizzards []wind, start, end point, totalRows, totalCols, stepsElapsedAlready int) (int, error) {
	cacheRoomStates := make(map[int]map[point]string, totalCols*totalRows)

	type moment struct {
		coords point
		steps  int
	}

	queue := []moment{}
	queue = append(queue, moment{
		coords: start,
		steps:  stepsElapsedAlready,
	})

	visited := map[[3]int]bool{}

	for len(queue) > 0 {
		curPoint := queue[0]
		queue = queue[1:]

		roomState := getRoomState(blizzards, curPoint.steps+1, totalRows, totalCols, cacheRoomStates)

		for _, diff := range directions {
			nextCoords := curPoint.coords.add(diff)

			if nextCoords == start {
				continue
			} else if nextCoords != end {
				if nextCoords.y < 0 || nextCoords.y >= totalRows ||
					nextCoords.x < 0 || nextCoords.x >= totalCols {
					continue
				}
			}

			// no point in processing a coordinate & steps pair that has already been seen
			hash := [3]int{nextCoords.x, nextCoords.y, curPoint.steps + 1}
			if visited[hash] {
				continue
			}

			visited[hash] = true

			if nextCoords != start && nextCoords != end {
				// if blocked, continue
				if roomState[nextCoords] != "." {
					continue
				} else if nextCoords.y < 0 || nextCoords.y >= totalRows ||
					nextCoords.x < 0 || nextCoords.x >= totalCols {
					// if out of bounds, continue
					continue
				}
			}

			// done
			if nextCoords == end {
				return curPoint.steps + 1, nil
			}

			queue = append(queue, moment{
				coords: nextCoords,
				steps:  curPoint.steps + 1,
			})
		}
		// if possible to stay still, add "wait" move
		if curPoint.coords == start ||
			roomState[curPoint.coords] == "." {
			queue = append(queue, moment{
				coords: curPoint.coords,
				steps:  curPoint.steps + 1,
			})
		}
	}

	return 0, fmt.Errorf("no path found")
}

func getRoomState(winds []wind, steps, totalRows, totalCols int, stateAt map[int]map[point]string) map[point]string {
	if s, ok := stateAt[steps]; ok {
		return s
	}

	state := make(map[point]string, totalRows*totalCols)

	for _, w := range winds {
		coords := w.extrapolatePosition(steps)
		state[coords] = w.char
	}

	for r := 0; r < totalRows; r++ {
		for c := 0; c < totalCols; c++ {
			if state[point{c, r}] == "" {
				state[point{c, r}] = "."
			}
		}
	}

	stateAt[steps] = state

	return state
}
