package exercises

import (
	"sort"
	"strconv"
	"strings"

	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

// Exercise for Advent of Code 2024 day 5.
type Exercise struct {
	common.BaseExercise
}

// parse splits the input into the ordering rule set and the list of updates.
// rules[[2]int{a,b}] = true means page a must come before page b.
func parse(instr string) (map[[2]int]bool, [][]int) {
	parts := strings.SplitN(strings.TrimRight(instr, "\n"), "\n\n", 2)

	rules := make(map[[2]int]bool)
	for _, line := range strings.Split(parts[0], "\n") {
		ab := strings.SplitN(line, "|", 2)
		a, _ := strconv.Atoi(ab[0])
		b, _ := strconv.Atoi(ab[1])
		rules[[2]int{a, b}] = true
	}

	var updates [][]int
	for _, line := range strings.Split(parts[1], "\n") {
		if line == "" {
			continue
		}
		var pages []int
		for _, s := range strings.Split(line, ",") {
			n, _ := strconv.Atoi(s)
			pages = append(pages, n)
		}
		updates = append(updates, pages)
	}
	return rules, updates
}

// ordered reports whether pages already satisfy every applicable rule.
func ordered(rules map[[2]int]bool, pages []int) bool {
	for i := 0; i < len(pages); i++ {
		for j := i + 1; j < len(pages); j++ {
			if rules[[2]int{pages[j], pages[i]}] {
				return false
			}
		}
	}
	return true
}

// One returns the answer to the first part of the exercise.
func (e Exercise) One(instr string) (any, error) {
	rules, updates := parse(instr)
	sum := 0
	for _, pages := range updates {
		if ordered(rules, pages) {
			sum += pages[len(pages)/2]
		}
	}
	return sum, nil
}

// Two returns the answer to the second part of the exercise.
func (e Exercise) Two(instr string) (any, error) {
	rules, updates := parse(instr)
	sum := 0
	for _, pages := range updates {
		if ordered(rules, pages) {
			continue
		}
		sort.SliceStable(pages, func(i, j int) bool {
			return rules[[2]int{pages[i], pages[j]}]
		})
		sum += pages[len(pages)/2]
	}
	return sum, nil
}
