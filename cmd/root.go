// Package cmd contains all CLI commands used by the application.
package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/dustin/go-humanize"
	"github.com/spf13/cobra"

	"github.com/asphaltbuffet/advent-of-code/pkg/exercise"
	"github.com/asphaltbuffet/advent-of-code/pkg/runners"
)

const (
	exerciseDir      = "exercises"
	exerciseInfoFile = "info.json"
)

var (
	rootCmd        *cobra.Command
	year           string
	day            int
	implementation string
	testOnly       bool
	noTest         bool
	visualize      bool

	exerciseInputString string
	selectedExercise    *exercise.Exercise
	selectedYear        string
	exerciseInfo        *exercise.Info
)

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(GetRootCommand().Execute())
}

// GetRootCommand returns the root command for the CLI.
func GetRootCommand() *cobra.Command {
	if rootCmd == nil {
		rootCmd = &cobra.Command{
			Use:               "advent-of-code [command]",
			Version:           "2.0.0",
			Short:             "advent-of-code is a collection of AoC solutions",
			Long:              `advent-of-code is a collection of AoC solutions`,
			PersistentPreRunE: getExerciseData,
			RunE:              RunRootCmd,
		}
	}

	rootCmd.Flags().BoolVarP(&testOnly, "test-only", "t", false, "only run test inputs")
	rootCmd.Flags().BoolVarP(&noTest, "no-test", "x", false, "do not run test inputs")
	rootCmd.Flags().BoolVarP(&visualize, "visualize", "g", false, "generate visualization")

	rootCmd.PersistentFlags().StringVarP(&year, "year", "y", "", "AoC year to use")
	rootCmd.PersistentFlags().IntVarP(&day, "day", "d", 0, "exercise day to use")
	rootCmd.PersistentFlags().StringVarP(&implementation, "implementation", "i", "", "implementation to use")

	rootCmd.AddCommand(GetBenchmarkCmd())
	rootCmd.AddCommand(GetGraphCmd())

	return rootCmd
}

// getExerciseData is a pre-run hook that loads the exercise data.
func getExerciseData(cmd *cobra.Command, args []string) error {
	var err error
	// List and select year
	selectedYear, err = selectYear(exerciseDir)
	if err != nil {
		return err
	}

	// List and select exercises
	selectedExercise, err = selectExercise(selectedYear)
	if err != nil {
		return err
	}

	// Load info.json file
	exerciseInfo, err = exercise.LoadExerciseInfo(filepath.Join(selectedExercise.Dir, exerciseInfoFile))
	if err != nil {
		return err
	}

	// Load exercise input
	exerciseInput, err := os.ReadFile(filepath.Join(selectedExercise.Dir, exerciseInfo.InputFile))
	if err != nil {
		return err
	}

	exerciseInputString = string(exerciseInput)

	return nil
}

// RunRootCmd is the entry point for the CLI.
func RunRootCmd(cmd *cobra.Command, args []string) error {
	// List and select implementations
	selectedImplementation, err := selectImplementation(selectedExercise)
	if err != nil {
		return err
	}

	//nolint:errcheck,gosec // printing to stdout
	bold.Printf("%s-%d %s (%s)\n\n",
		year,
		selectedExercise.Day,
		selectedExercise.Name,
		runners.RunnerNames[selectedImplementation],
	)

	runner := runners.Available[selectedImplementation](selectedExercise.Dir)
	if err := runner.Start(); err != nil {
		return err
	}

	defer func() {
		_ = runner.Stop()
		_ = runner.Cleanup()
	}()

	if visualize {
		return runVisualize(runner, exerciseInputString)
	}

	fmt.Print("Running...\n\n")

	if !noTest {
		if err := runTests(runner, exerciseInfo); err != nil {
			return err
		}
	}

	if !testOnly {
		if err := runMainTasks(runner, exerciseInputString); err != nil {
			return err
		}
	}

	return nil
}

func runVisualize(runner runners.Runner, exerciseInputString string) error {
	id := "vis"

	// directory the runner is run in, which is the exercise directory
	r, err := runner.Run(&runners.Task{
		TaskID:    id,
		Part:      runners.Visualize,
		Input:     exerciseInputString,
		OutputDir: ".",
	})
	if err != nil {
		return err
	}

	bold.Print("Visualization: ") //nolint:errcheck,gosec // printing to stdout

	var status, followUpText string

	if !r.Ok {
		status = incompleteLabel
		followUpText = fmt.Sprintf(" saying %q", r.Output)
	} else {
		status = passLabel
	}

	if followUpText == "" {
		followUpText = fmt.Sprintf(" in %s", humanize.SI(r.Duration, "s"))
	}

	fmt.Print(status)
	dimmed.Println(followUpText) //nolint:errcheck,gosec // printing to stdout

	return nil
}
