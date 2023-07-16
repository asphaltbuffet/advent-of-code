package runners

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"text/template"
)

const (
	golangInstallation              = "go"
	golangWrapperFilename           = "runtime-wrapper.go"
	golangWrapperExecutableFilename = "runtime-wrapper"
	golangBuildpathBase             = "github.com/asphaltbuffet/advent-of-code/exercises/%s/%s"
)

type golangRunner struct {
	dir                string
	cmd                *exec.Cmd
	wrapperFilepath    string
	executableFilepath string
	stdin              io.WriteCloser
}

func newGolangRunner(dir string) Runner {
	return &golangRunner{
		dir: dir,
	}
}

//go:embed interface/go.tmpl
var golangInterface []byte

func (g *golangRunner) Start() error {
	g.wrapperFilepath = filepath.Join(g.dir, golangWrapperFilename)
	g.executableFilepath = filepath.Join(g.dir, golangWrapperExecutableFilename)

	// windows requires .exe extension
	if runtime.GOOS == "windows" {
		g.executableFilepath += ".exe"
	}

	// determine package import path
	buildPath := fmt.Sprintf(
		golangBuildpathBase,
		filepath.Base(filepath.Dir(g.dir)),
		filepath.Base(g.dir))
	importPath := buildPath + "/go"

	// generate wrapper code from template
	var wrapperContent []byte
	{
		tpl := template.Must(template.New("").Parse(string(golangInterface)))
		b := new(bytes.Buffer)
		err := tpl.Execute(b, struct {
			ImportPath string
		}{importPath})
		if err != nil {
			return err
		}
		wrapperContent = b.Bytes()
	}

	// write wrapped code
	if err := os.WriteFile(g.wrapperFilepath, wrapperContent, 0o600); err != nil {
		return err
	}

	stderrBuffer := new(bytes.Buffer)

	//nolint:gosec // no user input
	cmd := exec.Command(
		golangInstallation,
		"build",
		"-tags", "runtime",
		"-o", g.executableFilepath,
		buildPath)

	cmd.Stderr = stderrBuffer
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("compilation failed: %s: %s", err, stderrBuffer.String())
	}

	if !cmd.ProcessState.Success() {
		return errors.New("compilation failed")
	}

	absExecPath, err := filepath.Abs(g.executableFilepath)
	if err != nil {
		return err
	}

	// run executable for exercise (wrapped)
	//nolint:gosec // no user input
	g.cmd = exec.Command(absExecPath)
	cmd.Dir = g.dir

	if stdin, err := setupBuffers(g.cmd); err != nil {
		return err
	} else {
		g.stdin = stdin
	}

	return g.cmd.Start()
}

func (g *golangRunner) Stop() error {
	if g.cmd == nil || g.cmd.Process == nil {
		return nil
	}

	return g.cmd.Process.Kill()
}

func (g *golangRunner) Cleanup() error {
	if g.wrapperFilepath != "" {
		_ = os.Remove(g.wrapperFilepath)
	}

	if g.executableFilepath != "" {
		_ = os.Remove(g.executableFilepath)
	}

	return nil
}

func (g *golangRunner) Run(task *Task) (*Result, error) {
	taskJSON, err := json.Marshal(task)
	if err != nil {
		return nil, err
	}

	_, _ = g.stdin.Write(append(taskJSON, '\n'))

	res := new(Result)

	if err := readJSONFromCommand(res, g.cmd); err != nil {
		return nil, err
	}

	return res, nil
}
