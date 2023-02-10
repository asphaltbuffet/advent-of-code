package aoc22_13 //nolint:revive,stylecheck // I don't care about the package name

import (
	"encoding/json"
	"fmt"
)

// IsOrdered compares two packet elements and returns true if they are ordered (left is "smaller").
func IsOrdered(left, right any) bool {
	return compare(left, right) <= 0
}

func compare(left, right any) int {
	l, lok := left.([]any)
	r, rok := right.([]any)

	switch {
	case !lok && !rok:
		return int(left.(float64) - right.(float64))
	case !lok:
		l = []any{left}
	case !rok:
		r = []any{right}
	}

	for i := 0; i < len(l) && i < len(r); i++ {
		if result := compare(l[i], r[i]); result != 0 {
			return result
		}
	}

	return len(l) - len(r)
}

// ParsePacket processes a packet string into a Packet.
func ParsePacket(packet string) (any, error) {
	var p any

	if err := json.Unmarshal([]byte(packet), &p); err != nil {
		return nil, fmt.Errorf("parsing packet: %w", err)
	}

	return p, nil
}
