package lexer

import (
	"PLB/token"
	"testing"
)

// TODO split into more tests

func TestNextToken(t *testing.T) {
	input := `let five = 5;
let float = 5.5;
let hex = 0x5;
let oct = 05;
let ten = 10;
let negten = -10;
let negfloat = -10.0;
let add = fn(x, y) {
x + y;
};
let result = add(five, ten);
`
	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.LET, "let"},
		{token.LABEL, "five"},
		{token.ASSIGN, "="},
		{token.INT, "5"},
		{token.SEMICOLON, ";"},
		{token.NEWLINE, "\n"},
		{token.LET, "let"},
		{token.LABEL, "float"},
		{token.ASSIGN, "="},
		{token.FLOAT, "5.5"},
		{token.SEMICOLON, ";"},
		{token.NEWLINE, "\n"},
		{token.LET, "let"},
		{token.LABEL, "hex"},
		{token.ASSIGN, "="},
		{token.HEX, "0x5"},
		{token.SEMICOLON, ";"},
		{token.NEWLINE, "\n"},
		{token.LET, "let"},
		{token.LABEL, "oct"},
		{token.ASSIGN, "="},
		{token.OCT, "05"},
		{token.SEMICOLON, ";"},
		{token.NEWLINE, "\n"},
		{token.LET, "let"},
		{token.LABEL, "ten"},
		{token.ASSIGN, "="},
		{token.INT, "10"},
		{token.SEMICOLON, ";"},
		{token.NEWLINE, "\n"},
		{token.LET, "let"},
		{token.LABEL, "negten"},
		{token.ASSIGN, "="},
		{token.INT, "-10"},
		{token.SEMICOLON, ";"},
		{token.NEWLINE, "\n"},
		{token.LET, "let"},
		{token.LABEL, "negfloat"},
		{token.ASSIGN, "="},
		{token.FLOAT, "-10.0"},
		{token.SEMICOLON, ";"},
		{token.NEWLINE, "\n"},
		{token.LET, "let"},
		{token.LABEL, "add"},
		{token.ASSIGN, "="},
		{token.FUNCTION, "fn"},
		{token.LPAREN, "("},
		{token.LABEL, "x"},
		{token.COMMA, ","},
		{token.LABEL, "y"},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.NEWLINE, "\n"},
		{token.LABEL, "x"},
		{token.PLUS, "+"},
		{token.LABEL, "y"},
		{token.SEMICOLON, ";"},
		{token.NEWLINE, "\n"},
		{token.RBRACE, "}"},
		{token.SEMICOLON, ";"},
		{token.NEWLINE, "\n"},
		{token.LET, "let"},
		{token.LABEL, "result"},
		{token.ASSIGN, "="},
		{token.LABEL, "add"},
		{token.LPAREN, "("},
		{token.LABEL, "five"},
		{token.COMMA, ","},
		{token.LABEL, "ten"},
		{token.RPAREN, ")"},
		{token.SEMICOLON, ";"},
		{token.NEWLINE, "\n"},
		{token.EOF, ""},
	}
	l := New(input)
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
