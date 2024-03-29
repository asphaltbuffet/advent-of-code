package exercises

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

// Directory represents a directory in the file system.
type Directory struct {
	Name        string
	Directories []Directory
	Size        int
	Files       []File
}

// File represents a file in the file system.
type File struct {
	Name string
	Size int
}

// Exercise for Advent of Code 2022 day 7.
type Exercise struct {
	common.BaseExercise
}

// One returns the answer to the first part of the exercise.
// incorrect: 1434712
// answer: 1243729
func (c Exercise) One(instr string) (any, error) {
	data := strings.Split(instr, "\n")
	files := mapFilesystem(data[1:]) // skip the first line

	// Sum the sizes of all directories <= 100k
	return sumSizes(files), nil
}

// Two returns the answer to the second part of the exercise.
// answer: 4443914
func (c Exercise) Two(instr string) (any, error) {
	data := strings.Split(instr, "\n")
	files := mapFilesystem(data[1:]) // skip the first line

	unused := 70000000 - files["root"]
	needed := 30000000 - unused
	minFound := files["root"]

	for _, s := range files {
		if s > needed && s < minFound {
			minFound = s
		}
	}
	// Sum the sizes of all directories <= 100k
	return minFound, nil
}

// mapFilesystem maps the filesystem from the input data.
func mapFilesystem(data []string) map[string]int {
	f := make(map[string]int)

	path := []string{"root"}

	for i, line := range data {
		tokens := strings.Split(line, " ")

		// Toss the first token if it's a $
		if tokens[0] == "$" {
			tokens = tokens[1:]
		}

		switch tokens[0] {
		case "ls":
			fallthrough
		case "dir":
			// fmt.Printf("skipping command: %s\n", tokens[0])
			continue
		case "cd": // manipulate the path
			if tokens[1] == ".." {
				// fmt.Printf("removing '%s' from path\n", path[len(path)-1:])
				path = path[:len(path)-1]
			} else {
				// fmt.Printf("adding '%s' to path\n", tokens[1])
				path = append(path, tokens[1])
			}
		default: // assume it's a file
			size, err := strconv.Atoi(tokens[0])
			if err != nil {
				fmt.Printf("error parsing line %d: %s\n", i, err)
				return nil
			}

			// fmt.Printf("adding %d to '%s'\n", size, strings.Join(path, "/"))

			// put this in a loop to add size to all directories in the path
			for i := 0; i < len(path); i++ {
				// fmt.Printf("adding %d to '%s'\n", size, strings.Join(path[:i+1], "/"))
				f[strings.Join(path[:i+1], "/")] += size
			}
		}
	}

	return f
}

// sumSizes sums the sizes of all directories <= 100k.
func sumSizes(f map[string]int) int {
	sum := 0

	for _, d := range f {
		if d <= 100000 {
			sum += d
		}
	}

	return sum
}
