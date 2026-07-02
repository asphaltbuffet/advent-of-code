package exercises

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Vis draws the key exchange as a diagram. Both the card and the door transform
// the subject number 7 by their secret loop size to get a public key; recovering
// one loop size (the discrete log that Part One brute-forces) and applying it to
// the other public key yields the same shared encryption key from either side.
// The two sides use distinct colorblind-safe colors, and the shared result is
// highlighted; labels carry all the numbers, so the diagram reads in grayscale.
func (e Exercise) Vis(instr, outdir string) error {
	card, door, err := parse(instr)
	if err != nil {
		return err
	}
	cardLoop := findLoop(card)
	doorLoop := findLoop(door)
	secret := modPow(door, cardLoop)

	const (
		W  = 900
		H  = 420
		mL = 40
	)
	cardCol := "#E69F00" // orange
	doorCol := "#0072B2" // blue
	secretCol := "#F0E442"

	var sb strings.Builder
	fmt.Fprintf(&sb, `<svg xmlns="http://www.w3.org/2000/svg" width="%d" height="%d" viewBox="0 0 %d %d" font-family="monospace">`, W, H, W, H)
	fmt.Fprintf(&sb, `<rect width="%d" height="%d" fill="#111418"/>`, W, H)
	fmt.Fprintf(&sb, `<text x="%d" y="34" fill="#e8ecf4" font-size="16">Combo Breaker: the key exchange (mod %d)</text>`, mL, modulus)

	box := func(x, y, w, h int, col, title, line1, line2 string) {
		fmt.Fprintf(&sb, `<rect x="%d" y="%d" width="%d" height="%d" rx="6" fill="none" stroke="%s" stroke-width="2"/>`, x, y, w, h, col)
		fmt.Fprintf(&sb, `<text x="%d" y="%d" fill="%s" font-size="14" font-weight="bold">%s</text>`, x+12, y+22, col, title)
		fmt.Fprintf(&sb, `<text x="%d" y="%d" fill="#e8ecf4" font-size="12">%s</text>`, x+12, y+42, line1)
		if line2 != "" {
			fmt.Fprintf(&sb, `<text x="%d" y="%d" fill="#c8d0dc" font-size="12">%s</text>`, x+12, y+60, line2)
		}
	}

	bw, bh := 340, 74
	lx, rx := mL, W-mL-bw

	// Row 1: each side's secret loop and public key.
	box(lx, 70, bw, bh, cardCol, "Card",
		fmt.Sprintf("secret loop: %d", cardLoop),
		fmt.Sprintf("public = 7^loop mod m = %d", card))
	box(rx, 70, bw, bh, doorCol, "Door",
		fmt.Sprintf("secret loop: %d", doorLoop),
		fmt.Sprintf("public = 7^loop mod m = %d", door))

	// Cross arrows: each side raises the OTHER's public key to its own loop.
	fmt.Fprintf(&sb, `<line x1="%d" y1="%d" x2="%d" y2="%d" stroke="%s" stroke-width="1.5" marker-end="url(#a1)"/>`, lx+bw, 130, rx, 200, cardCol)
	fmt.Fprintf(&sb, `<line x1="%d" y1="%d" x2="%d" y2="%d" stroke="%s" stroke-width="1.5" marker-end="url(#a2)"/>`, rx, 130, lx+bw, 200, doorCol)
	fmt.Fprintf(&sb, `<defs><marker id="a1" markerWidth="8" markerHeight="8" refX="6" refY="3" orient="auto"><path d="M0,0 L6,3 L0,6 Z" fill="%s"/></marker>`, cardCol)
	fmt.Fprintf(&sb, `<marker id="a2" markerWidth="8" markerHeight="8" refX="6" refY="3" orient="auto"><path d="M0,0 L6,3 L0,6 Z" fill="%s"/></marker></defs>`, doorCol)

	// Row 2: the two computations that both land on the secret.
	box(lx, 210, bw, bh, cardCol, "Card computes",
		"door_public ^ card_loop mod m",
		fmt.Sprintf("= %d ^ %d = %d", door, cardLoop, secret))
	box(rx, 210, bw, bh, doorCol, "Door computes",
		"card_public ^ door_loop mod m",
		fmt.Sprintf("= %d ^ %d = %d", card, doorLoop, secret))

	// Shared secret.
	sy := 330
	fmt.Fprintf(&sb, `<rect x="%d" y="%d" width="%d" height="50" rx="8" fill="%s"/>`, W/2-180, sy, 360, secretCol)
	fmt.Fprintf(&sb, `<text x="%d" y="%d" fill="#0d0f18" font-size="16" font-weight="bold" text-anchor="middle">shared encryption key = %d</text>`, W/2, sy+31, secret)

	sb.WriteString(`</svg>`)
	return os.WriteFile(filepath.Join(outdir, "combo-breaker.svg"), []byte(sb.String()), 0o600)
}
