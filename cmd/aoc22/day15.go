package aoc22

import (
	"github.com/asphaltbuffet/advent-of-code/exercises/aoc22_15"
	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

func init() { //nolint:gochecknoinits // init needed to register command
	common.NewDayCommand(2022, 15, aoc22_15.D15P1, aoc22_15.D15P2, Get2022Command())
}
