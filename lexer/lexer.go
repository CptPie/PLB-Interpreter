package lexer

import (
	"PLB/token"
	"strings"
)

type Lexer struct {
	input        string
	position     int  // current char position
	readPosition int  // next char position
	ch           byte // current char
}

func New(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition += 1
}

func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	} else {
		return l.input[l.readPosition]
	}
}
func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	switch {
	// Misc Symbols
	case l.ch == '$':
		tok = newToken(token.CURRENCY, l.ch)
	case l.ch == '#':
		tok = newToken(token.FORCING, l.ch)

	// Arithmetic operators
	case l.ch == '+':
		tok = newToken(token.PLUS, l.ch)
	case l.ch == '-' && l.peekChar() == ' ':
		tok = newToken(token.MINUS, l.ch)
	case l.ch == '/':
		tok = newToken(token.DIVIDE, l.ch)
	case l.ch == '*':
		if l.peekChar() == '*' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{
				Type:    token.EXPONENT,
				Literal: literal,
			}
		} else {
			tok = newToken(token.MULTIPLY, l.ch)
		}

	// Relational Operators
	case l.ch == '=':
		tok = newToken(token.EQ, l.ch)
	case l.ch == '<':
		if l.peekChar() == '>' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{
				Type:    token.NEQ,
				Literal: literal,
			}
		} else if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{
				Type:    token.LTEQ,
				Literal: literal,
			}
		} else {
			tok = newToken(token.LT, l.ch)
		}
	case l.ch == '>':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{
				Type:    token.GTEQ,
				Literal: literal,
			}
		} else {
			tok = newToken(token.GT, l.ch)
		}

	case l.ch == ';':
		tok = newToken(token.SEMICOLON, l.ch)
	case l.ch == '(':
		tok = newToken(token.LPAREN, l.ch)
	case l.ch == ')':
		tok = newToken(token.RPAREN, l.ch)
	case l.ch == ',' || l.ch == ':':
		tok = newToken(token.COMMA, l.ch)
	case l.ch == '\n' || l.ch == '\r':
		tok = newToken(token.NEWLINE, l.ch)
	case l.ch == ' ' || l.ch == '\t':
		tok = newToken(token.WHITESPACE, l.ch)
	case l.ch == 0:
		tok.Literal = ""
		tok.Type = token.EOF
	default:
		if isLetter(l.ch) {
			tok.Literal = l.readIdentifier()
			tok.Type = token.LookupIdent(tok.Literal)
			return tok
		} else if l.ch == '0' && isOctal(l.peekChar()) {
			tok.Type = token.OCT
			tok.Literal = l.readOctal()
			return tok
		} else if l.ch == '0' && l.peekChar() == 'x' && isHex(l.input[l.readPosition+1]) {
			tok.Type = token.HEX
			tok.Literal = l.readHex()
			return tok
		} else if (l.ch == '-' && isFloat(l.peekChar())) || isFloat(l.ch) {
			tok.Literal = l.readFloat()
			if strings.Contains(tok.Literal, ".") {
				tok.Type = token.FLOAT
				return tok
			} else {
				tok.Type = token.INT
				return tok
			}
		} else if (l.ch == '-' && isDigit(l.peekChar())) || isDigit(l.ch) {
			tok.Type = token.INT
			tok.Literal = l.readNumber()
			return tok
		} else {
			tok = newToken(token.ILLEGAL, l.ch)
		}
	}
	l.readChar()
	return tok
}

func newToken(tokenType token.TokenType, ch byte) token.Token {
	return token.Token{Type: tokenType, Literal: string(ch)}
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func (l *Lexer) readIdentifier() string {
	position := l.position
	for isLetter(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position] // Read until we encounter a non "isLetter" letter
}

func (l *Lexer) readNumber() string {
	position := l.position
	if l.ch == '-' {
		// if there is a leading minus, advance the readPointer to the next
		// position but keep var position at the minus
		l.readChar()
	}
	for isDigit(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func (l *Lexer) readOctal() string {
	position := l.position
	for isOctal(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

func isOctal(ch byte) bool {
	return '0' <= ch && ch <= '7'
}
func (l *Lexer) readHex() string {
	position := l.position
	for isHex(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

func isHex(ch byte) bool {
	return '0' <= ch && ch <= '9' || 'a' <= ch && ch <= 'f' || 'A' <= ch && ch <= 'F' || ch == 'x' || ch == 'X'
}

func (l *Lexer) readFloat() string {
	position := l.position
	if l.ch == '-' {
		// if there is a leading minus, advance the readPointer to the next
		// position but keep var position at the minus
		l.readChar()
	}
	for isFloat(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}
func isFloat(ch byte) bool {
	return '0' <= ch && ch <= '9' || ch == '.'
}
