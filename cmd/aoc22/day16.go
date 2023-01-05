package aoc22

import (
	"github.com/asphaltbuffet/advent-of-code/exercises/aoc22_16"
	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

func init() { //nolint:gochecknoinits // init needed to register command
	common.NewDayCommand(2022, 16, aoc22_16.D16P1, aoc22_16.D16P2, Get2022Command())
}
