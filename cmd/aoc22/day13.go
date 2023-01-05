package aoc22

import (
	"github.com/asphaltbuffet/advent-of-code/exercises/aoc22_13"
	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

func init() { //nolint:gochecknoinits // init needed to register command
	common.NewDayCommand(2022, 13, aoc22_13.D13P1, aoc22_13.D13P2, Get2022Command())
}
