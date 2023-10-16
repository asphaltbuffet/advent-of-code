package techpalace

import (
	"fmt"
	"strings"
)

// WelcomeMessage returns a welcome message for the customer.
func WelcomeMessage(customer string) string {
	const pre string = "Welcome to the Tech Palace, "

	return pre + strings.ToUpper(customer)
}

// AddBorder adds a border to a welcome message.
func AddBorder(welcomeMsg string, numStarsPerLine int) string {
	s := strings.Repeat("*", numStarsPerLine)

	return fmt.Sprintf("%s\n%s\n%s", s, welcomeMsg, s)
}

// CleanupMessage cleans up an old marketing message.
func CleanupMessage(oldMsg string) string {
	text := strings.ReplaceAll(oldMsg, "*", "")
	return strings.TrimSpace(text)
}
