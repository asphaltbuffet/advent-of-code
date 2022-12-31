package aoc22

import (
	"github.com/asphaltbuffet/advent-of-code/exercises/aoc22_03"
	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

func init() { //nolint:gochecknoinits // init needed to register command
	common.NewDayCommand(2022, 3, aoc22_03.D3P1, aoc22_03.D3P2, Get2022Command())
}
