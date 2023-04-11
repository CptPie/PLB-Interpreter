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
			input: "   \t\t\t",
			want: []tokens.Token{
				{Type: tokens.WHITESPACE, Literal: " "},
				{Type: tokens.WHITESPACE, Literal: " "},
				{Type: tokens.WHITESPACE, Literal: " "},
				{Type: tokens.WHITESPACE, Literal: "\t"},
				{Type: tokens.WHITESPACE, Literal: "\t"},
				{Type: tokens.WHITESPACE, Literal: "\t"},
			},
		},
		{
			name:  "Newlines and carriage returns",
			input: "\n\r\n\r",
			want: []tokens.Token{
				{Type: tokens.NEWLINE, Literal: "\n"},
				{Type: tokens.NEWLINE, Literal: "\r"},
				{Type: tokens.NEWLINE, Literal: "\n"},
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
			want:  &plbErrors.PLBError{ErrorCode: "Lexer", Message: "Invalid token type", File: "test", Line: 1, Column: 1, LineText: "!hello"},
		},
		{
			name:  "Invalid token 2",
			input: "1383!",
			want:  &plbErrors.PLBError{ErrorCode: "Lexer", Message: "Invalid token type", File: "test", Line: 1, Column: 5, LineText: "1383!"},
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
				{Type: tokens.ASTER, Literal: "*"},
				{Type: tokens.DNUM, Literal: "1"},
			},
		},
		{
			name:  "multiplication spaces",
			input: "1 * 1",
			want: []tokens.Token{
				{Type: tokens.DNUM, Literal: "1"},
				{Type: tokens.WHITESPACE, Literal: " "},
				{Type: tokens.ASTER, Literal: "*"},
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
				{Type: tokens.POW, Literal: "**"},
				{Type: tokens.DNUM, Literal: "1"},
			},
		},
		{
			name:  "exponent spaces",
			input: "1 ** 1",
			want: []tokens.Token{
				{Type: tokens.DNUM, Literal: "1"},
				{Type: tokens.WHITESPACE, Literal: " "},
				{Type: tokens.POW, Literal: "**"},
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
