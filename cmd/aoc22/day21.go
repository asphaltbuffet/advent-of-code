package aoc22

import (
	"github.com/asphaltbuffet/advent-of-code/exercises/aoc22_21"
	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

func init() { //nolint:gochecknoinits // init needed to register command
	common.NewDayCommand(2022, 21, aoc22_21.D21P1, aoc22_21.D21P2, Get2022Command())
}
