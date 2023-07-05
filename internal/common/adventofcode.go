package common

import (
	"errors"
	"fmt"
	"path/filepath"
)

type BaseExercise struct{}

func (c BaseExercise) One(instr string) (interface{}, error) {
	return nil, errors.New("not implemented")
}

func (c BaseExercise) Two(instr string) (interface{}, error) {
	return nil, errors.New("not implemented")
}

func (c BaseExercise) Vis(instr string, outdir string) error {
	return errors.New("not implemented")
}

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

	partOneInput, err := ReadInput(filepath.Clean(filepath.Join("inputs", year, fmt.Sprintf(partOneInputFile, day))))
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
