// Package utilities contains utility functions for Advent of Code solutions.
package utilities

import (
	"fmt"
	"strconv"
)

// ConvertStringSliceToIntSlice converts a slice of strings to a slice of ints.
func ConvertStringSliceToIntSlice(s []string) ([]int, error) {
	out := make([]int, 0, len(s))

	for _, v := range s {
		n, err := strconv.Atoi(v)
		if err != nil {
			return nil, fmt.Errorf("converting string to int: %w", err)
		}

		out = append(out, n)
	}

	return out, nil
}

// Map applies a function to each element of a slice.
// ref: https://github.com/sa-/slicefunk/blob/66981647c9612b24c7030d60edcb1215e43c4467/main.go#L3
func Map[T, U any](s []T, f func(T) U) []U {
	modified := make([]U, len(s))

	for i, v := range s {
		modified[i] = f(v)
	}

	return modified
}

// Filter returns a slice with only elements that match the predicate.
func Filter[T any](s []T, f func(T) bool) []T {
	r := make([]T, len(s))
	counter := 0

	for i := 0; i < len(s); i++ {
		if f(s[i]) {
			r[counter] = s[i]
			counter++
		}
	}

	return r[:counter]
}

// Unique returns a slice with only unique elements.
func Unique[T comparable](s []T) []T {
	inResult := make(map[T]bool)

	var result []T

	for _, str := range s {
		if _, ok := inResult[str]; !ok {
			inResult[str] = true

			result = append(result, str)
		}
	}

	return result
}
