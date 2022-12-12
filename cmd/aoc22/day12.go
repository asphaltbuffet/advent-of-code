package aoc22

import (
	"github.com/dominikbraun/graph"

	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

func init() { //nolint:gochecknoinits // init needed to register command
	common.NewDayCommand(2022, 12, D12P1, D12P2, Get2022Command())
}

// D12P1 returns the solution for 2022 day 12 part 1.
// answer:
func D12P1(data []string) string {
	g := graph.New(graph.IntHash, graph.Directed())

	start, end, err := PopulateFromInput(&g, data)

	// calculate the shortest path

	return ""
}

func PopulateFromInput(g *graph.Graph[int, int], data []string) (int, int, error) {
	panic("unimplemented")
}

// D12P2 returns the solution for 2022 day 12 part 2.
// answer:
func D12P2(data []string) string {
	return ""
}
