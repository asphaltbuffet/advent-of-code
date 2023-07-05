package cmd

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/fatih/color"
	au "github.com/logrusorgru/aurora" // TODO: replace with faitdh/color package

	"github.com/asphaltbuffet/advent-of-code/pkg/exercise"
	"github.com/asphaltbuffet/advent-of-code/pkg/runners"
)

var (
	passLabel       = color.New(color.FgHiGreen).Sprint("pass")
	failLabel       = color.New(color.FgHiRed).Sprint("fail")
	incompleteLabel = color.New(color.BgHiYellow).Sprint("did not complete")
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

	fmt.Print(
		color.New(color.Bold).Sprintf("Test %s: ",
			color.New(color.FgHiBlue).Sprintf("%d.%d", part, n),
		))

	passed := r.Output == testCase.Expected

	var status string
	var followUpText string
	if !r.Ok {
		status = incompleteLabel
		followUpText = "saying \"" + r.Output + "\""
	} else if passed {
		status = passLabel
	} else {
		status = failLabel
	}

	if followUpText == "" {
		followUpText = fmt.Sprintf("in %.4f seconds", r.Duration)
	}

	fmt.Print(status)
	fmt.Println(au.Gray(10, " "+followUpText))

	if !passed && r.Ok {
		fmt.Printf(" └ Expected %s, got %s\n", au.BrightBlue(testCase.Expected), au.BrightBlue(r.Output))
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

	fmt.Print(au.Bold(fmt.Sprintf("Part %d: ", au.Yellow(part))))

	if !r.Ok {
		fmt.Print(incompleteLabel)
		fmt.Println(au.Gray(10, " saying \""+r.Output+"\""))
	} else {
		fmt.Print(au.BrightBlue(r.Output))
		fmt.Println(au.Gray(10, fmt.Sprintf(" in %.4f seconds", r.Duration)))
	}
}
