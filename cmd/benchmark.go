package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/schollz/progressbar/v3"
	"github.com/spf13/cobra"

	"github.com/asphaltbuffet/advent-of-code/pkg/exercise"
	"github.com/asphaltbuffet/advent-of-code/pkg/runners"
	"github.com/asphaltbuffet/advent-of-code/pkg/utilities"
)

var (
	benchmarkCmd *cobra.Command
	iterations   int
)

func GetBenchmarkCmd() *cobra.Command {
	if benchmarkCmd == nil {
		benchmarkCmd = &cobra.Command{
			Use:     "benchmark [flags]",
			Aliases: []string{"bench", "b"},
			Short:   "generate benchmark data for an exercise",
			RunE: func(cmd *cobra.Command, args []string) error {
				return runBenchmark(selectedExercise, exerciseInputString, iterations)
			},
		}
	}

	benchmarkCmd.Flags().IntVarP(&iterations, "number", "n", 30, "number of benchmark iterations to run")

	return benchmarkCmd
}

func makeBenchmarkID(part runners.Part, number int) string {
	if number == -1 {
		return fmt.Sprintf("benchmark.part.%d", part)
	}

	return fmt.Sprintf("benchmark.part.%d.%d", part, number)
}

func runBenchmark(selectedExercise *exercise.Exercise, input string, numberRuns int) error {
	implementations, err := selectedExercise.GetImplementations()
	if err != nil {
		return err
	}

	var valueSets []*values

	for _, implementation := range implementations {
		v, implErr := benchmarkImplementation(implementation,
			selectedExercise.Dir,
			input,
			numberRuns)
		if implErr != nil {
			return implErr
		}

		valueSets = append(valueSets, v)
	}

	// make file
	jdata := make(map[string]interface{})
	jdata["day"] = selectedExercise.Number
	jdata["dir"] = selectedExercise.Dir
	jdata["numRuns"] = numberRuns
	jdata["implementations"] = make(map[string]interface{})

	for _, vs := range valueSets {
		x := make(map[string]interface{})

		for _, v := range vs.values {
			x[v.key] = v.value
		}

		(jdata["implementations"].(map[string]interface{}))[vs.implementation] = x
	}

	fpath := filepath.Join(selectedExercise.Dir, "benchmark.json")

	fmt.Println("Writing results to", fpath)

	jBytes, err := json.MarshalIndent(jdata, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(fpath, jBytes, 0o600)
}

type values struct {
	implementation string
	values         []kv
}

type kv struct {
	key   string
	value float64
}

func benchmarkImplementation(implementation string, dir string, inputString string, numberRuns int) (*values, error) {
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

	for _, task := range tasks {
		res, err := runner.Run(task)
		if err != nil {
			_ = pb.Close()
			return nil, err
		}

		results = append(results, res)
		_ = pb.Add(1)
	}

	fmt.Println()

	var (
		p1, p2 []float64
		p1id   = makeBenchmarkID(runners.PartOne, -1)
		p2id   = makeBenchmarkID(runners.PartTwo, -1)
	)

	for _, result := range results {
		if strings.HasPrefix(result.TaskID, p1id) {
			p1 = append(p1, result.Duration)
		} else if strings.HasPrefix(result.TaskID, p2id) {
			p2 = append(p2, result.Duration)
		}
	}

	return &values{
		implementation: runners.RunnerNames[implementation],
		values: []kv{
			{"part.1.avg", utilities.MeanFloatSlice(p1)},
			{"part.1.min", utilities.MinFloatSlice(p1)},
			{"part.1.max", utilities.MaxFloatSlice(p1)},
			{"part.2.avg", utilities.MeanFloatSlice(p2)},
			{"part.2.min", utilities.MinFloatSlice(p2)},
			{"part.2.max", utilities.MaxFloatSlice(p2)},
		},
	}, nil
}
