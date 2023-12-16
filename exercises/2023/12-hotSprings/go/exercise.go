package exercises

import (
	"strings"
	"sync"

	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

// Exercise for Advent of Code 2023 day 12.
type Exercise struct {
	common.BaseExercise
}

// One returns the answer to the first part of the exercise.
func (e Exercise) One(instr string) (any, error) {
	lines := strings.Split(instr, "\n")
	var sum int
	var wg sync.WaitGroup
	cnt := make(chan int, len(lines))

	for i, line := range lines {
		wg.Add(1)

		go func(i int, line string, cnt chan int) {
			defer wg.Done()

			r, _ := parseLine(line)

			count, _ := r.countCombinations()

			// fmt.Printf("line %d: %d\n", i, count)

			cnt <- count
		}(i, line, cnt)
	}

	wg.Wait()
	close(cnt)

	for c := range cnt {
		sum += c
	}

	return sum, nil
}

// Two returns the answer to the second part of the exercise.
func (e Exercise) Two(instr string) (any, error) {
	lines := strings.Split(instr, "\n")
	var sum int
	var wg sync.WaitGroup
	cnt := make(chan int, len(lines))

	for i, line := range lines {
		wg.Add(1)

		go func(i int, line string, cnt chan int) {
			defer wg.Done()

			r, _ := expandAndParseLine(line)

			count, _ := r.countCombinations()
			// fmt.Printf("line %d: %d\n", i, count)

			cnt <- count
		}(i, line, cnt)
	}

	wg.Wait()
	close(cnt)

	for c := range cnt {
		sum += c
	}

	return sum, nil
}
