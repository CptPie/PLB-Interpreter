package lexer

import (
	"PLB-Interpreter/plbErrors"
	"PLB-Interpreter/tokens"
	"bufio"
	"strings"
	"testing"
)

func TestLexer_NextToken_DOXNUM(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  []tokens.Token
	}{
		{
			name:  "Spaces and tabs",
			input: "2   \t\t\t",
			want: []tokens.Token{
				{Type: tokens.DNUM, Literal: "2"},
				{Type: tokens.WHITESPACE, Literal: "   \t\t\t"},
			},
		},
		{
			name:  "Newlines and carriage returns",
			input: "2\n2\r2\n2\r",
			want: []tokens.Token{
				{Type: tokens.DNUM, Literal: "2"},
				{Type: tokens.NEWLINE, Literal: "\n"},
				{Type: tokens.DNUM, Literal: "2"},
				{Type: tokens.NEWLINE, Literal: "\r"},
				{Type: tokens.DNUM, Literal: "2"},
				{Type: tokens.NEWLINE, Literal: "\n"},
				{Type: tokens.DNUM, Literal: "2"},
				{Type: tokens.NEWLINE, Literal: "\r"},
			},
		},
		{
			name:  "Hexadecimal numbers",
			input: "0x378\n0X38",
			want: []tokens.Token{
				{Type: tokens.XNUM, Literal: "0x378"},
				{Type: tokens.NEWLINE, Literal: "\n"},
				{Type: tokens.XNUM, Literal: "0X38"},
			},
		},
		{
			name:  "Octal numbers",
			input: "0377\n0377",
			want: []tokens.Token{
				{Type: tokens.ONUM, Literal: "0377"},
				{Type: tokens.NEWLINE, Literal: "\n"},
				{Type: tokens.ONUM, Literal: "0377"},
			},
		},
		{
			name:  "Decimal numbers",
			input: "123\n456",
			want: []tokens.Token{
				{Type: tokens.DNUM, Literal: "123"},
				{Type: tokens.NEWLINE, Literal: "\n"},
				{Type: tokens.DNUM, Literal: "456"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := bufio.NewReader(strings.NewReader(tt.input))
			l := New(reader, "")
			for i, want := range tt.want {
				got, _ := l.NextToken()
				if got != want {
					t.Errorf("Test %d: got %q, want %q", i, got, want)
				}
			}
		})
	}
}

func TestLexer_NextToken_Error(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  error
	}{
		{
			name:  "Invalid token 1",
			input: "!hello",
			want:  &plbErrors.PLBError{ErrorCode: "Lexer", Message: "Invalid token type", File: "test", LineNumber: 1, Column: 1, LineText: "!hello"},
		},
		{
			name:  "Invalid token 2",
			input: "1383!",
			want:  &plbErrors.PLBError{ErrorCode: "Lexer", Message: "Invalid token type", File: "test", LineNumber: 1, Column: 5, LineText: "1383!"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := bufio.NewReader(strings.NewReader(tt.input))
			l := New(reader, "test")
			cont := true
			for cont {
				_, err := l.NextToken()
				if err != nil {
					if err.Error() != tt.want.Error() {
						t.Errorf("got %q, want %q", err, tt.want)
					}
					cont = false
				}
			}

		})
	}
}

