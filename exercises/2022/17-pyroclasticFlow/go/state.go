package exercises

import (
	"fmt"
	"strings"
)

type state struct {
	chamber       [][]string
	settledHeight int
	curShape      []point
	nextShape     shape
	jets          []string
	jetNum        int
}

func (s *state) moveRight() {
	// check if can move right
	canMoveRight := true

	for _, c := range s.curShape {
		if c.x == 6 || s.chamber[c.y][c.x+1] != space {
			canMoveRight = false
		}
	}

	if canMoveRight {
		for i := range s.curShape {
			s.curShape[i].x++
		}
	}
}

func (s *state) moveLeft() {
	// check if can move left
	canMoveLeft := true

	for _, c := range s.curShape {
		if c.x == 0 || s.chamber[c.y][c.x-1] != space {
			canMoveLeft = false
		}
	}

	if canMoveLeft {
		for i := range s.curShape {
			s.curShape[i].x--
		}
	}
}

func (s *state) moveDown() {
	// move down
	canMoveDown := true

	for _, c := range s.curShape {
		if c.y == 0 || s.chamber[c.y-1][c.x] != space {
			canMoveDown = false
		}
	}
	// is blocked, draw onto settled then make nil
	if canMoveDown {
		for i := range s.curShape {
			s.curShape[i].y--
		}

		return
	}

	for _, c := range s.curShape {
		s.chamber[c.y][c.x] = rock
	}

	s.curShape = nil

	for r := len(s.chamber) - 1; r >= 0; r-- {
		if strings.Join(s.chamber[r], "") != "......." {
			s.settledHeight = r
			break
		}
	}
}

func (s *state) moveRock() {
	s.populateNextBaseCoords()

	highestRow := 0
	for _, c := range s.curShape {
		if c.y > highestRow {
			highestRow = c.y
		}
	}

	for len(s.chamber) <= highestRow {
		s.chamber = append(s.chamber, newEmptyRow())
	}

	// will be set back to nil when settled
	for s.curShape != nil {
		switch s.jets[s.jetNum%len(s.jets)] {
		case ">":
			s.moveRight()
		case "<":
			s.moveLeft()
		default:
			panic(s.jets[s.jetNum])
		}

		s.jetNum++

		s.moveDown()
	}
}

func newEmptyRow() []string {
	row := make([]string, 7)

	for i := range row {
		row[i] = "."
	}

	return row
}

func (s *state) populateNextBaseCoords() {
	copyCoords := make([]point, len(shapes[s.nextShape]))
	copy(copyCoords, shapes[s.nextShape])
	s.nextShape++
	s.nextShape %= 5

	for i := range copyCoords {
		copyCoords[i].y += s.settledHeight + 1 + 3
	}

	s.curShape = copyCoords
}

func (s *state) hash(topRowsToHash int) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprint(s.nextShape))

	for r := s.settledHeight; r >= 0 && r > s.settledHeight-topRowsToHash; r-- {
		sb.WriteString("\n" + strings.Join(s.chamber[r], ""))
	}

	return sb.String()
}
