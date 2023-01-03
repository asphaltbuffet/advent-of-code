package aoc22

import (
	"github.com/asphaltbuffet/advent-of-code/exercises/aoc22_17"
	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

func init() { //nolint:gochecknoinits // init needed to register command
	common.NewDayCommand(2022, 17, aoc22_17.D17P1, aoc22_17.D17P2, Get2022Command())
}