func TestLexer_NextToken_Symbols(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  []tokens.Token
	}{
		{
			name:  "prefix currency symbol",
			input: "$1234",
			want: []tokens.Token{
				{Type: tokens.CURRENCY, Literal: "$"},
				{Type: tokens.DNUM, Literal: "1234"},
			},
		},
		{
			name:  "postfix currency symbol",
			input: "1234$",
			want: []tokens.Token{
				{Type: tokens.DNUM, Literal: "1234"},
				{Type: tokens.CURRENCY, Literal: "$"},
			},
		},
		{
			name:  "spaces around currency symbol",
			input: "1234 $ ",
			want: []tokens.Token{
				{Type: tokens.DNUM, Literal: "1234"},
				{Type: tokens.WHITESPACE, Literal: " "},
				{Type: tokens.CURRENCY, Literal: "$"},
				{Type: tokens.WHITESPACE, Literal: " "},
			},
		},
		{
			name:  "forcing character",
			input: "1234#",
			want: []tokens.Token{
				{Type: tokens.DNUM, Literal: "1234"},
				{Type: tokens.FORCING, Literal: "#"},
			},
		},
		{
			name:  "comma separating numbers",
			input: "1234,5678",
			want: []tokens.Token{
				{Type: tokens.DNUM, Literal: "1234"},
				{Type: tokens.COMMA, Literal: ","},
				{Type: tokens.DNUM, Literal: "5678"},
			},
		},
		{
			name:  "comma separating numbers with spaces",
			input: "1234 , 5678",
			want: []tokens.Token{
				{Type: tokens.DNUM, Literal: "1234"},
				{Type: tokens.WHITESPACE, Literal: " "},
				{Type: tokens.COMMA, Literal: ","},
				{Type: tokens.WHITESPACE, Literal: " "},
				{Type: tokens.DNUM, Literal: "5678"},
			},
		},
		{
			name:  "comma separating numbers with spaces and currency symbol",
			input: "1234 , $5678",
			want: []tokens.Token{
				{Type: tokens.DNUM, Literal: "1234"},
				{Type: tokens.WHITESPACE, Literal: " "},
				{Type: tokens.COMMA, Literal: ","},
				{Type: tokens.WHITESPACE, Literal: " "},
				{Type: tokens.CURRENCY, Literal: "$"},
				{Type: tokens.DNUM, Literal: "5678"},
			},
		},
		{
			name:  "semicolon separating numbers",
			input: "1234;5678",
			want: []tokens.Token{
				{Type: tokens.DNUM, Literal: "1234"},
				{Type: tokens.SEMICOLON, Literal: ";"},
				{Type: tokens.DNUM, Literal: "5678"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := bufio.NewReader(strings.NewReader(tt.input))
			l := New(reader, "")
			for i, want := range tt.want {
				got, _ := l.NextToken()
				if got != want {
					t.Errorf("Test %d: got %q, want %q", i, got, want)
				}
			}
		})
	}
}

