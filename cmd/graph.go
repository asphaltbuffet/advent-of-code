package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"regexp"

	"github.com/fatih/color"
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/spf13/cobra"
)

var (
	graphCmd *cobra.Command

	output string
)

type BenchmarkData struct {
	Implementations map[string]map[string]float64 `json:"implementations"`
	Day             int                           `json:"day"`
}

type Metrics struct {
	Day            int
	Part1          float64
	Part2          float64
	Implementation string
}

func GetGraphCmd() *cobra.Command {
	if graphCmd == nil {
		graphCmd = &cobra.Command{
			Use:     "graph [flags]",
			Aliases: []string{"g"},
			Short:   "generate run-time graph for a year",
			RunE: func(cmd *cobra.Command, args []string) error {
				return runGraph()
			},
		}
	}

	graphCmd.Flags().StringVarP(&output, "output", "o", "run-times.html", "file to write output to")

	return graphCmd
}

func runGraph() error {
	if output == "" {
		output = fmt.Sprintf("%s_run-times.html", year)
	}

	exerciseDirRegex := regexp.MustCompile(`^(\d{2})-[a-zA-Z]+$`)

	var directories []string

	exercisePath := filepath.Join("exercises", year)

	files, err := os.ReadDir(exercisePath)
	if err != nil {
		return fmt.Errorf("failed to read directory: %w", err)
	}

	for _, file := range files {
		if file.IsDir() && exerciseDirRegex.MatchString(file.Name()) {
			directories = append(directories, file.Name())
		}
	}

	benchmarkFiles := make([]string, len(directories))
	for i, dir := range directories {
		benchmarkFiles[i] = filepath.Join(dir, "benchmark.json")
	}

	benchmarkData := make(map[string]map[string]float64)
	benchmarkData["Golang"] = make(map[string]float64)
	benchmarkData["Python"] = make(map[string]float64)

	metrics := []Metrics{}

	for _, filename := range benchmarkFiles {
		filepath := filepath.Join(exercisePath, filename)

		data, err := os.ReadFile(path.Clean(filepath))
		if err != nil {
			color.Yellow("Warning: missing file %q", filepath)
			continue
		}

		var benchmark BenchmarkData

		err = json.Unmarshal(data, &benchmark)
		if err != nil {
			color.HiRed("failed to unmarshal JSON: %w", err)
			continue
		}

		m := Metrics{}

		for lang, info := range benchmark.Implementations {
			m = Metrics{
				Day:            benchmark.Day,
				Part1:          info["part.1.avg"],
				Part2:          info["part.2.avg"],
				Implementation: lang,
			}
		}

		metrics = append(metrics, m)
	}

	chart := charts.NewLine()

	goPart1 := make([]opts.LineData, 25)
	goPart2 := make([]opts.LineData, 25)
	pyPart1 := make([]opts.LineData, 25)
	pyPart2 := make([]opts.LineData, 25)

	for _, m := range metrics {
		if m.Implementation == "Golang" {
			goPart1[m.Day-1] = opts.LineData{Value: m.Part1}
			goPart2[m.Day-1] = opts.LineData{Value: m.Part2}
		} else {
			pyPart1[m.Day-1] = opts.LineData{Value: m.Part1}
			pyPart2[m.Day-1] = opts.LineData{Value: m.Part2}
		}
	}

	days := []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "11", "12", "13", "14", "15", "16", "17", "18", "19", "20", "21", "22", "23", "24", "25"}
	goalMax := make([]opts.LineData, 0)
	for i := 0; i < 25; i++ {
		goalMax = append(goalMax, opts.LineData{Value: 15})
	}

	chart.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{Title: fmt.Sprintf("AoC %s run-times", year), Subtitle: "Average run-times for each day"}),
		charts.WithYAxisOpts(opts.YAxis{Type: "log", LogBase: 1}),
		charts.WithXAxisOpts(opts.XAxis{Name: "Day"}),
	)
	chart.SetXAxis(days).
		AddSeries("Golang Part 1", goPart1).
		AddSeries("Golang Part 2", goPart2).
		AddSeries("Python Part 1", pyPart1).
		AddSeries("Python Part 2", pyPart2)

	f, err := os.Create(path.Join(exercisePath, output))
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}

	err = chart.Render(f)
	if err != nil {
		return fmt.Errorf("failed to render chart: %w", err)
	}

	return nil
}
