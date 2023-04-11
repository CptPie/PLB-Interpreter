package lexer

import (
	"PLB-Interpreter/plbErrors"
	"PLB-Interpreter/tokens"
	"bufio"
	"bytes"
	"errors"
	"io"
)

type Lexer struct {
	input        *bufio.Reader
	position     int      // current position in input (points to current char)
	readPosition int      // current reading position in input (after current char)
	ch           byte     // current char under examination
	fileName     string   // filename of the file being lexed
	line         int      // currently processed line number
	col          int      // currently processed column/character/byte number
	lines        []string // lines of the file being lexed
	errors       []error  // errors encountered during lexing
}

// New Constructor for a new Lexer object, takes a bufio.Reader and the filename as inputs,
// advances the lexer to the first character and returns a pointer to the new Lexer object.
func New(is *bufio.Reader, filename string) *Lexer {
	l := &Lexer{input: is, fileName: filename}
	// setup pointers
	l.readChar()

	// setup vars
	l.line = 1
	l.col = 1

	// some buffer fuckery to get the lines out of the reader WITHOUT advancing the pointers of the reader
	var contents []byte
	contents = append(contents, l.ch)
	remainder, _ := is.Peek(is.Size())
	contents = append(contents, remainder...)
	lines := bufio.NewReader(bytes.NewReader(contents))
	for {
		line, err := lines.ReadString('\n')
		if errors.Is(err, io.EOF) {
			l.lines = append(l.lines, line)
			break
		}
		if err != nil {
			break
		}
		l.lines = append(l.lines, line)
	}

	return l
}

// readChar() reads the next character in the input string and advances the position of the lexer.
func (l *Lexer) readChar() {
	// save next byte in ch
	l.ch, _ = l.input.ReadByte()
	// update position
	l.position = l.readPosition
	// advance readPosition by 1
	l.readPosition += 1
	l.col++
}

// newToken returns a new token with the given type and literal.
func newToken(tokenType tokens.TokenType, ch byte) tokens.Token {
	return tokens.Token{Type: tokenType, Literal: string(ch)}
}

// NextToken returns the next token in the input stream.
// It advances the lexer to the next token and returns the token.
func (l *Lexer) NextToken() (tokens.Token, error) {
	var tok tokens.Token

	char := rune(l.ch)
	switch char {
	case ' ', '\t':
		tok = newToken(tokens.WHITESPACE, l.ch)
	case '\n', '\r':
		tok = newToken(tokens.NEWLINE, l.ch)
		l.line++
		l.col = 0
	case '-':
		if l.isDigit(l.peekChar()) {
			tok.Type = tokens.SIGNEDDNUM
			tok.Literal = l.readDec()
			return tok, nil
		}
	default:
		if l.isHexDigit(l.ch) {
			tok.Type = tokens.XNUM
			tok.Literal = l.readHex()
			return tok, nil
		} else if l.isOctDigit(l.ch) {
			tok.Type = tokens.ONUM
			tok.Literal = l.readOct()
			return tok, nil
		} else if l.isDigit(l.ch) {
			tok.Type = tokens.DNUM
			tok.Literal = l.readDec()
			return tok, nil
		} else {
			tok = newToken(tokens.ILLEGAL, l.ch)
			err := plbErrors.NewPLBError("Lexer", "Invalid token type", l.fileName, l.line, l.col, l.lines[l.line-1])
			l.errors = append(l.errors, err)
			return tok, err
		}
	}
	l.readChar()
	l.col++
	return tok, nil
}

// isDigit returns true if the given byte is a digit. (0-9)
func (l *Lexer) isDigit(ch byte) bool {
	if ch >= '0' && ch <= '9' {
		return true
	}
	return false
}

// isHexDigit returns true if the given byte is a hex digit. (0-9, a-f, A-F)
func (l *Lexer) isHexDigit(ch byte) bool {
	if ch == '0' && (l.peekChar() == 'x' || l.peekChar() == 'X') {
		return true
	}
	return false
}

// isOctDigit returns true if the given byte is an octal digit. (0-7)
func (l *Lexer) isOctDigit(ch byte) bool {
	peek := l.peekChar()
	if ch == '0' && (l.isDigit(peek) && peek != '8' && peek != '9') {
		return true
	}
	return false
}

// peekChar returns the next character in the input stream without advancing the lexer.
func (l *Lexer) peekChar() byte {
	ch, err := l.input.Peek(1)
	if err != nil {
		return 0
	}
	return ch[0]
}

// readHex reads a hex number from the input stream and returns it as a string.
func (l *Lexer) readHex() string {
	var hex string
	hex += string(l.ch)
	l.readChar()
	hex += string(l.ch)
	l.readChar()
	for l.isDigit(l.ch) || (l.ch >= 'a' && l.ch <= 'f') || (l.ch >= 'A' && l.ch <= 'F') {
		hex += string(l.ch)
		l.readChar()
	}
	return hex
}

// readOct reads an octal number from the input stream and returns it as a string.
func (l *Lexer) readOct() string {
	var oct string
	for l.isDigit(l.ch) && l.ch != '8' && l.ch != '9' {
		oct += string(l.ch)
		l.readChar()
	}
	return oct
}

// readDec reads a decimal number from the input stream and returns it as a string.
func (l *Lexer) readDec() string {
	var dec string
	// account for signed dnums
	if l.ch == '-' {
		dec += string(l.ch)
		l.readChar()
	}

	for l.isDigit(l.ch) {
		dec += string(l.ch)
		l.readChar()
	}

	return dec
}