func TestLexer_NextToken_Operators(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  []tokens.Token
	}{
		{
			name:  "addition",
			input: "1+1",
			want: []tokens.Token{
				{Type: tokens.DNUM, Literal: "1"},
				{Type: tokens.PLUS, Literal: "+"},
				{Type: tokens.DNUM, Literal: "1"},
			},
		},
		{
			name:  "addition spaces",
			input: "1 + 1",
			want: []tokens.Token{
				{Type: tokens.DNUM, Literal: "1"},
				{Type: tokens.WHITESPACE, Literal: " "},
				{Type: tokens.PLUS, Literal: "+"},
				{Type: tokens.WHITESPACE, Literal: " "},
				{Type: tokens.DNUM, Literal: "1"},
			},
		},
		{
			name:  "subtraction",
			input: "1-1",
			want: []tokens.Token{
				{Type: tokens.DNUM, Literal: "1"},
				{Type: tokens.MINUS, Literal: "-"},
				{Type: tokens.DNUM, Literal: "1"},
			},
		},
		{
			name:  "subtraction spaces",
			input: "1 - 1",
			want: []tokens.Token{
				{Type: tokens.DNUM, Literal: "1"},
				{Type: tokens.WHITESPACE, Literal: " "},
				{Type: tokens.MINUS, Literal: "-"},
				{Type: tokens.WHITESPACE, Literal: " "},
				{Type: tokens.DNUM, Literal: "1"},
			},
		},
		{
			name:  "multiplication",
			input: "1*1",
			want: []tokens.Token{
				{Type: tokens.DNUM, Literal: "1"},
				{Type: tokens.ASTERISK, Literal: "*"},
				{Type: tokens.DNUM, Literal: "1"},
			},
		},
		{
			name:  "multiplication spaces",
			input: "1 * 1",
			want: []tokens.Token{
				{Type: tokens.DNUM, Literal: "1"},
				{Type: tokens.WHITESPACE, Literal: " "},
				{Type: tokens.ASTERISK, Literal: "*"},
				{Type: tokens.WHITESPACE, Literal: " "},
				{Type: tokens.DNUM, Literal: "1"},
			},
		},
		{
			name:  "division",
			input: "1/1",
			want: []tokens.Token{
				{Type: tokens.DNUM, Literal: "1"},
				{Type: tokens.SLASH, Literal: "/"},
				{Type: tokens.DNUM, Literal: "1"},
			},
		},
		{
			name:  "division spaces",
			input: "1 / 1",
			want: []tokens.Token{
				{Type: tokens.DNUM, Literal: "1"},
				{Type: tokens.WHITESPACE, Literal: " "},
				{Type: tokens.SLASH, Literal: "/"},
				{Type: tokens.WHITESPACE, Literal: " "},
				{Type: tokens.DNUM, Literal: "1"},
			},
		},
		{
			name:  "exponent",
			input: "1**1",
			want: []tokens.Token{
				{Type: tokens.DNUM, Literal: "1"},
				{Type: tokens.POWER, Literal: "**"},
				{Type: tokens.DNUM, Literal: "1"},
			},
		},
		{
			name:  "exponent spaces",
			input: "1 ** 1",
			want: []tokens.Token{
				{Type: tokens.DNUM, Literal: "1"},
				{Type: tokens.WHITESPACE, Literal: " "},
				{Type: tokens.POWER, Literal: "**"},
				{Type: tokens.WHITESPACE, Literal: " "},
				{Type: tokens.DNUM, Literal: "1"},
			},
		},
		{
			name:  "equality",
			input: `1=1`,
			want: []tokens.Token{
				{Type: tokens.DNUM, Literal: "1"},
				{Type: tokens.EQ, Literal: "="},
				{Type: tokens.DNUM, Literal: "1"},
			},
		},
		{
			name:  "equality spaces",
			input: `1 = 1`,
			want: []tokens.Token{
				{Type: tokens.DNUM, Literal: "1"},
				{Type: tokens.WHITESPACE, Literal: " "},
				{Type: tokens.EQ, Literal: "="},
				{Type: tokens.WHITESPACE, Literal: " "},
				{Type: tokens.DNUM, Literal: "1"},
			},
		},
		{
			name:  "inequality",
			input: `1<>1`,
			want: []tokens.Token{
				{Type: tokens.DNUM, Literal: "1"},
				{Type: tokens.NEQ, Literal: "<>"},
				{Type: tokens.DNUM, Literal: "1"},
			},
		},
		{
			name:  "inequality spaces",
			input: `1 <> 1`,
			want: []tokens.Token{
				{Type: tokens.DNUM, Literal: "1"},
				{Type: tokens.WHITESPACE, Literal: " "},
				{Type: tokens.NEQ, Literal: "<>"},
				{Type: tokens.WHITESPACE, Literal: " "},
				{Type: tokens.DNUM, Literal: "1"},
			},
		},
		{
			name:  "less than",
			input: `1<1`,
			want: []tokens.Token{
				{Type: tokens.DNUM, Literal: "1"},
				{Type: tokens.LT, Literal: "<"},
				{Type: tokens.DNUM, Literal: "1"},
			},
		},
		{
			name:  "less than spaces",
			input: `1 < 1`,
			want: []tokens.Token{
				{Type: tokens.DNUM, Literal: "1"},
				{Type: tokens.WHITESPACE, Literal: " "},
				{Type: tokens.LT, Literal: "<"},
				{Type: tokens.WHITESPACE, Literal: " "},
				{Type: tokens.DNUM, Literal: "1"},
			},
		},
		{
			name:  "less than or equal",
			input: `1<=1`,
			want: []tokens.Token{
				{Type: tokens.DNUM, Literal: "1"},
				{Type: tokens.LEQ, Literal: "<="},
				{Type: tokens.DNUM, Literal: "1"},
			},
		},
		{
			name:  "less than or equal spaces",
			input: `1 <= 1`,
			want: []tokens.Token{
				{Type: tokens.DNUM, Literal: "1"},
				{Type: tokens.WHITESPACE, Literal: " "},
				{Type: tokens.LEQ, Literal: "<="},
				{Type: tokens.WHITESPACE, Literal: " "},
				{Type: tokens.DNUM, Literal: "1"},
			},
		},
		{
			name:  "greater than",
			input: `1>1`,
			want: []tokens.Token{
				{Type: tokens.DNUM, Literal: "1"},
				{Type: tokens.GT, Literal: ">"},
				{Type: tokens.DNUM, Literal: "1"},
			},
		},
		{
			name:  "greater than spaces",
			input: `1 > 1`,
			want: []tokens.Token{
				{Type: tokens.DNUM, Literal: "1"},
				{Type: tokens.WHITESPACE, Literal: " "},
				{Type: tokens.GT, Literal: ">"},
				{Type: tokens.WHITESPACE, Literal: " "},
				{Type: tokens.DNUM, Literal: "1"},
			},
		},
		{
			name:  "greater than or equal",
			input: `1>=1`,
			want: []tokens.Token{
				{Type: tokens.DNUM, Literal: "1"},
				{Type: tokens.GEQ, Literal: ">="},
				{Type: tokens.DNUM, Literal: "1"},
			},
		},
		{
			name:  "greater than or equal spaces",
			input: `1 >= 1`,
			want: []tokens.Token{
				{Type: tokens.DNUM, Literal: "1"},
				{Type: tokens.WHITESPACE, Literal: " "},
				{Type: tokens.GEQ, Literal: ">="},
				{Type: tokens.WHITESPACE, Literal: " "},
				{Type: tokens.DNUM, Literal: "1"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := bufio.NewReader(strings.NewReader(tt.input))
			l := New(reader, "")
			for i, want := range tt.want {
				got, _ := l.NextToken()
				if got != want {
					t.Errorf("Test %d: got %q, want %q", i, got, want)
				}
			}
		})
	}
}

