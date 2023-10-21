package blackjack

// ParseCard returns the integer value of a card following blackjack ruleset.
func ParseCard(card string) int {
	switch card {
	case "two":
		return 2
	case "three":
		return 3
	case "four":
		return 4
	case "five":
		return 5
	case "six":
		return 6
	case "seven":
		return 7
	case "eight":
		return 8
	case "nine":
		return 9
	case "ten", "jack", "queen", "king":
		return 10
	case "ace":
		return 11
	default:
		return 0

	}
}

// FirstTurn returns the decision for the first turn, given two cards of the
// player and one card of the dealer.
func FirstTurn(card1, card2, dealerCard string) string {
	const (
		stand string = "S"
		hit   string = "H"
		split string = "P"
		win   string = "W"
	)

	c1, c2, dc := ParseCard(card1), ParseCard(card2), ParseCard(dealerCard)

	sum := c1 + c2

	switch {
	case sum >= 22:
		return split
	case sum == 21 && dc < 10:
		return win
	case sum == 21:
		return stand
	case sum <= 20 && sum >= 17:
		return stand
	case sum <= 16 && sum >= 12 && dc >= 7:
		return hit
	case sum <= 16 && sum >= 12:
		return stand
	case sum <= 11:
		return hit
	default:
		panic("we shouldn't be here")
	}
}
