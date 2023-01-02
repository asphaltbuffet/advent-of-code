package aoc22

import (
	"github.com/asphaltbuffet/advent-of-code/exercises/aoc22_05"
	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

func init() { //nolint:gochecknoinits // init needed to register command
	common.NewDayCommand(2022, 5, aoc22_05.D5P1, aoc22_05.D5P2, Get2022Command())
}