func TestLexer_NextToken_KeywordsAndIdentifiers(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  []tokens.Token
	}{
		{
			name:  "and",
			input: `foo and bar`,
			want: []tokens.Token{
				{Type: tokens.IDENT, Literal: "foo"},
				{Type: tokens.WHITESPACE, Literal: " "},
				{Type: tokens.AND, Literal: "and"},
				{Type: tokens.WHITESPACE, Literal: " "},
				{Type: tokens.IDENT, Literal: "bar"},
			},
		},
		{
			name:  "AND",
			input: `foo AND bar`,
			want: []tokens.Token{
				{Type: tokens.IDENT, Literal: "foo"},
				{Type: tokens.WHITESPACE, Literal: " "},
				{Type: tokens.AND, Literal: "AND"},
				{Type: tokens.WHITESPACE, Literal: " "},
				{Type: tokens.IDENT, Literal: "bar"},
			},
		},
		{
			name:  "or",
			input: `foo or bar`,
			want: []tokens.Token{
				{Type: tokens.IDENT, Literal: "foo"},
				{Type: tokens.WHITESPACE, Literal: " "},
				{Type: tokens.OR, Literal: "or"},
				{Type: tokens.WHITESPACE, Literal: " "},
				{Type: tokens.IDENT, Literal: "bar"},
			},
		},
		{
			name:  "OR",
			input: `foo OR bar`,
			want: []tokens.Token{
				{Type: tokens.IDENT, Literal: "foo"},
				{Type: tokens.WHITESPACE, Literal: " "},
				{Type: tokens.OR, Literal: "OR"},
				{Type: tokens.WHITESPACE, Literal: " "},
				{Type: tokens.IDENT, Literal: "bar"},
			},
		},
		{
			name:  "NOT",
			input: `NOT bar`,
			want: []tokens.Token{
				{Type: tokens.NOT, Literal: "NOT"},
				{Type: tokens.WHITESPACE, Literal: " "},
				{Type: tokens.IDENT, Literal: "bar"},
			},
		},
		{
			name:  "preps",
			input: `FROM a TO b INTO c IN d BY e OF f WITH g USING h`,
			want: []tokens.Token{
				{Type: tokens.PREPOSITION, Literal: "FROM"},
				{Type: tokens.WHITESPACE, Literal: " "},
				{Type: tokens.IDENT, Literal: "a"},
				{Type: tokens.WHITESPACE, Literal: " "},
				{Type: tokens.PREPOSITION, Literal: "TO"},
				{Type: tokens.WHITESPACE, Literal: " "},
				{Type: tokens.IDENT, Literal: "b"},
				{Type: tokens.WHITESPACE, Literal: " "},
				{Type: tokens.PREPOSITION, Literal: "INTO"},
				{Type: tokens.WHITESPACE, Literal: " "},
				{Type: tokens.IDENT, Literal: "c"},
				{Type: tokens.WHITESPACE, Literal: " "},
				{Type: tokens.PREPOSITION, Literal: "IN"},
				{Type: tokens.WHITESPACE, Literal: " "},
				{Type: tokens.IDENT, Literal: "d"},
				{Type: tokens.WHITESPACE, Literal: " "},
				{Type: tokens.PREPOSITION, Literal: "BY"},
				{Type: tokens.WHITESPACE, Literal: " "},
				{Type: tokens.IDENT, Literal: "e"},
				{Type: tokens.WHITESPACE, Literal: " "},
				{Type: tokens.PREPOSITION, Literal: "OF"},
				{Type: tokens.WHITESPACE, Literal: " "},
				{Type: tokens.IDENT, Literal: "f"},
				{Type: tokens.WHITESPACE, Literal: " "},
				{Type: tokens.PREPOSITION, Literal: "WITH"},
				{Type: tokens.WHITESPACE, Literal: " "},
				{Type: tokens.IDENT, Literal: "g"},
				{Type: tokens.WHITESPACE, Literal: " "},
				{Type: tokens.PREPOSITION, Literal: "USING"},
				{Type: tokens.WHITESPACE, Literal: " "},
				{Type: tokens.IDENT, Literal: "h"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := bufio.NewReader(strings.NewReader(tt.input))
			l := New(reader, "")
			for i, want := range tt.want {
				got, _ := l.NextToken()
				if got != want {
					t.Errorf("Test %d: got %q, want %q", i, got, want)
				}
			}
		})
	}
}

