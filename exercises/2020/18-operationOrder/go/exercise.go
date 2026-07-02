package exercises

import (
	"fmt"
	"strings"

	"github.com/asphaltbuffet/advent-of-code/internal/common"
)

// Exercise for Advent of Code 2020 day 18.
type Exercise struct {
	common.BaseExercise
}

// tokenize splits an expression into number, operator, and parenthesis tokens.
func tokenize(line string) []string {
	var toks []string
	i := 0
	for i < len(line) {
		c := line[i]
		switch {
		case c == ' ':
			i++
		case c == '(' || c == ')' || c == '+' || c == '*':
			toks = append(toks, string(c))
			i++
		default: // number
			j := i
			for j < len(line) && line[j] >= '0' && line[j] <= '9' {
				j++
			}
			toks = append(toks, line[i:j])
			i = j
		}
	}
	return toks
}

// parser evaluates a token stream with precedence-climbing, using the provided
// operator precedence levels.
type parser struct {
	toks []string
	pos  int
	prec map[string]int
}

func (p *parser) peek() string {
	if p.pos < len(p.toks) {
		return p.toks[p.pos]
	}
	return ""
}

func (p *parser) next() string {
	t := p.toks[p.pos]
	p.pos++
	return t
}

// parseExpr evaluates an expression whose operators bind at least as tightly as
// minPrec (precedence climbing).
func (p *parser) parseExpr(minPrec int) int {
	left := p.parseAtom()
	for {
		op := p.peek()
		if op != "+" && op != "*" || p.prec[op] < minPrec {
			break
		}
		p.next()
		right := p.parseExpr(p.prec[op] + 1)
		if op == "+" {
			left += right
		} else {
			left *= right
		}
	}
	return left
}

// parseAtom evaluates a number or a parenthesized sub-expression.
func (p *parser) parseAtom() int {
	t := p.next()
	if t == "(" {
		v := p.parseExpr(1)
		p.next() // consume ')'
		return v
	}
	n := 0
	for _, ch := range t {
		n = n*10 + int(ch-'0')
	}
	return n
}

// evalAll sums the evaluation of every line using the given precedence map.
func evalAll(instr string, prec map[string]int) int {
	total := 0
	for _, line := range strings.Split(strings.TrimSpace(instr), "\n") {
		p := &parser{toks: tokenize(line), prec: prec}
		total += p.parseExpr(1)
	}
	return total
}

// One evaluates with + and * at equal precedence (left to right).
func (e Exercise) One(instr string) (any, error) {
	return fmt.Sprintf("%d", evalAll(instr, map[string]int{"+": 1, "*": 1})), nil
}

// Two evaluates with + binding tighter than * (addition before multiplication).
func (e Exercise) Two(instr string) (any, error) {
	return fmt.Sprintf("%d", evalAll(instr, map[string]int{"+": 2, "*": 1})), nil
}
