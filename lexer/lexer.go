package lexer

import (
	"PLB-Interpreter/plbErrors"
	"PLB-Interpreter/tokens"
	"bufio"
	"errors"
	"io"
	"strings"
	"unicode"
)

type Lexer struct {
	input        *bufio.Reader
	position     int      // current position in input (points to current char)
	readPosition int      // current reading position in input (after current char)
	ch           byte     // current char under examination
	fileName     string   // filename of the file being lexed
	lineNumber   int      // currently processed line number
	col          int      // currently processed column/character/byte number
	lines        []string // lines of the file being lexed
	lineHadNonWS bool     // whether the current line had non-whitespace characters
	errors       []error  // errors encountered during lexing
}

// New Constructor for a new Lexer object, takes a *bufio.Reader and the filename as inputs,
// advances the lexer to the first character and returns a pointer to the new Lexer object.
func New(is *bufio.Reader, filename string) *Lexer {
	l := &Lexer{input: is, fileName: filename}

	// read all lines from the input stream
	for {
		line, err := l.input.ReadString('\n')
		if errors.Is(err, io.EOF) {
			// we encountered the EOF, so we add the last line and break
			l.lines = append(l.lines, line)
			break
		}
		if err != nil {
			// we encountered an error
			break
		}
		l.lines = append(l.lines, line)
	}

	// reset the reader to the beginning by adding in all the lines again
	l.input.Reset(strings.NewReader(strings.Join(l.lines, "")))

	// setup pointers
	l.readChar()

	// setup vars
	l.lineNumber = 1
	l.col = 1

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
func (l *Lexer) newToken(tokenType tokens.TokenType, ch byte) tokens.Token {
	return tokens.Token{
		Type:     tokenType,
		Literal:  string(ch),
		Line:     l.lineNumber,
		Col:      l.col,
		FileName: l.fileName,
		LineTxt:  l.lines[l.lineNumber-1],
	}
}

// NextToken returns the next token in the input stream.
// It advances the lexer to the next token and returns the token.
func (l *Lexer) NextToken() (tokens.Token, error) {
	var tok tokens.Token

	switch l.ch {
	case 0:
		tok = l.newToken(tokens.EOF, '\x00')
	case '.':
		if !l.lineHadNonWS {
			tok = l.handleComment()
		}
	case ' ', '\t':
		tok = tokens.Token{
			Line:     l.lineNumber,
			Col:      l.col,
			LineTxt:  l.lines[l.lineNumber-1],
			FileName: l.fileName,
		}
		tok.Type = tokens.WHITESPACE
		tok.Literal = l.consumeWhiteSpace()

	case '\n', '\r':
		if !l.lineHadNonWS {
			tok = l.newToken(tokens.NULLLINE, l.ch)
			l.consumeLine()
			l.lineNumber++
			l.col = 0
			l.lineHadNonWS = false
		} else {
			tok = l.newToken(tokens.NEWLINE, l.ch)
			l.lineNumber++
			l.col = 0
			l.lineHadNonWS = false
		}
	case '$':
		if !l.isLetter(l.peekChar()) {
			tok = l.newToken(tokens.CURRENCY, l.ch)
			l.lineHadNonWS = true
		}
	case '#':
		tok = l.newToken(tokens.FORCING, l.ch)
		l.lineHadNonWS = true
	case ',', ':':
		tok = l.newToken(tokens.COMMA, l.ch)
		l.lineHadNonWS = true
	case ';':
		tok = l.newToken(tokens.SEMICOLON, l.ch)
		l.lineHadNonWS = true
	case '(':
		tok = l.newToken(tokens.LPAREN, l.ch)
		l.lineHadNonWS = true
	case ')':
		tok = l.newToken(tokens.RPAREN, l.ch)
		l.lineHadNonWS = true
	case '*':
		if !l.lineHadNonWS {
			tok = l.handleComment()
		} else {
			l.lineHadNonWS = true
			if l.peekChar() == '*' {
				tok = tokens.Token{
					Line:     l.lineNumber,
					Col:      l.col,
					LineTxt:  l.lines[l.lineNumber-1],
					FileName: l.fileName,
				}
				tok.Type = tokens.POWER
				tok.Literal = "**"
				l.readChar()
				l.col++
			} else {
				tok = l.newToken(tokens.ASTERISK, l.ch)
			}
		}
	case '/':
		l.lineHadNonWS = true
		tok = l.newToken(tokens.SLASH, l.ch)
	case '+':
		if !l.lineHadNonWS {
			tok = l.handleComment()
		} else {
			l.lineHadNonWS = true
			tok = l.newToken(tokens.PLUS, l.ch)
		}
	case '-':
		l.lineHadNonWS = true
		tok = l.newToken(tokens.MINUS, l.ch)
	case '<':
		l.lineHadNonWS = true
		if l.peekChar() == '=' {
			tok = tokens.Token{
				Line:     l.lineNumber,
				Col:      l.col,
				LineTxt:  l.lines[l.lineNumber-1],
				FileName: l.fileName,
			}
			tok.Type = tokens.LEQ
			tok.Literal = "<="
			l.readChar()
			l.col++
		} else if l.peekChar() == '>' {
			tok = tokens.Token{
				Line:     l.lineNumber,
				Col:      l.col,
				LineTxt:  l.lines[l.lineNumber-1],
				FileName: l.fileName,
			}
			tok.Type = tokens.NEQ
			tok.Literal = "<>"
			l.readChar()
			l.col++
		} else {
			tok = l.newToken(tokens.LT, l.ch)
		}
	case '>':
		l.lineHadNonWS = true
		if l.peekChar() == '=' {
			tok = tokens.Token{
				Line:     l.lineNumber,
				Col:      l.col,
				LineTxt:  l.lines[l.lineNumber-1],
				FileName: l.fileName,
			}
			tok.Type = tokens.GEQ
			tok.Literal = ">="
			l.readChar()
			l.col++
		} else {
			tok = l.newToken(tokens.GT, l.ch)
		}
	case '=':
		l.lineHadNonWS = true
		tok = l.newToken(tokens.EQ, l.ch)
	case '"':
		if l.peekChar() == '-' || l.isDigit(l.peekChar()) || l.peekChar() == '.' {
			tok = tokens.Token{
				Line:     l.lineNumber,
				Col:      l.col,
				LineTxt:  l.lines[l.lineNumber-1],
				FileName: l.fileName,
			}
			tok.Type = tokens.NUMERICLITERAL
			tok.Literal = l.readLiteral()
			if l.containsNonNumerics(tok.Literal) {
				tok.Type = tokens.LITERAL
			}
			l.lineHadNonWS = true
		} else {
			tok = tokens.Token{
				Line:     l.lineNumber,
				Col:      l.col,
				LineTxt:  l.lines[l.lineNumber-1],
				FileName: l.fileName,
			}
			tok.Type = tokens.LITERAL
			tok.Literal = l.readLiteral()
			l.lineHadNonWS = true
		}
	default:
		l.lineHadNonWS = true
		if l.isHexDigit(l.ch) {
			tok = tokens.Token{
				Line:     l.lineNumber,
				Col:      l.col,
				LineTxt:  l.lines[l.lineNumber-1],
				FileName: l.fileName,
			}
			tok.Type = tokens.XNUM
			tok.Literal = l.readHex()
			return tok, nil
		} else if l.isOctDigit(l.ch) {
			tok = tokens.Token{
				Line:     l.lineNumber,
				Col:      l.col,
				LineTxt:  l.lines[l.lineNumber-1],
				FileName: l.fileName,
			}
			tok.Type = tokens.ONUM
			tok.Literal = l.readOct()
			return tok, nil
		} else if l.isDigit(l.ch) {
			tok = tokens.Token{
				Line:     l.lineNumber,
				Col:      l.col,
				LineTxt:  l.lines[l.lineNumber-1],
				FileName: l.fileName,
			}
			tok.Type = tokens.DNUM
			tok.Literal = l.readDec()
			return tok, nil
		} else if l.isLetter(l.ch) {
			tok = tokens.Token{
				Line:     l.lineNumber,
				Col:      l.col,
				LineTxt:  l.lines[l.lineNumber-1],
				FileName: l.fileName,
			}
			tok.Literal = l.readIdentifier()
			tok.Type = tokens.LookupIdent(tok.Literal)
			return tok, nil
		} else {
			tok = tokens.Token{
				Line:     l.lineNumber,
				Col:      l.col,
				LineTxt:  l.lines[l.lineNumber-1],
				FileName: l.fileName,
			}
			tok = l.newToken(tokens.ILLEGAL, l.ch)
			err := plbErrors.NewPLBError(
				"Lexer",
				"Invalid token type",
				l.fileName,
				l.lineNumber,
				l.col,
				strings.TrimRight(l.lines[l.lineNumber-1], "\r\n"),
			)
			l.errors = append(l.errors, err)
			return tok, err
		}
	}
	l.readChar()
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

	for l.isDigit(l.ch) {
		dec += string(l.ch)
		l.readChar()
	}

	return dec
}

// readLiteral consumes a string literal from the input stream and returns it as a string.
// It continues consuming until it finds a non Letter or non Digit character.
func (l *Lexer) readIdentifier() string {
	var lit string
	for l.isLetter(l.ch) || l.isDigit(l.ch) {
		lit += string(l.ch)
		l.readChar()
	}
	return lit
}

// isLetter returns true if the given byte is a letter. (a-z, A-Z, _, $)
func (l *Lexer) isLetter(ch byte) bool {
	if unicode.IsLetter(rune(ch)) || ch == '_' || ch == '$' {
		return true
	}
	return false
}

// consumeLine consumes the rest of the current lineNumber. The lexer pointers are advanced accordingly.
// Newline characters are not consumed
func (l *Lexer) consumeLine() {
	for l.ch != '\n' && l.ch != '\r' {
		l.readChar()
		if l.peekChar() == 0 {
			break
		}
	}
}

// readLiteral consumes a string literal from the input stream and returns it as a string.
// It continues consuming until it finds an unescaped closing quote. A double quote is escaped. ("" -> ")
// The opening and closing quotes are consumed but not returned.
func (l *Lexer) readLiteral() string {
	var lit string
	// Skip the opening quote
	l.position = l.readPosition
	l.readPosition++
	for {
		l.readChar()
		if l.ch == '"' {
			if l.peekChar() == '"' {
				l.readChar()
				// Skip the second quote
				l.position = l.readPosition
				l.readPosition++
			} else {
				// Skip the closing quote
				l.position = l.readPosition
				l.readPosition++
				break
			}
		}
		lit += string(l.ch)
	}

	return lit
}

// handleComment consumes the rest of the current lineNumber and returns a COMMENT token.
// The lexer pointers are advanced accordingly. The newline characters are not consumed.
func (l *Lexer) handleComment() tokens.Token {
	var tok tokens.Token
	tok.Type = tokens.COMMENT
	currLine := l.lines[l.lineNumber-1]
	if currLine[len(currLine)-1] == '\r' || currLine[len(currLine)-1] == '\n' {
		currLine = currLine[:len(currLine)-1]
	}
	tok.Literal = strings.TrimSpace(currLine)
	l.consumeLine()
	l.lineNumber++
	l.col = 0
	return tok
}

func (l *Lexer) containsNonNumerics(literal string) bool {
	for _, character := range literal {
		if character != '-' && character != '.' {
			if !unicode.IsNumber(character) {
				return true
			}
		}
	}
	return false
}

func (l *Lexer) consumeWhiteSpace() string {
	var whiteSpace string
	whiteSpace += string(l.ch)
	for l.ch == ' ' || l.ch == '\t' {
		if l.peekChar() != ' ' && l.peekChar() != '\t' {
			break
		}
		l.readChar()
		whiteSpace += string(l.ch)
	}
	return whiteSpace
}
