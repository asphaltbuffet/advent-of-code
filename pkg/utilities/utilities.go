// Package utilities contains utility functions for Advent of Code solutions.
package utilities

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
)

// TestData contains the test data for a given day.
type TestData struct {
	Input  []int
	Output int
}

// Exercise contains the exercise for a given day.
type Exercise struct {
	Test  TestData
	Input []int
}

// ReadTestData reads the test data for the given day.
func ReadTestData(year, day string) (*Exercise, error) {
	testInput, err := readInput(filepath.Join("examples", year, fmt.Sprintf("day-%s-input.txt", day)))
	if err != nil {
		return nil, fmt.Errorf("reading test input: %w", err)
	}

	testOutput, err := readInput(filepath.Clean(filepath.Join("examples", year, fmt.Sprintf("day-%s-solution.txt", day))))
	if err != nil {
		return nil, fmt.Errorf("reading test output: %w", err)
	}

	problemInput, err := readInput(filepath.Clean(filepath.Join("inputs", year, fmt.Sprintf("day-%s.txt", day))))
	if err != nil {
		return nil, fmt.Errorf("reading problem input: %w", err)
	}

	return &Exercise{
		Test: TestData{
			Input:  testInput,
			Output: testOutput[0],
		},
		Input: problemInput,
	}, nil
}

func readInput(f string) ([]int, error) {
	file, err := os.Open(filepath.Clean(f))
	if err != nil {
		return nil, fmt.Errorf("opening input file: %w", err)
	}
	defer file.Close() //nolint: errcheck, gosec // not important right now

	lines := []int{}

	scanner := bufio.NewScanner(file)
	// optionally, resize scanner's capacity for lines over 64K, see next example
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
