package exercises

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

// Exercise for Advent of Code 2021 day 4.
type Exercise struct {
	common.BaseExercise
}

const boardSize = 5

type board struct {
	cells  [boardSize][boardSize]int
	marked [boardSize][boardSize]bool
	won    bool
}

// mark flips the cell holding n (if any) and reports whether that completed a
// row or column.
func (b *board) mark(n int) bool {
	for r := 0; r < boardSize; r++ {
		for c := 0; c < boardSize; c++ {
			if b.cells[r][c] == n {
				b.marked[r][c] = true
				return b.wins(r, c)
			}
		}
	}

	return false
}

// wins checks only the row and column of the just-marked cell.
func (b *board) wins(row, col int) bool {
	rowFull, colFull := true, true

	for i := 0; i < boardSize; i++ {
		if !b.marked[row][i] {
			rowFull = false
		}
		if !b.marked[i][col] {
			colFull = false
		}
	}

	return rowFull || colFull
}

// unmarkedSum totals the cells that have not been marked.
func (b *board) unmarkedSum() int {
	sum := 0
	for r := 0; r < boardSize; r++ {
		for c := 0; c < boardSize; c++ {
			if !b.marked[r][c] {
				sum += b.cells[r][c]
			}
		}
	}

	return sum
}

func parse(instr string) ([]int, []*board, error) {
	blocks := strings.Split(strings.TrimSpace(instr), "\n\n")
	if len(blocks) < 2 {
		return nil, nil, fmt.Errorf("expected a draw list and at least one board")
	}

	var draws []int
	for _, s := range strings.Split(strings.TrimSpace(blocks[0]), ",") {
		n, err := strconv.Atoi(strings.TrimSpace(s))
		if err != nil {
			return nil, nil, fmt.Errorf("parsing draw %q: %w", s, err)
		}
		draws = append(draws, n)
	}

	var boards []*board
	for _, blk := range blocks[1:] {
		b := &board{}
		for r, line := range strings.Split(strings.TrimSpace(blk), "\n") {
			if r >= boardSize {
				return nil, nil, fmt.Errorf("board has more than %d rows", boardSize)
			}
			for c, f := range strings.Fields(line) {
				n, err := strconv.Atoi(f)
				if err != nil {
					return nil, nil, fmt.Errorf("parsing board cell %q: %w", f, err)
				}
				b.cells[r][c] = n
			}
		}
		boards = append(boards, b)
	}

	return draws, boards, nil
}

// One returns the answer to the first part of the exercise.
func (e Exercise) One(instr string) (any, error) {
	draws, boards, err := parse(instr)
	if err != nil {
		return nil, fmt.Errorf("parsing input: %w", err)
	}

	for _, n := range draws {
		for _, b := range boards {
			if b.mark(n) {
				return fmt.Sprintf("%d", b.unmarkedSum()*n), nil
			}
		}
	}

	return nil, fmt.Errorf("no board ever won")
}

// Two returns the answer to the second part of the exercise.
func (e Exercise) Two(instr string) (any, error) {
	draws, boards, err := parse(instr)
	if err != nil {
		return nil, fmt.Errorf("parsing input: %w", err)
	}

	lastScore := 0
	remaining := len(boards)

	for _, n := range draws {
		for _, b := range boards {
			if b.won {
				continue
			}
			if b.mark(n) {
				b.won = true
				remaining--
				if remaining == 0 {
					lastScore = b.unmarkedSum() * n
				}
			}
		}
	}

	if lastScore == 0 {
		return nil, fmt.Errorf("not all boards won")
	}

	return fmt.Sprintf("%d", lastScore), nil
}
