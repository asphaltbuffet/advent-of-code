package exercise

import (
	"fmt"
	"os"
	"strings"

	"github.com/asphaltbuffet/advent-of-code/pkg/runners"
)

// GetImplementations returns a list of available implementations for the exercise.
func (c *Exercise) GetImplementations() ([]string, error) {
	dirEntries, err := os.ReadDir(c.Dir)
	if err != nil {
		return nil, fmt.Errorf("getting implementations for exercise: %w", err)
	}

	var impls []string

	for _, entry := range dirEntries {
		if !entry.IsDir() {
			continue
		}

		if _, ok := runners.Available[strings.ToLower(entry.Name())]; ok {
			impls = append(impls, entry.Name())
		} // TODO: log warning if there are directories that don't match any known runners
	}

	return impls, nil
}
