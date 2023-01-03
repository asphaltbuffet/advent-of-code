package aoc22

import (
	"github.com/asphaltbuffet/advent-of-code/exercises/aoc22_23"
	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

func init() { //nolint:gochecknoinits // init needed to register command
	common.NewDayCommand(2022, 23, aoc22_23.D23P1, aoc22_23.D23P2, Get2022Command())
}
