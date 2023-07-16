package runners

type Part uint8

const (
	PartOne Part = iota + 1
	PartTwo
	Visualize
)

type Runner interface {
	Start() error
	Stop() error
	Cleanup() error
	Run(task *Task) (*Result, error)
}

type ResultOrError struct {
	Result *Result
	Error  error
}

type RunnerCreator func(dir string) Runner

var Available = map[string]RunnerCreator{
	"go": newGolangRunner,
	"py": newPythonRunner,
}

var RunnerNames = map[string]string{
	"go": "Golang",
	"py": "Python",
}
