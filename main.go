// Package main is the entry point for the CLI
package main

import (
	"github.com/asphaltbuffet/advent-of-code/cmd"
	_ "github.com/asphaltbuffet/advent-of-code/cmd/aoc21"
	_ "github.com/asphaltbuffet/advent-of-code/cmd/aoc22"
)

func main() {
	cmd.Execute()
}
