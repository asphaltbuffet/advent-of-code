package aoc22_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/asphaltbuffet/advent-of-code/cmd/aoc22"
	aoc "github.com/asphaltbuffet/advent-of-code/pkg/utilities"
)

func Test_2022(t *testing.T) {
	tests := []struct {
		name  string
		code  aoc.ExerciseFunc
		input []string
		want  string
	}{
		{"2022-1a Example", aoc22.D1P1, []string{"1000", "2000", "3000", "", "4000", "", "5000", "6000", "", "7000", "8000", "9000", "", "10000", ""}, "24000"},
		// 43554 \ 60470 \ 35853 \ 51812
		{"2022-1a Subset", aoc22.D1P1, []string{
			"20576", "21113", "1865", "",
			"2343", "3759", "4671", "3514", "6866", "4546", "3609", "6326", "5906", "5442", "5195", "5583", "2710", "",
			"16332", "2699", "3741", "7185", "5896", "",
			"2267", "3893", "2980", "2947", "3050", "4802", "3632", "3782", "3496", "2039", "5480", "4251", "1354", "4110", "3729", "",
		}, "60470"},
		{"2022-1b Example", aoc22.D1P2, []string{"1000", "2000", "3000", "", "4000", "", "5000", "6000", "", "7000", "8000", "9000", "", "10000", ""}, "45000"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.code(tt.input); got != tt.want {
				assert.Equal(t, tt.want, got)
			}
		})
	}
}
