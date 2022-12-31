package aoc22

import (
	"github.com/asphaltbuffet/advent-of-code/exercises/aoc22_12"
	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

func init() { //nolint:gochecknoinits // init needed to register command
	common.NewDayCommand(2022, 12, aoc22_12.D12P1, aoc22_12.D12P2, Get2022Command())
}
