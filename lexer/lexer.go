package lexer

import (
	"tgoc/token"
)

// Lexer type
type Lexer struct {
	Tokens []token.Token
	Pos    int
	Input  string
}

// New lexer create
func New(input string) *Lexer {
	t := []token.Token{}
	return &Lexer{Tokens: t, Input: input, Pos: 0}
}

// Analyze the input string ans split the token sequences
func (l *Lexer) Analyze() {
	for ; l.Pos < len(l.Input); l.Pos++ {
		l.skip()

		var tok token.Token

		switch l.Input[l.Pos] {
		case '+':
			tok = token.New(token.ADD, "+")
		case '-':
			tok = token.New(token.SUB, "-")
		case '*':
			tok = token.New(token.MUL, "*")
		case '/':
			tok = token.New(token.DIV, "/")
		case '%':
			tok = token.New(token.REM, "%")
		case '(':
			tok = token.New(token.LPAREN, "(")
		case ')':
			tok = token.New(token.RPAREN, ")")
		default:
			if isDigit(l.Input[l.Pos]) {
				tok = token.New(token.INT, l.readDigit())
			}
		}

		l.Tokens = append(l.Tokens, tok)
	}
}

func isDigit(c byte) bool {
	return ('0' <= c) && (c <= '9')
}

func isSpace(c byte) bool {
	return c == ' ' || c == '\t' || c == '\n'
}

func (l *Lexer) skip() {
	for isSpace(l.Input[l.Pos]) {
		l.Pos++
	}
}

func (l *Lexer) readDigit() string {
	head := l.Pos
	tail := l.Pos + 1
	for ; tail < len(l.Input); tail++ {
		if !isDigit(l.Input[tail]) {
			break
		}
	}

	l.Pos = tail - 1
	return l.Input[head:tail]
}
