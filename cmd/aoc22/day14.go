package aoc22

import (
	"github.com/asphaltbuffet/advent-of-code/exercises/aoc22_14"
	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

func init() { //nolint:gochecknoinits // init needed to register command
	common.NewDayCommand(2022, 14, aoc22_14.D14P1, aoc22_14.D14P2, Get2022Command())
}
