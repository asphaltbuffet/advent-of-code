package aoc22

import (
	"github.com/asphaltbuffet/advent-of-code/exercises/aoc22_10"
	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

func init() { //nolint:gochecknoinits // init needed to register command
	common.NewDayCommand(2022, 10, aoc22_10.D10P1, aoc22_10.D10P2, Get2022Command())
}