func TestLexer_NextToken_Comments(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  []tokens.Token
	}{
		{
			name:  "single line comment form .",
			input: `. this is a comment`,
			want: []tokens.Token{
				{Type: tokens.COMMENT, Literal: ". this is a comment"},
			},
		},
		{
			name: "single line comment form.  with newline",
			input: ` . this is a comment
`,
			want: []tokens.Token{
				{Type: tokens.WHITESPACE, Literal: " "},
				{Type: tokens.COMMENT, Literal: ". this is a comment"},
				{Type: tokens.EOF, Literal: "\x00"},
			},
		},
		{
			name: "single line comment form . enclosed with code",
			input: `foo

   			. this is a comment
bar`,
			want: []tokens.Token{
				{Type: tokens.IDENT, Literal: "foo"},
				{Type: tokens.NEWLINE, Literal: "\n"},
				{Type: tokens.NULLLINE, Literal: "\n"},
				{Type: tokens.WHITESPACE, Literal: "   \t\t\t"},
				{Type: tokens.COMMENT, Literal: ". this is a comment"},
				{Type: tokens.IDENT, Literal: "bar"},
				{Type: tokens.EOF, Literal: "\x00"},
			},
		},
		{
			name:  "single line comment form *",
			input: `* this is a comment`,
			want: []tokens.Token{
				{Type: tokens.COMMENT, Literal: "* this is a comment"},
			},
		},
		{
			name: "single line comment form * with newline",
			input: ` * this is a comment
`,
			want: []tokens.Token{
				{Type: tokens.WHITESPACE, Literal: " "},
				{Type: tokens.COMMENT, Literal: "* this is a comment"},
				{Type: tokens.EOF, Literal: "\x00"},
			},
		},
		{
			name: "single line comment form * enclosed with code",
			input: `foo

   			* this is a comment
bar`,
			want: []tokens.Token{
				{Type: tokens.IDENT, Literal: "foo"},
				{Type: tokens.NEWLINE, Literal: "\n"},
				{Type: tokens.NULLLINE, Literal: "\n"},
				{Type: tokens.WHITESPACE, Literal: "   \t\t\t"},
				{Type: tokens.COMMENT, Literal: "* this is a comment"},
				{Type: tokens.IDENT, Literal: "bar"},
				{Type: tokens.EOF, Literal: "\x00"},
			},
		},
		{
			name:  "single line comment form +",
			input: `+ this is a comment`,
			want: []tokens.Token{
				{Type: tokens.COMMENT, Literal: "+ this is a comment"},
			},
		},
		{
			name: "single line comment form + with newline",
			input: ` + this is a comment
`,
			want: []tokens.Token{
				{Type: tokens.WHITESPACE, Literal: " "},
				{Type: tokens.COMMENT, Literal: "+ this is a comment"},
				{Type: tokens.EOF, Literal: "\x00"},
			},
		},
		{
			name: "single line comment form +gc enclosed with code",
			input: `foo

   			+ this is a comment
bar`,
			want: []tokens.Token{
				{Type: tokens.IDENT, Literal: "foo"},
				{Type: tokens.NEWLINE, Literal: "\n"},
				{Type: tokens.NULLLINE, Literal: "\n"},
				{Type: tokens.WHITESPACE, Literal: "   \t\t\t"},
				{Type: tokens.COMMENT, Literal: "+ this is a comment"},
				{Type: tokens.IDENT, Literal: "bar"},
				{Type: tokens.EOF, Literal: "\x00"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := bufio.NewReader(strings.NewReader(tt.input))
			l := New(reader, "")
			for i, want := range tt.want {
				got, _ := l.NextToken()
				if got != want {
					t.Errorf("Test %d: got %q, want %q", i, got, want)
				}
			}
		})
	}
}

