package exercises

import (
	"fmt"
	"log/slog"
	"math"
	"strconv"
	"strings"
)

type Race struct {
	ID       int
	Time     int
	Distance int
}

func parseRaces(s string) []Race {
	lines := strings.Split(s, "\n")
	if len(lines) != 2 {
		slog.Error("invalid input", slog.Any("lines", lines))
		return nil
	}

	rawTimes, _ := strings.CutPrefix(lines[0], "Time:")
	rawDists, _ := strings.CutPrefix(lines[1], "Distance:")

	races := []Race{}

	t := strings.Fields(rawTimes)
	d := strings.Fields(rawDists)

	for i := 0; i < len(t); i++ {
		tt, errT := strconv.Atoi(t[i])
		dd, errD := strconv.Atoi(d[i])
		if errT != nil || errD != nil {
			slog.Error("invalid input", slog.Any("time", errT), slog.Any("dist", errD))
			return nil
		}

		slog.Debug("new race", slog.Int("Race", i), slog.Int("time", tt), slog.Int("dist", dd))

		races = append(races, Race{ID: i, Time: tt, Distance: dd})
	}

	return races
}

func parseBigRace(s string) *Race {
	lines := strings.Split(s, "\n")
	if len(lines) != 2 {
		slog.Error("invalid input", slog.Any("lines", lines))
		return nil
	}

	rawTimes, _ := strings.CutPrefix(lines[0], "Time:")
	rawDists, _ := strings.CutPrefix(lines[1], "Distance:")

	t := strings.ReplaceAll(rawTimes, " ", "")
	d := strings.ReplaceAll(rawDists, " ", "")

	tt, errT := strconv.Atoi(t)
	dd, errD := strconv.Atoi(d)
	if errT != nil || errD != nil {
		slog.Error("invalid input", slog.Any("time", errT), slog.Any("dist", errD))
		return nil
	}

	slog.Debug("big race", slog.Int("time", tt), slog.Int("dist", dd))

	return &Race{ID: 0, Time: tt, Distance: dd}
}

func (r *Race) String() string {
	return fmt.Sprintf("Race %d", r.ID)
}

func (r *Race) CalculateDistances() (int, []int) {
	distances := make([]int, 0, r.Time)
	n := 0

	for i := 0; i <= r.Time; i++ {
		travel := i * (r.Time - i)

		distances = append(distances, travel)

		if travel > r.Distance {
			n++
		}
	}

	return n, distances
}

func (r *Race) CountFasterTimes() int {
	if r.Time <= 0 || r.Distance <= 0 {
		return -1
	}

	a := -1.0
	b := float64(r.Time)
	c := -float64(r.Distance)

	discriminant := b*b - 4*a*c

	// if discriminant is negative, there are no real roots
	if discriminant < 0 {
		return -1
	}

	r1 := (-b + math.Sqrt(discriminant)) / (2 * a)
	r2 := (-b - math.Sqrt(discriminant)) / (2 * a)

	return int(math.Abs(math.Ceil(r1) - math.Ceil(r2)))
}
