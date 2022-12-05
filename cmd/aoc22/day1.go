package aoc22

import (
	"sort"
	"strconv"

	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

func init() { //nolint:gochecknoinits // init needed to register command
	common.NewDayCommand(2022, 1, D1P1, D1P2, Get2022Command())
}

// D1P1 returns the solution for 2022 day 1 part 1
// answer: 70720
func D1P1(data []string) string {
	sum := 0
	calories := sort.IntSlice{}

	for i := 0; i < len(data); i++ {
		if data[i] == "" {
			calories = append(calories, sum)
			sum = 0
		} else {
			n, _ := strconv.Atoi(data[i])
			sum += n
		}
	}

	sort.Slice(calories, func(i, j int) bool {
		return calories[i] > calories[j]
	})

	return strconv.Itoa(calories[0])
}

// D1P2 returns the solution for 2022 day 1 part 2
// answer: 207148
func D1P2(data []string) string {
	sum := 0
	calories := sort.IntSlice{}

	for i := 0; i < len(data); i++ {
		if data[i] == "" {
			calories = append(calories, sum)
			sum = 0
		} else {
			n, _ := strconv.Atoi(data[i])
			sum += n
		}
	}

	sort.Slice(calories, func(i, j int) bool {
		return calories[i] > calories[j]
	})

	return strconv.Itoa(calories[0] + calories[1] + calories[2])
}
