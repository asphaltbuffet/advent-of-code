package aoc22

import (
	"github.com/asphaltbuffet/advent-of-code/exercises/aoc22_18"
	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

func init() { //nolint:gochecknoinits // init needed to register command
	common.NewDayCommand(2022, 18, aoc22_18.D18P1, aoc22_18.D18P2, Get2022Command())
}
