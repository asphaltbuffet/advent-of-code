package exercises

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

// Exercise for Advent of Code 2020 day 9.
type Exercise struct {
	common.BaseExercise
}

func parse(instr string) ([]int, error) {
	lines := strings.Split(strings.TrimSpace(instr), "\n")
	nums := make([]int, 0, len(lines))
	for _, line := range lines {
		n, err := strconv.Atoi(strings.TrimSpace(line))
		if err != nil {
			return nil, fmt.Errorf("parsing %q: %w", line, err)
		}
		nums = append(nums, n)
	}
	return nums, nil
}

// windowFor picks the preamble size: 5 for the small example, 25 for the full
// input. The example has 20 numbers; the real input has hundreds.
func windowFor(nums []int) int {
	if len(nums) <= 25 {
		return 5
	}
	return 25
}

// hasPairSum reports whether two distinct entries of window sum to target.
func hasPairSum(window []int, target int) bool {
	seen := map[int]bool{}
	for _, v := range window {
		if seen[target-v] {
			return true
		}
		seen[v] = true
	}
	return false
}

// firstInvalid returns the first number that is not a sum of two of the preceding
// `w` numbers, and its index.
func firstInvalid(nums []int, w int) (int, int) {
	for i := w; i < len(nums); i++ {
		if !hasPairSum(nums[i-w:i], nums[i]) {
			return nums[i], i
		}
	}
	return 0, -1
}

// One returns the first number that breaks the XMAS pair-sum rule.
func (e Exercise) One(instr string) (any, error) {
	nums, err := parse(instr)
	if err != nil {
		return nil, err
	}
	val, _ := firstInvalid(nums, windowFor(nums))
	return fmt.Sprintf("%d", val), nil
}

// Two finds the contiguous range summing to the part-one target and returns the
// sum of that range's smallest and largest members (the encryption weakness).
func (e Exercise) Two(instr string) (any, error) {
	nums, err := parse(instr)
	if err != nil {
		return nil, err
	}
	target, _ := firstInvalid(nums, windowFor(nums))

	// Sliding window: grow the right edge, shrink the left while over target.
	lo, sum := 0, 0
	for hi := 0; hi < len(nums); hi++ {
		sum += nums[hi]
		for sum > target && lo < hi {
			sum -= nums[lo]
			lo++
		}
		if sum == target && hi > lo {
			min, max := nums[lo], nums[lo]
			for _, v := range nums[lo : hi+1] {
				if v < min {
					min = v
				}
				if v > max {
					max = v
				}
			}
			return fmt.Sprintf("%d", min+max), nil
		}
	}

	return nil, fmt.Errorf("no contiguous range sums to %d", target)
}
