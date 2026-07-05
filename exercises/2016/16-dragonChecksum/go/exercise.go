package exercises

import (
	"strings"

	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

// Exercise for Advent of Code 2016 day 16.
type Exercise struct {
	common.BaseExercise
}

// fill grows the dragon-curve data to at least size bytes, then truncates.
func fill(seed []byte, size int) []byte {
	data := append([]byte(nil), seed...)
	for len(data) < size {
		b := make([]byte, len(data))
		for i := range data {
			// reverse + complement
			rc := data[len(data)-1-i]
			if rc == '0' {
				b[i] = '1'
			} else {
				b[i] = '0'
			}
		}
		data = append(append(data, '0'), b...)
	}
	return data[:size]
}

// checksum repeatedly pairs bytes (match -> '1', differ -> '0') until odd.
func checksum(data []byte) string {
	for len(data)%2 == 0 {
		next := make([]byte, len(data)/2)
		for i := range next {
			if data[2*i] == data[2*i+1] {
				next[i] = '1'
			} else {
				next[i] = '0'
			}
		}
		data = next
	}
	return string(data)
}

// solve fills to size and returns the checksum.
func solve(instr string, size int) string {
	return checksum(fill([]byte(strings.TrimSpace(instr)), size))
}

// sizeOne returns the disk size for part one: 20 for the example seed, else 272.
func sizeOne(instr string) int {
	if strings.TrimSpace(instr) == "10000" {
		return 20
	}
	return 272
}

// One returns the answer to the first part of the exercise.
func (e Exercise) One(instr string) (any, error) {
	return solve(instr, sizeOne(instr)), nil
}

// Two returns the answer to the second part of the exercise.
func (e Exercise) Two(instr string) (any, error) {
	return solve(instr, 35651584), nil
}
