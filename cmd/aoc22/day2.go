package aoc22

import (
	"github.com/asphaltbuffet/advent-of-code/exercises/aoc22_02"
	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

func init() { //nolint:gochecknoinits // init needed to register command
	common.NewDayCommand(2022, 2, aoc22_02.D2P1, aoc22_02.D2P2, Get2022Command())
}
