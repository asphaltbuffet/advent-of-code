package exercise

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/asphaltbuffet/advent-of-code/pkg/utilities"
)

type Exercise struct {
	Number int
	Name   string
	Dir    string
}

func (c *Exercise) String() string {
	return fmt.Sprintf("%d - %s", c.Number, c.Name)
}

var exerciseDirRegexp = regexp.MustCompile(`(?m)^(\d{2})-([a-zA-Z]+)$`)

func ListingFromDir(sourceDir string) ([]*Exercise, error) {
	dirEntries, err := os.ReadDir(sourceDir)
	if err != nil {
		return nil, err
	}

	var out []*Exercise
	for _, entry := range dirEntries {
		if entry.IsDir() && exerciseDirRegexp.MatchString(entry.Name()) {
			dir := entry.Name()

			x := strings.Split(dir, "-")
			dayInt, _ := strconv.Atoi(x[0]) // error ignored because regex should have ensured this is ok
			dayTitle := utilities.CamelToTitle(x[1])
			out = append(out, &Exercise{
				Number: dayInt,
				Name:   dayTitle,
				Dir:    filepath.Join(sourceDir, dir),
			})
		}
	}

	return out, nil
}