func TestLexer_NextToken_Strings(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  []tokens.Token
	}{
		{
			name:  "single line string",
			input: `"this is a string"`,
			want: []tokens.Token{
				{Type: tokens.LITERAL, Literal: `this is a string`},
			},
		},
		{
			name: "single line string with newline",
			input: `"this is a string"
`,
			want: []tokens.Token{
				{Type: tokens.LITERAL, Literal: `this is a string`},
				{Type: tokens.NEWLINE, Literal: "\n"},
				{Type: tokens.EOF, Literal: "\x00"},
			},
		},
		{
			name: "single line string enclosed with code",
			input: `foo

   			"this is a string"
bar`,
			want: []tokens.Token{
				{Type: tokens.IDENT, Literal: "foo"},
				{Type: tokens.NEWLINE, Literal: "\n"},
				{Type: tokens.NULLLINE, Literal: "\n"},
				{Type: tokens.WHITESPACE, Literal: "   \t\t\t"},
				{Type: tokens.LITERAL, Literal: `this is a string`},
				{Type: tokens.NEWLINE, Literal: "\n"},
				{Type: tokens.IDENT, Literal: "bar"},
				{Type: tokens.EOF, Literal: "\x00"},
			},
		},
		{
			name:  "single line string with escaped quote",
			input: `"this is a ""string"""`,
			want: []tokens.Token{
				{Type: tokens.LITERAL, Literal: `this is a "string"`},
			},
		},
		{
			name:  "single line string with escaped quote 2",
			input: `"this is a ""string"" with some more text"`,
			want: []tokens.Token{
				{Type: tokens.LITERAL, Literal: `this is a "string" with some more text`},
			},
		},
		{
			name:  "numericliteral starting with '-'",
			input: `"-1"`,
			want: []tokens.Token{
				{Type: tokens.NUMERICLITERAL, Literal: `-1`},
			},
		},
		{
			name:  "numericliteral starting with '.'",
			input: `".1"`,
			want: []tokens.Token{
				{Type: tokens.NUMERICLITERAL, Literal: `.1`},
			},
		},
		{
			name:  "numericliteral starting with '1'",
			input: `"1189"`,
			want: []tokens.Token{
				{Type: tokens.NUMERICLITERAL, Literal: `1189`},
			},
		},
		{
			name:  "literal starting with '-'",
			input: `"-asdf"`,
			want: []tokens.Token{
				{Type: tokens.LITERAL, Literal: `-asdf`},
			},
		},
		{
			name:  "literal starting with '.'",
			input: `".asdf"`,
			want: []tokens.Token{
				{Type: tokens.LITERAL, Literal: `.asdf`},
			},
		},
		{
			name:  "literal starting with '1'",
			input: `"1asdf"`,
			want: []tokens.Token{
				{Type: tokens.LITERAL, Literal: `1asdf`},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := bufio.NewReader(strings.NewReader(tt.input))
			l := New(reader, "")
			for i, want := range tt.want {
				got, _ := l.NextToken()
				if got != want {
					t.Errorf("Test %d: got %q, want %q", i, got, want)
				}
			}
		})
	}
}
