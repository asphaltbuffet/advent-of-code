package aoc22

import (
	"github.com/asphaltbuffet/advent-of-code/exercises/aoc22_07"
	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

func init() { //nolint:gochecknoinits // init needed to register command
	common.NewDayCommand(2022, 7, aoc22_07.D7P1, aoc22_07.D7P2, Get2022Command())
}
