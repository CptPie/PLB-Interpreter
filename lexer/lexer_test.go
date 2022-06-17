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
		{token.IDENT, "five"},
		{token.ASSIGN, "="},
		{token.INT, "5"},
		{token.SEMICOLON, ";"},
		{token.LET, "let"},
		{token.IDENT, "float"},
		{token.ASSIGN, "="},
		{token.FLOAT, "5.5"},
		{token.SEMICOLON, ";"},
		{token.LET, "let"},
		{token.IDENT, "hex"},
		{token.ASSIGN, "="},
		{token.HEX, "0x5"},
		{token.SEMICOLON, ";"},
		{token.LET, "let"},
		{token.IDENT, "oct"},
		{token.ASSIGN, "="},
		{token.OCT, "05"},
		{token.SEMICOLON, ";"},
		{token.LET, "let"},
		{token.IDENT, "ten"},
		{token.ASSIGN, "="},
		{token.INT, "10"},
		{token.SEMICOLON, ";"},
		{token.LET, "let"},
		{token.IDENT, "negten"},
		{token.ASSIGN, "="},
		{token.INT, "-10"},
		{token.SEMICOLON, ";"},
		{token.LET, "let"},
		{token.IDENT, "negfloat"},
		{token.ASSIGN, "="},
		{token.FLOAT, "-10.0"},
		{token.SEMICOLON, ";"},
		{token.LET, "let"},
		{token.IDENT, "add"},
		{token.ASSIGN, "="},
		{token.FUNCTION, "fn"},
		{token.LPAREN, "("},
		{token.IDENT, "x"},
		{token.COMMA, ","},
		{token.IDENT, "y"},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.IDENT, "x"},
		{token.PLUS, "+"},
		{token.IDENT, "y"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},
		{token.SEMICOLON, ";"},
		{token.LET, "let"},
		{token.IDENT, "result"},
		{token.ASSIGN, "="},
		{token.IDENT, "add"},
		{token.LPAREN, "("},
		{token.IDENT, "five"},
		{token.COMMA, ","},
		{token.IDENT, "ten"},
		{token.RPAREN, ")"},
		{token.SEMICOLON, ";"},
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
