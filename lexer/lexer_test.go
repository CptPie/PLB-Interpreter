package lexer

import (
	plbErrors "PLB-Interpreter/plbErrors"
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
		{
			name:  "Signed decimal numbers",
			input: "-123\n-456",
			want: []tokens.Token{
				{Type: tokens.SIGNEDDNUM, Literal: "-123"},
				{Type: tokens.NEWLINE, Literal: "\n"},
				{Type: tokens.SIGNEDDNUM, Literal: "-456"},
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
