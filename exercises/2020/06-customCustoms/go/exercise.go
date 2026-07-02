package exercises

import (
	"fmt"
	"strings"

	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

// Exercise for Advent of Code 2020 day 6.
type Exercise struct {
	common.BaseExercise
}

// parse splits the input into groups, each a slice of one person's answer lines.
func parse(instr string) [][]string {
	blocks := strings.Split(strings.TrimSpace(instr), "\n\n")
	groups := make([][]string, 0, len(blocks))
	for _, block := range blocks {
		groups = append(groups, strings.Fields(block))
	}
	return groups
}

// One sums, per group, the number of questions anyone answered yes (the union of
// answers across the group).
func (e Exercise) One(instr string) (any, error) {
	total := 0
	for _, group := range parse(instr) {
		yes := map[rune]bool{}
		for _, person := range group {
			for _, q := range person {
				yes[q] = true
			}
		}
		total += len(yes)
	}
	return fmt.Sprintf("%d", total), nil
}

// Two sums, per group, the number of questions everyone answered yes (the
// intersection): count answers whose tally equals the group size.
func (e Exercise) Two(instr string) (any, error) {
	total := 0
	for _, group := range parse(instr) {
		count := map[rune]int{}
		for _, person := range group {
			for _, q := range person {
				count[q]++
			}
		}
		for _, n := range count {
			if n == len(group) {
				total++
			}
		}
	}
	return fmt.Sprintf("%d", total), nil
}
