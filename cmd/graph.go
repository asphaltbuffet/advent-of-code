package cmd

import (
	"encoding/json"
	"fmt"
	"image/color"
	"io/fs"
	"math"
	"os"
	"path/filepath"

	"github.com/dustin/go-humanize"
	"github.com/spf13/cobra"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
	"gonum.org/v1/plot/vg/vgimg"
)

var (
	graphCmd *cobra.Command

	outfile string

	langColor = map[string]color.Color{
		"Golang": color.RGBA{R: 0, G: 173, B: 216, A: 255},
		"Python": color.RGBA{R: 55, G: 118, B: 171, A: 255},
	}
)

func GetGraphCmd() *cobra.Command {
	if graphCmd == nil {
		graphCmd = &cobra.Command{
			Use:   "graph year [flags]",
			Args:  cobra.ExactArgs(1),
			Short: "generate run-time graph for a year",
			PersistentPreRun: func(cmd *cobra.Command, args []string) {
				fmt.Println("gathering data...")
			},
			RunE: func(cmd *cobra.Command, args []string) error {
				return runGraph(args[0])
			},
		}
	}

	graphCmd.Flags().StringVarP(&outfile, "output", "o", "run-times.png", "file to write output to")

	return graphCmd
}

func getBenchmarkFilesByYear(year string) ([]string, error) {
	benchFiles := []string{}

	// use filepath.walk to get all benchmark.json files
	err := filepath.WalkDir(filepath.Join("exercises", year), func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil
		}

		if filepath.Base(path) == "benchmark.json" {
			benchFiles = append(benchFiles, path)
		}

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("searching for benchmark files: %w", err)
	}

	return benchFiles, nil
}

func runGraph(year string) error {
	files, err := getBenchmarkFilesByYear(year)
	if err != nil {
		return fmt.Errorf("getting benchmark files: %w", err)
	}

	fmt.Printf("found %d benchmark files:\n", len(files))
	benchData := make([]*BenchmarkData, len(files))

	for i, bf := range files {
		var data *BenchmarkData

		data, err = readBenchmarkFile(bf)
		if err != nil {
			return fmt.Errorf("reading benchmark file: %w", err)
		}

		data.Year = year
		benchData[i] = data
	}

	err = generateGraph(benchData, outfile)
	if err != nil {
		return fmt.Errorf("generating graph: %w", err)
	}

	fmt.Printf("wrote %s graph to %s\n", year, outfile)

	return nil
}

func readBenchmarkFile(path string) (*BenchmarkData, error) {
	var bd BenchmarkData

	f, err := os.ReadFile(filepath.Clean(path))
	if err != nil {
		return nil, fmt.Errorf("reading file: %w", err)
	}

	err = json.Unmarshal(f, &bd)
	if err != nil {
		return nil, fmt.Errorf("unmarshalling json: %w", err)
	}

	return &bd, nil
}

func mapBenchmarkData(benchData []*BenchmarkData) map[string]map[int]plotter.XYs {
	langs := map[string]bool{}
	days := map[int]bool{}

	// initialize maps to store data more effectively for plotting
	for _, bd := range benchData {
		for _, impl := range bd.Implementations {
			langs[impl.Name] = true
			days[bd.Day] = true
		}
	}

	dataMap := make(map[string]map[int]plotter.XYs, len(langs))

	for lang := range langs {
		dataMap[lang] = make(map[int]plotter.XYs, 2)
		dataMap[lang][0] = plotter.XYs{}
		dataMap[lang][1] = plotter.XYs{}
	}

	for _, bd := range benchData {
		for _, impl := range bd.Implementations {
			dataMap[impl.Name][0] = append(dataMap[impl.Name][0], plotter.XY{
				X: float64(bd.Day),
				Y: impl.PartOne.Mean,
			})

			if impl.PartTwo == nil {
				continue
			}

			dataMap[impl.Name][1] = append(dataMap[impl.Name][1],
				plotter.XY{
					X: float64(bd.Day),
					Y: impl.PartTwo.Mean,
				})
		}
	}

	return dataMap
}

