package aoc22

import (
	"github.com/asphaltbuffet/advent-of-code/exercises/aoc22_08"
	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

func init() { //nolint:gochecknoinits // init needed to register command
	common.NewDayCommand(2022, 8, aoc22_08.D8P1, aoc22_08.D8P2, Get2022Command())
}
