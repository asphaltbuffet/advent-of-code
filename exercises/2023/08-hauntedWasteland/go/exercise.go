package exercises

import (
	"strings"
	"sync"

	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

// Exercise for Advent of Code 2023 day 8.
type Exercise struct {
	common.BaseExercise
}

// One returns the answer to the first part of the exercise.
func (e Exercise) One(instr string) (any, error) {
	mo, ma, _ := strings.Cut(instr, "\n\n")
	moves := parseMovement(mo)
	maps := parseMap(ma)

	loc := "AAA"
	count := 0

	for {
		curMove := moves.Element

		loc = maps[loc][curMove]
		count++

		if loc == "ZZZ" {
			break
		}

		moves = moves.Next()
	}

	return count, nil
}

// Two returns the answer to the second part of the exercise.
func (e Exercise) Two(instr string) (any, error) {
	mo, ma, _ := strings.Cut(instr, "\n\n")
	moves := parseMovement(mo)
	maps := parseMap(ma)

	locs := []string{}
	for k := range maps {
		if strings.HasSuffix(k, "A") {
			locs = append(locs, k)
		}
	}

	// fmt.Printf("starting positions: %v\n", locs)

	counts := make([]int, len(locs))
	wg := sync.WaitGroup{}
	for i, loc := range locs {
		wg.Add(1)
		go func(i int, loc string) {
			curLoc := loc
			myMoves := moves
			count := 0

			defer wg.Done()

			for {
				curMove := myMoves.Element

				curLoc = maps[curLoc][curMove]
				count++

				if strings.HasSuffix(curLoc, "Z") {
					break
				}

				myMoves = myMoves.Next()
			}

			counts[i] = count
		}(i, loc)
	}

	wg.Wait()

	// fmt.Printf("counts: %v\n", counts)

	total := lcm(counts[0], counts[1], counts[2:]...)

	return total, nil
}

// greatest common divisor (gcd) via Euclidean algorithm
func gcd(a, b int) int {
	for b != 0 {
		t := b
		b = a % b
		a = t
	}

	return a
}

// least common multiple (lcm)
func lcm(a, b int, c ...int) int {
	m := a * b / gcd(a, b)

	for _, cc := range c {
		m = lcm(m, cc)
	}

	return m
}
