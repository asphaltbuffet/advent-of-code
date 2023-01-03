package aoc22

import (
	"github.com/asphaltbuffet/advent-of-code/exercises/aoc22_06"
	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

func init() { //nolint:gochecknoinits // init needed to register command
	common.NewDayCommand(2022, 6, aoc22_06.D6P1, aoc22_06.D6P2, Get2022Command())
}
