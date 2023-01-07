package aoc22_25 //nolint:revive,stylecheck // I don't care about the package name

import (
	"strconv"
)

// Decode converts a base 5 SNAFU string to an integer in base 10.
func Decode(s string) int {
	rev := Reverse([]rune(s))
	n := 0

	for i, r := range rev {
		switch r {
		case '-':
			n += -1 * Pow(5, i)
		case '=':
			n += -2 * Pow(5, i)
		case '0':
			n += 0
		case '1':
			n += 1 * Pow(5, i)
		case '2':
			n += 2 * Pow(5, i)
		default:
			return 0
		}
	}

	return n
}

// Encode converts an integer in base 10 to a base 5 SNAFU string.
func Encode(n int) string {
	if n == 0 {
		return "0"
	}

	rs := make(map[int]int, 0)

	// convert to little-endian base 5.
	// for i := 0; n != 0; i++ {
	// 	r := n % 5
	// 	n /= 5

	// 	rs[i] = r
	// }
	b := strconv.FormatInt(int64(n), 5)
	for i, d := range b {
		rs[len(b)-1-i], _ = strconv.Atoi(string(d))
	}

	s := []byte{}

	// convert to SNAFU.
	for i := 0; i < len(rs); i++ {
		switch rs[i] {
		case 4:
			s = append(s, '-')
			rs[i+1]++
		case 3:
			s = append(s, '=')
			rs[i+1]++
		case 2:
			s = append(s, '2')
		case 1:
			s = append(s, '1')
		case 0:
			s = append(s, '0')
		case 5: // carry digits
			s = append(s, '0')
			rs[i+1]++
		}
	}

	return string(Reverse(s))
}

// Reverse returns a new slice with the elements of s in Reverse order.
func Reverse[T any](s []T) []T {
	r := s
	for i, j := 0, len(r)-1; i <= j; i, j = i+1, j-1 {
		r[i], r[j] = r[j], r[i]
	}

	return r
}

// Pow returns x^n, the base-x exponential of n.
func Pow(x, n int) int {
	switch {
	case n == 0:
		return 1
	case n%2 == 0:
		return Pow(x*x, n/2)
	default:
		return x * Pow(x*x, (n-1)/2)
	}
}