func generateGraph(benchData []*BenchmarkData, outfile string) error {
	plots, err := NewBenchmarkPlots(benchData[0].Year)
	if err != nil {
		return fmt.Errorf("creating plots: %w", err)
	}

	dataMap := mapBenchmarkData(benchData)

	for lang, parts := range dataMap {
		for part, xys := range parts {
			var (
				ln *plotter.Line
				pt *plotter.Scatter
			)

			ln, pt, err = plotter.NewLinePoints(xys)
			if err != nil {
				return fmt.Errorf("filling %s part %d plot: %w", lang, part, err)
			}

			ln.Color = langColor[lang]
			pt.Shape = draw.CircleGlyph{}
			pt.Color = langColor[lang]

			plots[0][part].Add(ln, pt)
			plots[0][part].Legend.Add(lang, ln, pt)
		}
	}

	// make sure both plots have the same Y axis for alignment
	max := max(plots[0][0].Y.Max, plots[0][1].Y.Max, 60)
	plots[0][0].Y.Max = max
	plots[0][1].Y.Max = max

	min := min(plots[0][0].Y.Min, plots[0][1].Y.Min)
	plots[0][0].Y.Min = min
	plots[0][1].Y.Min = min

	img := vgimg.NewWith(vgimg.UseWH(12.5*vg.Inch, 5*vg.Inch), vgimg.UseDPI(300))
	dc := draw.New(img)

	const rows, cols = 1, 2

	t := draw.Tiles{
		Rows:      rows,
		Cols:      cols,
		PadX:      vg.Points(20),
		PadRight:  vg.Points(10),
		PadLeft:   vg.Points(10),
		PadBottom: vg.Points(10),
		PadTop:    vg.Points(10),
	}

	canvases := plot.Align(plots, t, dc)

	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			if plots[r][c] != nil {
				plots[r][c].Draw(canvases[r][c])
			}
		}
	}

	path := filepath.Join("exercises", benchData[0].Year, outfile)
	fmt.Printf("writing graph to %s\n", path)

	w, err := os.Create(filepath.Clean(path))
	if err != nil {
		return fmt.Errorf("creating image file: %w", err)
	}

	png := vgimg.PngCanvas{Canvas: img}
	if _, err := png.WriteTo(w); err != nil {
		return fmt.Errorf("writing image file: %w", err)
	}

	return nil
}

func NewBenchmarkPlots(year string) ([][]*plot.Plot, error) {
	const rows, cols = 1, 2
	plots := make([][]*plot.Plot, rows)

	for j := 0; j < rows; j++ {
		plots[j] = make([]*plot.Plot, cols)

		for i := 0; i < cols; i++ {
			p := plot.New()

			p.X.Label.Text = "Day"

			// p.Y.Label.Text = "Running time (seconds)"
			p.Y.Tick.Marker = HumanizedLogTicks{}
			p.X.Tick.Marker = plot.TickerFunc(func(min, max float64) []plot.Tick {
				ticks := []plot.Tick{}

				for i := min; i <= max; i++ {
					ticks = append(
						ticks,
						plot.Tick{
							Value: i,
							Label: fmt.Sprintf("%d", int(i)),
						},
					)
				}

				return ticks
			})
			p.Y.Scale = plot.LogScale{}
			p.Y.Min = 0.000001
			// part1Plot.Y.Max = +10
			// part1Plot.X.Label.Position = draw.PosRight
			// part1Plot.Y.Label.Position = draw.PosTop

			plots[j][i] = p
		}
	}

	part1Plot := plots[0][0]
	part2Plot := plots[0][1]

	part1Plot.Title.Text = fmt.Sprintf(
		"Average Exercise Running Time\nAdvent of Code %s: Part One",
		year)
	part2Plot.Title.Text = fmt.Sprintf(
		"Average Exercise Running Time\nAdvent of Code %s: Part Two",
		year)

	g := plotter.NewGrid()
	g.Vertical.Color = color.Transparent
	part1Plot.Add(g)
	part2Plot.Add(g)

	redline := plotter.NewFunction(func(x float64) float64 { return 15 })
	redline.Color = color.RGBA{R: 255, G: 0, B: 0, A: 255}
	redline.Dashes = plotutil.Dashes(2)
	part1Plot.Add(redline)
	part2Plot.Add(redline)

	return plots, nil
}

// HumanizedLogTicks is suitable for the Tick.Marker field of an Axis,
// it returns tick marks suitable for a log-scale axis which have been
// humanized.
type HumanizedLogTicks struct {
	// Prec specifies the precision of tick rendering
	// according to the documentation for strconv.FormatFloat.
	Prec int
}

var _ plot.Ticker = HumanizedLogTicks{}

// Ticks returns Ticks in a specified range
func (t HumanizedLogTicks) Ticks(min, max float64) []plot.Tick {
	if min <= 0 || max <= 0 {
		panic("Values must be greater than 0 for a log scale.")
	}

	val := math.Pow10(int(math.Log10(min)))
	max = math.Pow10(int(math.Ceil(math.Log10(max))))

	var ticks []plot.Tick

	for val < max {
		for i := 1; i < 10; i++ {
			if i == 1 {
				ticks = append(
					ticks,
					plot.Tick{
						Value: val,
						Label: humanize.SIWithDigits(val, 0, "s"),
					})
			}

			ticks = append(ticks, plot.Tick{Value: val * float64(i)})
		}

		val *= 10
	}

	ticks = append(ticks,
		plot.Tick{
			Value: val,
			Label: humanize.SIWithDigits(val, 0, "s"),
		})

	return ticks
}
