package exercises

import (
	"strconv"
	"strings"

	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

// Exercise for Advent of Code 2017 day 2.
type Exercise struct {
	common.BaseExercise
}

// parseRows reads the spreadsheet into rows of ints. Values are whitespace-
// delimited, so this handles both the tab-separated real input and the
// space-separated examples.
func parseRows(instr string) [][]int {
	var rows [][]int
	for _, line := range strings.Split(strings.TrimSpace(instr), "\n") {
		fields := strings.Fields(line)
		if len(fields) == 0 {
			continue
		}
		row := make([]int, 0, len(fields))
		for _, f := range fields {
			n, _ := strconv.Atoi(f)
			row = append(row, n)
		}
		rows = append(rows, row)
	}
	return rows
}

// One sums each row's spread (largest minus smallest value).
func (e Exercise) One(instr string) (any, error) {
	sum := 0
	for _, row := range parseRows(instr) {
		lo, hi := row[0], row[0]
		for _, v := range row {
			if v < lo {
				lo = v
			}
			if v > hi {
				hi = v
			}
		}
		sum += hi - lo
	}
	return sum, nil
}

// Two sums the quotient of the one evenly-divisible pair in each row.
func (e Exercise) Two(instr string) (any, error) {
	sum := 0
	for _, row := range parseRows(instr) {
		for i := 0; i < len(row); i++ {
			for j := i + 1; j < len(row); j++ {
				a, b := row[i], row[j]
				if a < b {
					a, b = b, a
				}
				if b != 0 && a%b == 0 {
					sum += a / b
				}
			}
		}
	}
	return sum, nil
}
