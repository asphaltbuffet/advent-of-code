package exercises

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_getCalibrationValue(t *testing.T) {
	type args struct {
		line string
	}

	tests := []struct {
		name string
		args args
		want int
	}{
		{"example 1-1", args{line: "1abc2"}, 12},
		{"example 1-2", args{line: "pqr3stu8vwx"}, 38},
		{"example 1-3", args{line: "a1b2c3d4e5f"}, 15},
		{"example 1-4", args{line: "treb7uchet"}, 77},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, getCalibrationValue(tt.args.line))
		})
	}
}

func Test_part2(t *testing.T) {
	type args struct {
		line string
	}

	tests := []struct {
		name string
		args args
		want int
	}{
		{"two1nine", args{line: "two1nine"}, 29},                   // 29
		{"neightwothree", args{line: "neightwothree"}, 83},         // 83
		{"nabcone2threexyz", args{line: "nabcone2threexyz"}, 13},   // 13
		{"nxtwone3four", args{line: "nxtwone3four"}, 24},           // 24
		{"n4nineeightseven2", args{line: "n4nineeightseven2"}, 42}, // 42
		{"nzoneight234", args{line: "nzoneight234"}, 14},           // 14
		{"n7pqrstsixteen", args{line: "n7pqrstsixteen"}, 76},       // 76
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, part2(tt.args.line))
		})
	}
}

func Test_firstNumber(t *testing.T) {
	type args struct {
		s       string
		numbers []string
	}

	type wants struct {
		number int
		ok     bool
	}

	all := []string{
		"one", "two", "three", "four", "five", "six", "seven", "eight", "nine",
		"1", "2", "3", "4", "5", "6", "7", "8", "9",
	}

	tests := []struct {
		name  string
		args  args
		wants wants
	}{
		{"two1nine", args{s: "two1nine", numbers: all}, wants{number: 2, ok: true}},
		{"n4nineeightseven2", args{s: "n4nineeightseven2", numbers: all}, wants{number: 4, ok: true}},
		{"abcdefghik", args{s: "abcdefghik", numbers: all}, wants{number: -1, ok: false}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotNumber, gotOk := firstNumber(tt.args.s, tt.args.numbers)
			assert.Equal(t, tt.wants.number, gotNumber)
			assert.Equal(t, tt.wants.ok, gotOk)
		})
	}
}

func Test_lastNumber(t *testing.T) {
	type args struct {
		s       string
		numbers []string
	}

	type wants struct {
		number int
		ok     bool
	}

	all := []string{
		"one", "two", "three", "four", "five", "six", "seven", "eight", "nine",
		"1", "2", "3", "4", "5", "6", "7", "8", "9",
	}

	tests := []struct {
		name  string
		args  args
		wants wants
	}{
		{"two1nine", args{s: "two1nine", numbers: all}, wants{number: 9, ok: true}},
		{"n4nineeightseven2", args{s: "n4nineeightseven2", numbers: all}, wants{number: 2, ok: true}},
		{"abcdefghik", args{s: "abcdefghik", numbers: all}, wants{number: -1, ok: false}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotNumber, gotOk := lastNumber(tt.args.s, tt.args.numbers)
			assert.Equal(t, tt.wants.number, gotNumber)
			assert.Equal(t, tt.wants.ok, gotOk)
		})
	}
}

func Test_stringToNumber(t *testing.T) {
	type args struct {
		s string
	}

	tests := []struct {
		name      string
		args      args
		want      int
		assertion require.ErrorAssertionFunc
	}{
		{"text", args{s: "four"}, 4, require.NoError},
		{"digit", args{s: "6"}, 6, require.NoError},
		{"not a number", args{s: "fake"}, -1, require.Error},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := stringToNumber(tt.args.s)
			tt.assertion(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
