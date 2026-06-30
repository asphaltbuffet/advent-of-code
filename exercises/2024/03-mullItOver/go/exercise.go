package exercises

import (
	"regexp"
	"strconv"

	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

var mulRe = regexp.MustCompile(`mul\((\d{1,3}),(\d{1,3})\)`)

var instrRe = regexp.MustCompile(`mul\((\d{1,3}),(\d{1,3})\)|do\(\)|don't\(\)`)

// Exercise for Advent of Code 2024 day 3.
type Exercise struct {
	common.BaseExercise
}

// One returns the answer to the first part of the exercise.
// One returns the answer to the first part of the exercise.
func (e Exercise) One(instr string) (any, error) {
	sum := 0
	for _, m := range mulRe.FindAllStringSubmatch(instr, -1) {
		a, _ := strconv.Atoi(m[1])
		b, _ := strconv.Atoi(m[2])
		sum += a * b
	}
	return sum, nil
}

// Two returns the answer to the second part of the exercise.
// Two returns the answer to the second part of the exercise.
func (e Exercise) Two(instr string) (any, error) {
	sum := 0
	enabled := true
	for _, m := range instrRe.FindAllStringSubmatch(instr, -1) {
		switch m[0] {
		case "do()":
			enabled = true
		case "don't()":
			enabled = false
		default:
			if enabled {
				a, _ := strconv.Atoi(m[1])
				b, _ := strconv.Atoi(m[2])
				sum += a * b
			}
		}
	}
	return sum, nil
}
