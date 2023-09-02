package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/greenpau/go-calculator"
	"github.com/schollz/progressbar/v3"
	"github.com/spf13/cobra"

	"github.com/asphaltbuffet/advent-of-code/pkg/exercise"
	"github.com/asphaltbuffet/advent-of-code/pkg/runners"
)

var (
	benchmarkCmd *cobra.Command

	iterations int
)

type BenchmarkData struct {
	Date            time.Time             `json:"date,omitempty"`
	Dir             string                `json:"dir"`
	Year            string                `json:"year,omitempty"`
	Day             int                   `json:"day"`
	Runs            int                   `json:"numRuns"`
	Implementations []*ImplementationData `json:"implementations"`
}

type ImplementationData struct {
	Name    string    `json:"name"`
	PartOne *PartData `json:"part-one"`
	PartTwo *PartData `json:"part-two,omitempty"`
}

type PartData struct {
	Mean   float64 `json:"mean"`
	Median float64 `json:"median"`
	Min    float64 `json:"min"`
	Max    float64 `json:"max"`
}

func GetBenchmarkCmd() *cobra.Command {
	if benchmarkCmd == nil {
		benchmarkCmd = &cobra.Command{
			Use:     "benchmark [flags]",
			Aliases: []string{"bench", "b"},
			Short:   "generate benchmark data for an exercise",
			PreRunE: getExerciseData,
			RunE: func(cmd *cobra.Command, args []string) error {
				return runBenchmark(selectedExercise, exerciseInputString, iterations)
			},
		}
	}

	benchmarkCmd.Flags().IntVarP(&iterations, "number", "n", 30, "number of benchmark iterations to run")
	// TODO: add flag to compare to previous benchmark data
	// TODO: add flag to compare to other implementations
	// TODO: add flag to skip writing to file
	// TODO: add flag to print calculated results to stdout

	return benchmarkCmd
}

func makeBenchmarkID(part runners.Part, number int) string {
	if number == -1 {
		return fmt.Sprintf("benchmark.part.%d", part)
	}

	return fmt.Sprintf("benchmark.part.%d.%d", part, number)
}

// ImplementationError indicates that the implementation task failed.
type ImplementationError struct {
	Impl string
}

func (e *ImplementationError) Error() string {
	return fmt.Sprintf("%s run failed", e.Impl)
}

func runBenchmark(selectedExercise *exercise.Exercise, input string, numberRuns int) error {
	implementations, err := selectedExercise.GetImplementations()
	if err != nil {
		return err
	}

	var (
		implData []*ImplementationData
		ie       *ImplementationError
	)

	for _, implementation := range implementations {
		d, implErr := benchmarkImplementation(implementation,
			selectedExercise.Dir,
			input,
			numberRuns)
		if errors.As(implErr, &ie) {
			fmt.Println()
			fmt.Printf("Skipping %s due to error: %v\n",
				runners.RunnerNames[implementation],
				implErr)

			continue
		} else if err != nil {
			return implErr
		}

		if d != nil {
			implData = append(implData, d)
		}

		fmt.Println()
	}

	benchmarkData := &BenchmarkData{
		Implementations: implData,
		Day:             selectedExercise.Day,
		Runs:            numberRuns,
		Year:            year,
		Date:            time.Now().UTC(),
		Dir:             selectedExercise.Dir,
	}

	fpath := filepath.Join(selectedExercise.Dir, "benchmark.json")

	fmt.Println("Writing results to", fpath)

	jBytes, err := json.MarshalIndent(benchmarkData, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(fpath, jBytes, 0o600)
}

func benchmarkImplementation(implementation string, dir string, inputString string, numberRuns int) (*ImplementationData, error) {
	var (
		tasks   []*runners.Task
		results []*runners.Result
	)

	runner := runners.Available[implementation](dir)

	for i := 0; i < numberRuns; i++ {
		tasks = append(tasks, &runners.Task{
			TaskID: makeBenchmarkID(runners.PartOne, i),
			Part:   runners.PartOne,
			Input:  inputString,
		}, &runners.Task{
			TaskID: makeBenchmarkID(runners.PartTwo, i),
			Part:   runners.PartTwo,
			Input:  inputString,
		})
	}

	pb := progressbar.NewOptions(
		numberRuns*2, // two parts means 2x the number of runs
		progressbar.OptionSetDescription(
			fmt.Sprintf("Running %s benchmarks", runners.RunnerNames[implementation]),
		),
	)

	if err := runner.Start(); err != nil {
		return nil, err
	}

	defer func() {
		_ = runner.Stop()
		_ = runner.Cleanup()
	}()

	// TODO: add a timecheck and if bechmarking takes too long, limit number of runs
	for _, task := range tasks {
		res, err := runner.Run(task)
		if err != nil {
			_ = pb.Close()
			return nil, err
		}

		// bad results are not recorded
		if !res.Ok {
			_ = pb.Close()
			return nil, &ImplementationError{Impl: runners.RunnerNames[implementation]}
		}

		results = append(results, res)
		_ = pb.Add(1)
	}

	p1Stats, p2Stats, err := resultsToStats(results)
	if err != nil {
		return nil, err
	}

	return &ImplementationData{
		Name:    runners.RunnerNames[implementation],
		PartOne: p1Stats,
		PartTwo: p2Stats,
	}, nil
}

func resultsToStats(results []*runners.Result) (*PartData, *PartData, error) {
	var (
		p1, p2           []float64
		p1id             = makeBenchmarkID(runners.PartOne, -1)
		p2id             = makeBenchmarkID(runners.PartTwo, -1)
		p1Stats, p2Stats *PartData
	)

	for _, result := range results {
		if strings.HasPrefix(result.TaskID, p1id) {
			p1 = append(p1, result.Duration)
		} else if strings.HasPrefix(result.TaskID, p2id) {
			p2 = append(p2, result.Duration)
		}
	}

	if len(p1) == 0 && len(p2) == 0 {
		return nil, nil, &ImplementationError{Impl: runners.RunnerNames[implementation]}
	}

	if len(p1) > 0 {
		p1calc := calculator.New(p1).RunAll()
		if p1calc == nil {
			return nil, nil, fmt.Errorf("calculating part one results: %s", runners.RunnerNames[implementation])
		}

		p1Stats = &PartData{
			Mean:   p1calc.Register.Mean,
			Min:    p1calc.Register.MinValue,
			Max:    p1calc.Register.MaxValue,
			Median: p1calc.Register.Median,
		}
	} else {
		fmt.Printf(
			"No results for %s part one\n",
			runners.RunnerNames[implementation],
		)
	}

	if len(p2) > 0 {
		p2calc := calculator.New(p2).RunAll()
		if p2calc == nil {
			return nil, nil,
				fmt.Errorf(
					"calculating part two results: %s",
					runners.RunnerNames[implementation],
				)
		}

		p2Stats = &PartData{
			Mean:   p2calc.Register.Mean,
			Min:    p2calc.Register.MinValue,
			Max:    p2calc.Register.MaxValue,
			Median: p2calc.Register.Median,
		}
	} else {
		fmt.Printf(
			"No results for %s part two\n",
			runners.RunnerNames[implementation],
		)
	}

	return p1Stats, p2Stats, nil
}
