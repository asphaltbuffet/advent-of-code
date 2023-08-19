package utilities

import (
	"unicode"
)

// TODO: this fails on hyphenated words: Intern-Elves -> Intern- Elves
func CamelToTitle(x string) string {
	var out string
	for i, char := range x {
		if i == 0 {
			out += string(unicode.ToUpper(char))
		} else if unicode.IsUpper(char) {
			out += " " + string(char)
		} else {
			out += string(char)
		}
	}
	return out
}
