package aoc22

import (
	"github.com/asphaltbuffet/advent-of-code/exercises/aoc22_25"
	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

func init() { //nolint:gochecknoinits // init needed to register command
	common.NewDayCommand(2022, 25, aoc22_25.D25P1, aoc22_25.D25P2, Get2022Command())
}
