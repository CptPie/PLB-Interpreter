package lexer

import (
	"PLB/token"
	"os"
	"testing"
)

// TODO split into more tests

func TestNextToken_SimpleTokens(t *testing.T) {
	rawInput, err := os.ReadFile("../sources/simpleTokens")
	if err != nil {
		t.Fatalf("Error while opening test source file: %s", err.Error())
	}
	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.LET, "let"},
		{token.WHITESPACE, " "},
		{token.INT, "5"},
		{token.WHITESPACE, " "},
		{token.PREPOSITION, "into"},
		{token.WHITESPACE, " "},
		{token.LABEL, "integer"},
		{token.NEWLINE, "\n"},

		{token.LET, "let"},
		{token.WHITESPACE, " "},
		{token.FLOAT, "5.5"},
		{token.WHITESPACE, " "},
		{token.PREPOSITION, "to"},
		{token.WHITESPACE, " "},
		{token.LABEL, "float"},
		{token.NEWLINE, "\n"},

		{token.LET, "let"},
		{token.WHITESPACE, " "},
		{token.INT, "-5"},
		{token.WHITESPACE, " "},
		{token.PREPOSITION, "INTO"},
		{token.WHITESPACE, " "},
		{token.LABEL, "negative_int"},
		{token.NEWLINE, "\n"},

		{token.LET, "let"},
		{token.WHITESPACE, " "},
		{token.FLOAT, "-5.5"},
		{token.WHITESPACE, " "},
		{token.PREPOSITION, "TO"},
		{token.WHITESPACE, " "},
		{token.LABEL, "negativeFloat"},
		{token.NEWLINE, "\n"},

		{token.LET, "let"},
		{token.WHITESPACE, " "},
		{token.OCT, "07"},
		{token.WHITESPACE, " "},
		{token.PREPOSITION, "in"},
		{token.WHITESPACE, " "},
		{token.LABEL, "oct"},
		{token.NEWLINE, "\n"},

		{token.LET, "let"},
		{token.WHITESPACE, " "},
		{token.HEX, "0x1283F"},
		{token.WHITESPACE, " "},
		{token.PREPOSITION, "IN"},
		{token.WHITESPACE, " "},
		{token.LABEL, "hex"},
		{token.NEWLINE, "\n"},

		{token.EOF, ""},
	}
	l := New(string(rawInput))
	for i, tt := range tests {
		tok := l.NextToken()
		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q, literal=%q",
				i, tt.expectedType, tok.Type, tok.Literal)
		}
		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}

	}
}
func TestNextToken_Operators(t *testing.T) {
	rawInput, err := os.ReadFile("../sources/operators")
	if err != nil {
		t.Fatalf("Error while opening test source file: %s", err.Error())
	}
	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.CURRENCY, "$"},
		{token.WHITESPACE, " "},
		{token.FORCING, "#"},
		{token.NEWLINE, "\n"},

		{token.PLUS, "+"},
		{token.WHITESPACE, " "},
		{token.MINUS, "-"},
		{token.WHITESPACE, " "},
		{token.MULTIPLY, "*"},
		{token.WHITESPACE, " "},
		{token.DIVIDE, "/"},
		{token.WHITESPACE, " "},
		{token.EXPONENT, "**"},
		{token.NEWLINE, "\n"},

		{token.EQ, "="},
		{token.WHITESPACE, " "},
		{token.NEQ, "<>"},
		{token.WHITESPACE, " "},
		{token.LT, "<"},
		{token.WHITESPACE, " "},
		{token.GT, ">"},
		{token.WHITESPACE, " "},
		{token.LTEQ, "<="},
		{token.WHITESPACE, " "},
		{token.GTEQ, ">="},
		{token.NEWLINE, "\n"},

		{token.AND, "AND"},
		{token.WHITESPACE, " "},
		{token.AND, "and"},
		{token.WHITESPACE, " "},
		{token.OR, "OR"},
		{token.WHITESPACE, " "},
		{token.OR, "or"},
		{token.WHITESPACE, " "},
		{token.NOT, "NOT"},
		{token.WHITESPACE, " "},
		{token.NOT, "not"},
		{token.NEWLINE, "\n"},

		{token.EOF, ""},
	}
	l := New(string(rawInput))
	for i, tt := range tests {
		tok := l.NextToken()
		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q, literal=%q",
				i, tt.expectedType, tok.Type, tok.Literal)
		}
		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}

	}
}
