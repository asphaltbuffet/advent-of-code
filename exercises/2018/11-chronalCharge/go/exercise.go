package exercises

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

// Exercise for Advent of Code 2018 day 11.
type Exercise struct {
	common.BaseExercise
}

const gridSize = 300

// cellPower is the fuel cell's power level at (x, y) for the given serial.
func cellPower(x, y, serial int) int {
	rackID := x + 10
	power := (rackID*y + serial) * rackID
	return (power/100)%10 - 5
}

// summedAreaTable builds a prefix-sum table so any square's total power is an O(1)
// lookup. sat[y][x] holds the sum of all cells from (1,1) to (x,y) inclusive.
func summedAreaTable(serial int) [][]int {
	sat := make([][]int, gridSize+1)
	for i := range sat {
		sat[i] = make([]int, gridSize+1)
	}

	for y := 1; y <= gridSize; y++ {
		for x := 1; x <= gridSize; x++ {
			sat[y][x] = cellPower(x, y, serial) + sat[y-1][x] + sat[y][x-1] - sat[y-1][x-1]
		}
	}

	return sat
}

// squareSum totals the size×size square whose top-left corner is (x, y).
func squareSum(sat [][]int, x, y, size int) int {
	x2, y2 := x+size-1, y+size-1
	return sat[y2][x2] - sat[y-1][x2] - sat[y2][x-1] + sat[y-1][x-1]
}

// best finds the top-left corner of the highest-power square of the given size.
func best(sat [][]int, size int) (bx, by, sum int) {
	sum = -1 << 62
	for y := 1; y <= gridSize-size+1; y++ {
		for x := 1; x <= gridSize-size+1; x++ {
			if s := squareSum(sat, x, y, size); s > sum {
				sum, bx, by = s, x, y
			}
		}
	}
	return bx, by, sum
}

func parseSerial(instr string) (int, error) {
	return strconv.Atoi(strings.TrimSpace(instr))
}

// One returns the answer to the first part of the exercise.
// Answer: 235,16
func (e Exercise) One(instr string) (any, error) {
	serial, err := parseSerial(instr)
	if err != nil {
		return nil, err
	}

	sat := summedAreaTable(serial)
	x, y, _ := best(sat, 3)

	return fmt.Sprintf("%d,%d", x, y), nil
}

// Two returns the answer to the second part of the exercise.
// Answer: 236,227,14
func (e Exercise) Two(instr string) (any, error) {
	serial, err := parseSerial(instr)
	if err != nil {
		return nil, err
	}

	sat := summedAreaTable(serial)

	// Scan every square size for the single highest-power square. The summed-area
	// table keeps each square total O(1), so the whole search is O(gridSize^3).
	bx, by, bSize, bSum := 0, 0, 0, -1<<62
	for size := 1; size <= gridSize; size++ {
		x, y, sum := best(sat, size)
		if sum > bSum {
			bSum, bx, by, bSize = sum, x, y, size
		}
	}

	return fmt.Sprintf("%d,%d,%d", bx, by, bSize), nil
}
