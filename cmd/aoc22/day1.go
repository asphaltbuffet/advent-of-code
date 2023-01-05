package aoc22

import (
	"github.com/asphaltbuffet/advent-of-code/exercises/aoc22_01"
	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

func init() { //nolint:gochecknoinits // init needed to register command
	common.NewDayCommand(2022, 1, aoc22_01.D1P1, aoc22_01.D1P2, Get2022Command())
}
