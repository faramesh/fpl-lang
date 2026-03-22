package lexer

import (
	"strings"
	"unicode"

	"github.com/faramesh/fpl-lang/reference/go/fpl/token"
)

type Lexer struct {
	input  []rune
	pos    int
	line   int
	column int
}

func New(input string) *Lexer {
	return &Lexer{input: []rune(input), line: 1, column: 1}
}

func (l *Lexer) NextToken() token.Token {
	l.skipWhitespaceNoNewline()

	if l.pos >= len(l.input) {
		return token.Token{Type: token.EOF, Line: l.line, Column: l.column}
	}

	ch := l.input[l.pos]
	line, col := l.line, l.column

	switch ch {
	case '{':
		l.advance()
		return token.Token{Type: token.LBrace, Literal: "{", Line: line, Column: col}
	case '}':
		l.advance()
		return token.Token{Type: token.RBrace, Literal: "}", Line: line, Column: col}
	case '(':
		l.advance()
		return token.Token{Type: token.LParen, Literal: "(", Line: line, Column: col}
	case ')':
		l.advance()
		return token.Token{Type: token.RParen, Literal: ")", Line: line, Column: col}
	case '[':
		l.advance()
		return token.Token{Type: token.LBracket, Literal: "[", Line: line, Column: col}
	case ']':
		l.advance()
		return token.Token{Type: token.RBracket, Literal: "]", Line: line, Column: col}
	case ',':
		l.advance()
		return token.Token{Type: token.Comma, Literal: ",", Line: line, Column: col}
	case '*':
		l.advance()
		return token.Token{Type: token.Star, Literal: "*", Line: line, Column: col}
	case '\n':
		l.advance()
		return token.Token{Type: token.Newline, Literal: "\\n", Line: line, Column: col}
	case '#':
		l.skipComment()
		return l.NextToken()
	case '"':
		lit := l.readString()
		return token.Token{Type: token.String, Literal: lit, Line: line, Column: col}
	case ':':
		l.advance()
		return token.Token{Type: token.Colon, Literal: ":", Line: line, Column: col}
	case '>':
		if l.peek() == '=' {
			l.advance()
			l.advance()
			return token.Token{Type: token.GE, Literal: ">=", Line: line, Column: col}
		}
		l.advance()
		return token.Token{Type: token.GT, Literal: ">", Line: line, Column: col}
	case '<':
		if l.peek() == '=' {
			l.advance()
			l.advance()
			return token.Token{Type: token.LE, Literal: "<=", Line: line, Column: col}
		}
		l.advance()
		return token.Token{Type: token.LT, Literal: "<", Line: line, Column: col}
	case '=':
		if l.peek() == '=' {
			l.advance()
			l.advance()
			return token.Token{Type: token.EQ, Literal: "==", Line: line, Column: col}
		}
		l.advance()
		return token.Token{Type: token.Illegal, Literal: "=", Line: line, Column: col}
	case '!':
		if l.peek() == '=' {
			l.advance()
			l.advance()
			return token.Token{Type: token.NE, Literal: "!=", Line: line, Column: col}
		}
		l.advance()
		return token.Token{Type: token.Illegal, Literal: "!", Line: line, Column: col}
	case '$':
		lit := l.readCurrency()
		return token.Token{Type: token.Currency, Literal: lit, Line: line, Column: col}
	default:
		if unicode.IsDigit(ch) {
			lit := l.readNumber()
			return token.Token{Type: token.Number, Literal: lit, Line: line, Column: col}
		}
		if isIdentStart(ch) {
			ident := l.readIdent()
			if ident == "deny!" {
				return token.Token{Type: token.DenyBang, Literal: ident, Line: line, Column: col}
			}
			return token.Token{Type: token.LookupIdent(ident), Literal: ident, Line: line, Column: col}
		}
		l.advance()
		return token.Token{Type: token.Illegal, Literal: string(ch), Line: line, Column: col}
	}
}

func (l *Lexer) skipWhitespaceNoNewline() {
	for l.pos < len(l.input) {
		if l.input[l.pos] == '\n' {
			return
		}
		if !unicode.IsSpace(l.input[l.pos]) {
			return
		}
		l.advance()
	}
}

func (l *Lexer) skipComment() {
	for l.pos < len(l.input) && l.input[l.pos] != '\n' {
		l.advance()
	}
}

func (l *Lexer) readString() string {
	l.advance()
	start := l.pos
	for l.pos < len(l.input) && l.input[l.pos] != '"' {
		l.advance()
	}
	lit := string(l.input[start:l.pos])
	if l.pos < len(l.input) {
		l.advance()
	}
	return lit
}

func (l *Lexer) readNumber() string {
	start := l.pos
	hasDot := false
	for l.pos < len(l.input) {
		ch := l.input[l.pos]
		if unicode.IsDigit(ch) {
			l.advance()
			continue
		}
		if ch == '.' && !hasDot {
			hasDot = true
			l.advance()
			continue
		}
		break
	}
	return string(l.input[start:l.pos])
}

func (l *Lexer) readCurrency() string {
	start := l.pos
	l.advance()
	for l.pos < len(l.input) {
		ch := l.input[l.pos]
		if unicode.IsDigit(ch) || ch == '.' {
			l.advance()
			continue
		}
		break
	}
	return string(l.input[start:l.pos])
}

func (l *Lexer) readIdent() string {
	start := l.pos
	for l.pos < len(l.input) && isIdentPart(l.input[l.pos]) {
		l.advance()
	}
	return strings.TrimSpace(string(l.input[start:l.pos]))
}

func isIdentStart(ch rune) bool {
	return unicode.IsLetter(ch)
}

func isIdentPart(ch rune) bool {
	return unicode.IsLetter(ch) || unicode.IsDigit(ch) || ch == '_' || ch == '-' || ch == '!' || ch == '/' || ch == '.'
}

func (l *Lexer) peek() rune {
	if l.pos+1 >= len(l.input) {
		return 0
	}
	return l.input[l.pos+1]
}

func (l *Lexer) advance() {
	if l.pos >= len(l.input) {
		return
	}
	if l.input[l.pos] == '\n' {
		l.line++
		l.column = 1
	} else {
		l.column++
	}
	l.pos++
}
