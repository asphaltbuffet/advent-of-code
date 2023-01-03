package aoc22

import (
	"github.com/asphaltbuffet/advent-of-code/exercises/aoc22_24"
	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

func init() { //nolint:gochecknoinits // init needed to register command
	common.NewDayCommand(2022, 24, aoc22_24.D24P1, aoc22_24.D24P2, Get2022Command())
}
