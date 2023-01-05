package aoc22

import (
	"github.com/asphaltbuffet/advent-of-code/exercises/aoc22_11"
	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

func init() { //nolint:gochecknoinits // init needed to register command
	common.NewDayCommand(2022, 11, aoc22_11.D11P1, aoc22_11.D11P2, Get2022Command())
}
