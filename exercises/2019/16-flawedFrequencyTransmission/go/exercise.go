package exercises

import (
	"strings"

	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

// Exercise for Advent of Code 2019 day 16.
type Exercise struct {
	common.BaseExercise
}

var basePattern = [4]int{0, 1, 0, -1}

func fft(signal []int) []int {
	n := len(signal)
	out := make([]int, n)
	for i := 0; i < n; i++ {
		sum := 0
		for j := 0; j < n; j++ {
			p := basePattern[((j+1)/(i+1))%4]
			sum += signal[j] * p
		}
		if sum < 0 {
			sum = -sum
		}
		out[i] = sum % 10
	}
	return out
}

// One returns the answer to the first part of the exercise.
func (e Exercise) One(instr string) (any, error) {
	s := strings.TrimSpace(instr)
	signal := make([]int, len(s))
	for i, b := range []byte(s) {
		signal[i] = int(b - '0')
	}

	for phase := 0; phase < 100; phase++ {
		signal = fft(signal)
	}

	result := make([]byte, 8)
	for i := 0; i < 8; i++ {
		result[i] = byte(signal[i]) + '0'
	}
	return string(result), nil
}

// Two returns the answer to the second part of the exercise.
func (e Exercise) Two(instr string) (any, error) {
	s := strings.TrimSpace(instr)
	n := len(s)

	// Parse offset from first 7 digits
	offset := 0
	for i := 0; i < 7; i++ {
		offset = offset*10 + int(s[i]-'0')
	}

	total := n * 10000
	suffixLen := total - offset

	// Build suffix slice: digits from offset to end of repeated signal
	suffix := make([]int, suffixLen)
	for i := 0; i < suffixLen; i++ {
		suffix[i] = int(s[(offset+i)%n] - '0')
	}

	// 100 phases: suffix sum right-to-left, mod 10
	for phase := 0; phase < 100; phase++ {
		for i := suffixLen - 2; i >= 0; i-- {
			suffix[i] = (suffix[i] + suffix[i+1]) % 10
		}
	}

	result := make([]byte, 8)
	for i := 0; i < 8; i++ {
		result[i] = byte(suffix[i]) + '0'
	}
	return string(result), nil
}
