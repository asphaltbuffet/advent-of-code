package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/AlecAivazis/survey/v2"

	"github.com/asphaltbuffet/advent-of-code/pkg/exercise"
	"github.com/asphaltbuffet/advent-of-code/pkg/runners"
)

func userSelect(question string, choices []string) (int, error) {
	var o string
	prompt := &survey.Select{
		Message: question,
		Options: choices,
	}

	err := survey.AskOne(prompt, &o)
	if err != nil {
		return 0, err
	}

	for i, x := range choices {
		if x == o {
			return i, nil
		}
	}

	return -1, nil
}

func selectYear(dir string) (string, error) {
	var opts []string

	entries, err := os.ReadDir(dir)
	if err != nil {
		return "", fmt.Errorf("reading directory for years: %w", err)
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		opts = append(opts, entry.Name())
	}

	if len(opts) == 0 {
		return "", errors.New("no years to use")
	}

	if year != "" {
		for _, x := range opts {
			if x == year {
				return filepath.Join(dir, x), nil
			}
		}

		fmt.Printf("Could not locate year %q\n", year)
	}

	var selectedYearIndex int

	if x := len(opts); x == 1 {
		selectedYearIndex = 0
	} else {
		selectedYearIndex, err = userSelect("Which year do you want to use?", opts)
		if err != nil {
			return "", err
		}
	}

	return filepath.Join(dir, opts[selectedYearIndex]), nil
}

func selectExercise(dir string) (*exercise.Exercise, error) {
	exercises, err := exercise.ListingFromDir(dir)
	if err != nil {
		return nil, fmt.Errorf("listing exercises from %q: %w", dir, err)
	}

	if len(exercises) == 0 {
		return nil, fmt.Errorf("no exercises to run in %q", dir)
	}

	if day != 0 {
		for _, ch := range exercises {
			if ch.Day == day {
				return ch, nil
			}
		}
		fmt.Printf("Could not locate day %d\n", day)
	}

	var selectedExerciseIndex int

	if x := len(exercises); x == 1 {
		selectedExerciseIndex = 0
	} else {
		var opts []string
		for _, c := range exercises {
			opts = append(opts, c.String())
		}

		selectedExerciseIndex, err = userSelect("Which exercise do you want to run?", opts)
		if err != nil {
			return nil, fmt.Errorf("selecting exercise: %w", err)
		}
	}

	return exercises[selectedExerciseIndex], nil
}

func selectImplementation(ch *exercise.Exercise) (string, error) {
	implementations, err := ch.GetImplementations()
	if err != nil {
		return "", fmt.Errorf("getting implementations in %q: %w", ch.Dir, err)
	}

	if len(implementations) == 0 {
		return "", fmt.Errorf("no implementations to use in %q", ch.Dir)
	}

	if implementation != "" {
		for _, im := range implementations {
			if strings.EqualFold(im, implementation) {
				return im, nil
			}
		}

		fmt.Printf("Could not locate implementation %#v\n", implementation)
	}

	var selectedImplementationIndex int

	if len(implementations) == 1 {
		selectedImplementationIndex = 0
	} else {
		var opts []string
		for _, i := range implementations {
			opts = append(opts, runners.RunnerNames[i])
		}

		selectedImplementationIndex, err = userSelect("Which implementation do you want to use?", opts)
		if err != nil {
			return "", fmt.Errorf("selecting implementation: %w", err)
		}
	}

	return implementations[selectedImplementationIndex], nil
}
