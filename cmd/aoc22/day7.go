package aoc22

import (
	"fmt"
	"strconv"
	"strings"

	// "github.com/kylelemons/godebug/pretty".

	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

func init() { //nolint:gochecknoinits // init needed to register command
	common.NewDayCommand(2022, 7, D7P1, D7P2, Get2022Command())
}

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

// D7P1 returns the solution for 2022 day 7 part 1.
// incorrect: 1434712
// answer: 1243729
func D7P1(data []string) string {
	files := MapFilesystem(data[1:]) // skip the first line

	// pretty.Print(files)

	// Sum the sizes of all directories <= 100k
	return strconv.Itoa(SumSizes(files))
}

// D7P2 returns the solution for 2022 day 7 part 2.
// answer: 4443914
func D7P2(data []string) string {
	files := MapFilesystem(data[1:]) // skip the first line

	// pretty.Print(files)

	unused := 70000000 - files["root"]
	needed := 30000000 - unused
	minFound := files["root"]

	for _, s := range files {
		if s > needed && s < minFound {
			minFound = s
		}
	}
	// Sum the sizes of all directories <= 100k
	return strconv.Itoa(minFound)
}

// MapFilesystem maps the filesystem from the input data.
func MapFilesystem(data []string) map[string]int {
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

// SumSizes sums the sizes of all directories <= 100k.
func SumSizes(f map[string]int) int {
	sum := 0

	for _, d := range f {
		if d <= 100000 {
			sum += d
		}
	}

	return sum
}
