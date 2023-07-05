package exercise

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
)

// Info contains the relative path to exercise input and the specific test case data for an exercise.
type Info struct {
	InputFile string `json:"inputFile"`
	TestCases struct {
		One []*TestCase `json:"one"`
		Two []*TestCase `json:"two"`
	} `json:"testCases"`
}

// TestCase contains the input and expected output for a test case.
type TestCase struct {
	Input    string `json:"input"`
	Expected string `json:"expected"`
}

// LoadExerciseInfo loads the input and test cases for an exercise from the given json file.
func LoadExerciseInfo(fname string) (*Info, error) {
	fcont, err := os.ReadFile(path.Clean(fname))
	if err != nil {
		return nil, fmt.Errorf("loading exercise information: %w", err)
	}

	c := new(Info)

	err = json.Unmarshal(fcont, c)
	if err != nil {
		return nil, err
	}

	return c, nil
}
