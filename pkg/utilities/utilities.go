// Package utilities contains utility functions for Advent of Code solutions.
package utilities

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"go.uber.org/multierr"
)

// TestData contains the test data for a given day.
type Data struct {
	Input    []int
	Solution int
}

// Exercise contains the exercise for a given day.
type Exercise struct {
	Year    string
	Day     string
	Test    Data
	PartOne Data
	PartTwo Data
}

var (
	testInputFile    string = "day-%s-test.txt"
	testSolutionFile string = "day-%s-test-solution.txt"
	partOneInputFile string = "day-%s-part1.txt"
	partTwoInputFile string = "day-%s-part2.txt"
)

// NewExercise creates a new exercise for a given date.
func NewExercise(year, day string) (*Exercise, error) {
	err := checkForInputs(year, day)
	if err != nil {
		return nil, fmt.Errorf("checking for input files: %w", err)
	}

	testInput, err := readInput(filepath.Join("examples", year, fmt.Sprintf(testInputFile, day)))
	if err != nil {
		return nil, fmt.Errorf("reading test input: %w", err)
	}

	testOutput, err := readInput(filepath.Clean(filepath.Join("examples", year, fmt.Sprintf(testSolutionFile, day))))
	if err != nil {
		return nil, fmt.Errorf("reading test output: %w", err)
	}

	partOneInput, err := readInput(filepath.Clean(filepath.Join("inputs", year, fmt.Sprintf(partOneInputFile, day))))
	if err != nil {
		return nil, fmt.Errorf("reading problem input: %w", err)
	}

	partTwoInput, err := readInput(filepath.Clean(filepath.Join("inputs", year, fmt.Sprintf(partTwoInputFile, day))))
	if err != nil {
		return nil, fmt.Errorf("reading problem input: %w", err)
	}

	return &Exercise{
		Year: year,
		Day:  day,
		Test: Data{
			Input:    testInput,
			Solution: testOutput[0],
		},
		PartOne: Data{
			Input:    partOneInput,
			Solution: 0,
		},
		PartTwo: Data{
			Input:    partTwoInput,
			Solution: 0,
		},
	}, nil
}

func checkForInputs(year, day string) error {
	_, testInputErr := os.Stat(filepath.Clean(filepath.Join("inputs", year, fmt.Sprintf(testInputFile, day))))
	_, testSolutionErr := os.Stat(filepath.Clean(filepath.Join("inputs", year, fmt.Sprintf(testSolutionFile, day))))
	_, partOneInputErr := os.Stat(filepath.Clean(filepath.Join("inputs", year, fmt.Sprintf(partOneInputFile, day))))
	_, partTwoInputErr := os.Stat(filepath.Clean(filepath.Join("inputs", year, fmt.Sprintf("day-%s-part2.txt", day))))

	return multierr.Combine(testInputErr, testSolutionErr, partOneInputErr, partTwoInputErr)
}

func readInput(f string) (i []int, err error) {
	file, err := os.Open(filepath.Clean(f))
	if err != nil {
		return nil, fmt.Errorf("opening input file: %w", err)
	}
	defer func() {
		// ref: https://pkg.go.dev/github.com/uber-go/multierr#Append
		err = multierr.Append(err, file.Close())
	}()

	lines := []int{}

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		n, scanErr := strconv.Atoi(scanner.Text())
		if scanErr != nil {
			return nil, fmt.Errorf("converting input to int: %w", scanErr)
		}

		lines = append(lines, n)
	}

	err = scanner.Err()
	if err != nil {
		return nil, fmt.Errorf("unable to read input file: %w", err)
	}

	return lines, nil
}
