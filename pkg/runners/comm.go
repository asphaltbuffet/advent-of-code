package runners

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os/exec"
	"strings"
	"sync"
	"time"

	au "github.com/logrusorgru/aurora"
)

type Task struct {
	TaskID    string `json:"task_id"`
	Part      Part   `json:"part"`
	Input     string `json:"input"`
	OutputDir string `json:"output_dir,omitempty"`
}

type Result struct {
	TaskID   string  `json:"task_id"`
	Ok       bool    `json:"ok"`
	Output   string  `json:"output"`
	Duration float64 `json:"duration"`
}

type customWriter struct {
	pending []byte
	entries [][]byte
	mux     sync.Mutex
}

func (c *customWriter) Write(b []byte) (int, error) {
	var n int

	c.mux.Lock()
	for _, x := range b {
		if x == '\n' {
			c.entries = append(c.entries, c.pending)
			c.pending = nil
		} else {
			c.pending = append(c.pending, x)
		}
		n++
	}
	c.mux.Unlock()

	return n, nil
}

func (c *customWriter) GetEntry() ([]byte, error) {
	c.mux.Lock()
	defer c.mux.Unlock()

	if len(c.entries) == 0 {
		return nil, errors.New("no entries")
	}

	var x []byte
	x, c.entries = c.entries[0], c.entries[1:]

	return x, nil
}

func setupBuffers(cmd *exec.Cmd) (io.WriteCloser, error) {
	stdoutWriter := &customWriter{}
	cmd.Stdout = stdoutWriter
	cmd.Stderr = new(bytes.Buffer)

	return cmd.StdinPipe()
}

func checkWait(cmd *exec.Cmd) ([]byte, error) {
	//nolint:errcheck // we will handle errors in the loop
	c := cmd.Stdout.(*customWriter)

	for {
		e, err := c.GetEntry()
		if err == nil {
			return e, nil
		}

		if cmd.ProcessState != nil {
			return nil, fmt.Errorf("run failed with exit code %d: %s", cmd.ProcessState.ExitCode(), cmd.Stderr.(*bytes.Buffer).String())
		}

		time.Sleep(time.Millisecond * 10)
	}
}

func readJSONFromCommand(res interface{}, cmd *exec.Cmd) error {
	for {
		inp, err := checkWait(cmd)
		if err != nil {
			return err
		}

		err = json.Unmarshal(inp, res)
		if err != nil {
			// anything returned as an error is considered a debug message
			fmt.Printf("[%s] %v\n", au.BrightRed("DBG"), strings.TrimSpace(string(inp)))
		} else {
			break
		}
	}

	return nil
}
