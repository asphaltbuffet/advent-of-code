package exercises

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

// Exercise for Advent of Code 2023 day 2.
type Exercise struct {
	common.BaseExercise
}

// One returns the answer to the first part of the exercise.
func (e Exercise) One(instr string) (any, error) {
	var sum int

	for _, line := range strings.Split(instr, "\n") {
		re := regexp.MustCompile(`(((1[3-9]|[2-9]\d)\ r)|((1[4-9]|[2-9]\d)\ g)|((1[5-9]|[2-9]\d)\ b))`)

		if !re.MatchString(line) {
			reID := regexp.MustCompile(`Game\ (\d+):`)

			matches := reID.FindStringSubmatch(line)

			id, _ := strconv.Atoi(matches[1])
			// fmt.Printf("Game %d is possible\n", id)
			sum += id
		}
	}

	return sum, nil
}

// Two returns the answer to the second part of the exercise.
func (e Exercise) Two(instr string) (any, error) {
	totalPower := 0

	for _, line := range strings.Split(instr, "\n") {
		reRed := regexp.MustCompile(`\ (\d+)\ red`)
		reds := reRed.FindAllStringSubmatch(line, -1)
		gameRed := 0
		for _, r := range reds {
			n, _ := strconv.Atoi(r[1])
			if n > gameRed {
				gameRed = n
			}
		}

		reGreen := regexp.MustCompile(`\ (\d+)\ green`)
		greens := reGreen.FindAllStringSubmatch(line, -1)
		gameGreen := 0
		for _, g := range greens {
			n, _ := strconv.Atoi(g[1])
			if n > gameGreen {
				gameGreen = n
			}
		}

		reBlue := regexp.MustCompile(`\ (\d+)\ blue`)
		blues := reBlue.FindAllStringSubmatch(line, -1)
		gameBlue := 0
		for _, b := range blues {
			n, _ := strconv.Atoi(b[1])
			if n > gameBlue {
				gameBlue = n
			}
		}

		power := gameRed * gameGreen * gameBlue

		// fmt.Printf("Game %d: red=%d, green=%d, blue=%d, power=%d\n", id, gameRed, gameGreen, gameBlue, power)

		totalPower += power
	}

	// fmt.Println("Total power:", totalPower)

	return totalPower, nil
}
