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

// ExerciseFunc is a function that solves an exercise.
type ExerciseFunc func([]string) string

// Data contains the test data for a given day.
type Data struct {
	Input []string
}

// Exercise contains the exercise for a given day.
type Exercise struct {
	Year    string
	Day     string
	PartOne Data
	PartTwo Data
}

const (
	partOneInputFile = "day-%s-part1.txt"
	partTwoInputFile = "day-%s-part2.txt"
)

// NewExercise creates a new exercise for a given date.
func NewExercise(year, day string) (*Exercise, error) {
	err := checkForInputs(year, day)
	if err != nil {
		return nil, fmt.Errorf("checking for input files: %w", err)
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

		PartOne: Data{
			Input: partOneInput,
		},
		PartTwo: Data{
			Input: partTwoInput,
		},
	}, nil
}

func checkForInputs(year, day string) error {
	_, partOneInputErr := os.Stat(filepath.Clean(filepath.Join("inputs", year, fmt.Sprintf(partOneInputFile, day))))
	_, partTwoInputErr := os.Stat(filepath.Clean(filepath.Join("inputs", year, fmt.Sprintf(partTwoInputFile, day))))

	return multierr.Combine(partOneInputErr, partTwoInputErr)
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

// ConvertStringSliceToIntSlice converts a slice of strings to a slice of ints.
func ConvertStringSliceToIntSlice(s []string) ([]int, error) {
	out := make([]int, len(s))

	for i, v := range s {
		n, err := strconv.Atoi(v)
		if err != nil {
			return nil, fmt.Errorf("converting string to int: %w", err)
		}

		out[i] = n
	}

	return out, nil
}

// Map applies a function to each element of a slice.
// ref: https://github.com/sa-/slicefunk/blob/66981647c9612b24c7030d60edcb1215e43c4467/main.go#L3
func Map[T, U any](s []T, f func(T) U) []U {
	modified := make([]U, len(s))

	for i, v := range s {
		modified[i] = f(v)
	}

	return modified
}
