package exercises

import (
	"crypto/md5"
	"encoding/hex"
	"strconv"

	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

// Exercise for Advent of Code 2015 day 4.
type Exercise struct {
	common.BaseExercise
}

// One returns the answer to the first part of the exercise.
// answer:
func (e Exercise) One(instr string) (any, error) {
	var i int64 = 0

	k := []byte(instr)

	for ; ; i++ {
		m := md5.Sum(strconv.AppendInt(k, i, 10)) //nolint:gosec // md5 is not used for security
		h := hex.EncodeToString(m[:])

		if h[:5] == "00000" {
			break
		}
	}

	return i, nil
}

// Two returns the answer to the second part of the exercise.
// answer:
func (e Exercise) Two(instr string) (any, error) {
	var i int64 = 0

	k := []byte(instr)

	for ; ; i++ {
		m := md5.Sum(strconv.AppendInt(k, i, 10)) //nolint:gosec // md5 is not used for security
		h := hex.EncodeToString(m[:])

		if h[:6] == "000000" {
			break
		}
	}

	return i, nil
}
