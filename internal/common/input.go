package common

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"

	"go.uber.org/multierr"
)

func checkForInputs(year, day string) error {
	_, partOneInputErr := os.Stat(filepath.Clean(filepath.Join("inputs", year, fmt.Sprintf(partOneInputFile, day))))

	return partOneInputErr
}

func readInput(f string) (s []string, err error) {
	file, err := os.Open(filepath.Clean(f))
	if err != nil {
		return nil, fmt.Errorf("opening input file: %w", err)
	}

	defer func() {
		// ref: https://pkg.go.dev/github.com/uber-go/multierr#Append
		err = multierr.Append(err, file.Close())
	}()

	lines := []string{}

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		n := scanner.Text()

		lines = append(lines, n)
	}

	err = scanner.Err()
	if err != nil {
		return nil, fmt.Errorf("unable to read input file: %w", err)
	}

	return lines, nil
}
