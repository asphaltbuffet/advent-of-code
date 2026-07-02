package exercises

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

// Exercise for Advent of Code 2020 day 10.
type Exercise struct {
	common.BaseExercise
}

// chain returns the sorted joltages including the 0-jolt outlet at the front and
// the device (max + 3) at the end.
func chain(instr string) ([]int, error) {
	lines := strings.Split(strings.TrimSpace(instr), "\n")
	nums := make([]int, 0, len(lines)+2)
	nums = append(nums, 0) // charging outlet
	for _, line := range lines {
		n, err := strconv.Atoi(strings.TrimSpace(line))
		if err != nil {
			return nil, fmt.Errorf("parsing %q: %w", line, err)
		}
		nums = append(nums, n)
	}
	sort.Ints(nums)
	nums = append(nums, nums[len(nums)-1]+3) // device
	return nums, nil
}

// One returns the count of 1-jolt differences times the count of 3-jolt
// differences across the full adapter chain.
func (e Exercise) One(instr string) (any, error) {
	nums, err := chain(instr)
	if err != nil {
		return nil, err
	}

	var ones, threes int
	for i := 1; i < len(nums); i++ {
		switch nums[i] - nums[i-1] {
		case 1:
			ones++
		case 3:
			threes++
		}
	}

	return fmt.Sprintf("%d", ones*threes), nil
}

// Two counts the distinct arrangements of adapters that connect the outlet to the
// device. ways[i] accumulates the paths reaching adapter i from any of the up-to-
// three predecessors within 3 jolts.
func (e Exercise) Two(instr string) (any, error) {
	nums, err := chain(instr)
	if err != nil {
		return nil, err
	}

	ways := make([]int, len(nums))
	ways[0] = 1
	for i := 1; i < len(nums); i++ {
		for j := i - 1; j >= 0 && nums[i]-nums[j] <= 3; j-- {
			ways[i] += ways[j]
		}
	}

	return fmt.Sprintf("%d", ways[len(ways)-1]), nil
}
