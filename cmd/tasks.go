package cmd

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/dustin/go-humanize"
	"github.com/fatih/color"

	"github.com/asphaltbuffet/advent-of-code/pkg/exercise"
	"github.com/asphaltbuffet/advent-of-code/pkg/runners"
)

var (
	passLabel       = color.New(color.FgHiGreen).Sprint("pass")
	failLabel       = color.New(color.FgHiRed).Sprint("fail")
	incompleteLabel = color.New(color.BgHiYellow).Sprint("did not complete")
	missingLabel    = color.New(color.FgHiYellow, color.Italic).Sprint("empty")

	bold       = color.New(color.Bold)
	dimmed     = color.New(color.FgHiBlack, color.Italic)
	brightBlue = color.New(color.FgHiBlue)
	boldBlue   = color.New(color.Bold, color.FgHiBlue)
	boldYellow = color.New(color.Bold, color.FgHiYellow)
)

func makeTestID(part runners.Part, n int) string {
	return fmt.Sprintf("test.%d.%d", part, n)
}

func parseTestID(x string) (runners.Part, int) {
	y := strings.Split(x, ".")
	p, _ := strconv.Atoi(y[1])
	n, _ := strconv.Atoi(y[2])

	return runners.Part(p), n
}

func makeMainID(part runners.Part) string {
	return fmt.Sprintf("main.%d", part)
}

func parseMainID(x string) runners.Part {
	y := strings.Split(x, ".")
	p, _ := strconv.Atoi(y[1])

	return runners.Part(p)
}

func runTests(runner runners.Runner, info *exercise.Info) error {
	for i, testCase := range info.TestCases.One {
		id := makeTestID(runners.PartOne, i)

		if testCase.Input == "" && testCase.Expected == "" {
			handleTestResult(&runners.Result{
				TaskID: id,
				Ok:     false,
				Output: "empty input or expected output",
			}, testCase)

			continue
		}

		result, err := runner.Run(&runners.Task{
			TaskID: id,
			Part:   runners.PartOne,
			Input:  testCase.Input,
		})
		if err != nil {
			return err
		}

		handleTestResult(result, testCase)
	}

	for i, testCase := range info.TestCases.Two {
		id := makeTestID(runners.PartTwo, i)

		if testCase.Input == "" && testCase.Expected == "" {
			handleTestResult(&runners.Result{
				TaskID: id,
				Ok:     false,
				Output: "empty input or expected output",
			}, testCase)

			continue
		}

		result, err := runner.Run(&runners.Task{
			TaskID: id,
			Part:   runners.PartTwo,
			Input:  testCase.Input,
		})
		if err != nil {
			return err
		}

		handleTestResult(result, testCase)
	}

	return nil
}

func handleTestResult(r *runners.Result, testCase *exercise.TestCase) {
	part, n := parseTestID(r.TaskID)

	bold.Print("Test ")               //nolint:errcheck,gosec // printing to stdout
	boldBlue.Printf("%d.%d", part, n) //nolint:errcheck,gosec // printing to stdout
	bold.Print(": ")                  //nolint:errcheck,gosec // printing to stdout

	passed := r.Output == testCase.Expected
	missing := testCase.Input == "" && testCase.Expected == ""

	var status, followUpText string

	switch {
	case missing:
		status = missingLabel

	case !r.Ok:
		status = incompleteLabel
		followUpText = fmt.Sprintf(" saying %q", r.Output)

	case passed:
		status = passLabel

	default:
		status = failLabel
	}

	if followUpText == "" && !missing {
		followUpText = fmt.Sprintf(" in %s", humanize.SIWithDigits(r.Duration, 1, "s"))
	}

	fmt.Print(status)
	dimmed.Println(followUpText) //nolint:errcheck,gosec // printing to stdout

	if !passed && r.Ok {
		fmt.Printf(" └ Expected %s, got %s\n", brightBlue.Sprint(testCase.Expected), brightBlue.Sprint(r.Output))
	}
}

func runMainTasks(runner runners.Runner, input string) error {
	for part := runners.PartOne; part <= runners.PartTwo; part += 1 {
		id := makeMainID(part)

		result, err := runner.Run(&runners.Task{
			TaskID: id,
			Part:   part,
			Input:  input,
		})
		if err != nil {
			return err
		}

		handleMainResult(result)
	}

	return nil
}

func handleMainResult(r *runners.Result) {
	part := parseMainID(r.TaskID)

	bold.Print("Part ")             //nolint:errcheck,gosec // printing to stdout
	boldYellow.Printf("%d: ", part) //nolint:errcheck,gosec // printing to stdout

	if !r.Ok {
		fmt.Print(incompleteLabel)
		dimmed.Printf(" saying %q\n", r.Output) //nolint:errcheck,gosec // printing to stdout
	} else {
		brightBlue.Print(r.Output)                                           //nolint:errcheck,gosec // printing to stdout
		dimmed.Printf(" in %s\n", humanize.SIWithDigits(r.Duration, 1, "s")) //nolint:errcheck,gosec // printing to stdout
	}
}
