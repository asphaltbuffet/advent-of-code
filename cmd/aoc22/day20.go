package aoc22

import (
	"github.com/asphaltbuffet/advent-of-code/exercises/aoc22_20"
	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

func init() { //nolint:gochecknoinits // init needed to register command
	common.NewDayCommand(2022, 20, aoc22_20.D20P1, aoc22_20.D20P2, Get2022Command())
}
