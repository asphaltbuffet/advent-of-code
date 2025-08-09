package logs

import (
	"unicode/utf8"
)

// Application identifies the application emitting the given log.
func Application(log string) string {
	for _, r := range log {
		switch {
		case r == '❗':
			return "recommendation"

		case r == '🔍':
			return "search"

		case r == '☀':
			return "weather"
		}
	}

	return "default"
}

// Replace replaces all occurrences of old with new, returning the modified log
// to the caller.
func Replace(log string, oldRune, newRune rune) string {
	var fixed string

	for _, r := range log {
		if r == oldRune {
			fixed += string(newRune)
		} else {
			fixed += string(r)
		}
	}

	return fixed
}

// WithinLimit determines whether or not the number of characters in log is
// within the limit.
func WithinLimit(log string, limit int) bool {
	return utf8.RuneCountInString(log) <= limit
}
