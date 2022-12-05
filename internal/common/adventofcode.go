package common

import (
	"fmt"
	"path/filepath"
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
	partOneInputFile = "day-%s.txt"
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

	return &Exercise{
		Year: year,
		Day:  day,

		PartOne: Data{
			Input: partOneInput,
		},
		PartTwo: Data{
			Input: partOneInput,
		},
	}, nil
}
