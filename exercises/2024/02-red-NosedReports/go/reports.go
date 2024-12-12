package exercises

import (
	"fmt"
	"strconv"
	"strings"
)

type Report struct {
	values []int
}

func parse(s string) ([]Report, error) {
	var reports []Report

	for _, line := range strings.Split(s, "\n") {
		vals := strings.Fields(line)
		r := Report{
			values: make([]int, 0, len(vals)),
		}

		for i, level := range vals {
			v, err := strconv.Atoi(level)
			if err != nil {
				return nil, fmt.Errorf("failed to parse line %d: %w", i, err)
			}

			r.values = append(r.values, v)
		}
		reports = append(reports, r)

	}

	return reports, nil
}

func (r *Report) isSafe() bool {
	if r.values[0] < r.values[1] {
		for i, j := 0, 1; j < len(r.values); i, j = i+1, j+1 {
			if r.values[j]-r.values[i] < 1 || r.values[j]-r.values[i] > 3 {
				return false
			}
		}
	} else {
		for i, j := 0, 1; j < len(r.values); i, j = i+1, j+1 {
			if r.values[i]-r.values[j] < 1 || r.values[i]-r.values[j] > 3 {
				return false
			}
		}
	}

	return true
}

func makeSafe(r *Report, idx int) (int, bool) {
	fmt.Println(r.values)
	for i := 0; i < len(r.values); i++ {
		tmp := make([]int, 0, len(r.values)-1)
		tmp = append(tmp, r.values[:i]...)
		tmp = append(tmp, r.values[i+1:]...)

		fmt.Println(tmp)
	}
	return -1, true
}

func removeAt(r []int, idx int) []int {
	tmp := make([]int, len(r))
	copy(tmp, r)
	tmp = append(tmp[:idx], tmp[idx+1:]...)
	return tmp
}
