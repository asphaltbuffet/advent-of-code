package aoc22

import (
	"github.com/asphaltbuffet/advent-of-code/exercises/aoc22_19"
	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

func init() { //nolint:gochecknoinits // init needed to register command
	common.NewDayCommand(2022, 19, aoc22_19.D19P1, aoc22_19.D19P2, Get2022Command())
}
