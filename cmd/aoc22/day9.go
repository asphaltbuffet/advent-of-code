package aoc22

import (
	"github.com/asphaltbuffet/advent-of-code/exercises/aoc22_09"
	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

func init() { //nolint:gochecknoinits // init needed to register command
	common.NewDayCommand(2022, 9, aoc22_09.D9P1, aoc22_09.D9P2, Get2022Command())
}
