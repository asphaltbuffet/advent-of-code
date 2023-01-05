package aoc22

import (
	"github.com/asphaltbuffet/advent-of-code/exercises/aoc22_04"
	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

func init() { //nolint:gochecknoinits // init needed to register command
	common.NewDayCommand(2022, 4, aoc22_04.D4P1, aoc22_04.D4P2, Get2022Command())
}
