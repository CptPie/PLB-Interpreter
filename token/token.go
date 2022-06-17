package token

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
}

const (
	ILLEGAL    = "ILLEGAL"
	WHITESPACE = "WHITESPACE" // {space} \t
	NEWLINE    = "NEWLINE"    // \n \r
	EOF        = "EOF"

	// Basic elements (doxnum)
	INT   = "INT"   // (signed)-dnum
	OCT   = "OCT"   // onum
	HEX   = "HEX"   // xnum
	FLOAT = "FLOAT" // numeric constant

	// LABELS + literals
	LABEL    = "LABEL"    // data_label
	LOCATION = "LOCATION" // execution_label
	LITERAL  = "LITERAL"  // literal (strings ...)

	// Special Chars
	CURRENCY = "CURRENCY" // $ (¤)
	FORCING  = "FORCING"  // # (£)

	// COMMENT LINE
	COMMENT = "COMMENT" // .{} / *{} / +{}

	// ARITHMETHIC OPERATORS
	PLUS     = "+"
	MINUS    = "-"
	MULTIPLY = "*"
	DIVIDE   = "/"
	EXPONENT = "**"

	// RELATIONAL OPERATORS
	EQUAL     = "="
	NOTEQUAL  = "<>"
	LESS      = "<"
	GREATER   = ">"
	LESSEQ    = "<="
	GREATEREQ = ">="

	// LOGICAL OPERATORS
	AND = "AND"
	OR  = "OR"
	NOT = "NOT"

	// Delimiters
	COMMA     = ","
	SEMICOLON = ";"
	LPAREN    = "("
	RPAREN    = ")"
	//	LBRACE    = "{"
	//	RBRACE    = "}"

	// Keywords
	PREPOSITION = "PREPOSITION"
	FUNCTION    = "FUNCTION"
	LET         = "LET"
)

var keywords = map[string]TokenType{
	"fn":  FUNCTION,
	"let": LET,

	"from":  PREPOSITION,
	"to":    PREPOSITION,
	"into":  PREPOSITION,
	"in":    PREPOSITION,
	"by":    PREPOSITION,
	"of":    PREPOSITION,
	"with":  PREPOSITION,
	"using": PREPOSITION,

	"FROM":  PREPOSITION,
	"TO":    PREPOSITION,
	"INTO":  PREPOSITION,
	"IN":    PREPOSITION,
	"BY":    PREPOSITION,
	"OF":    PREPOSITION,
	"WITH":  PREPOSITION,
	"USING": PREPOSITION,

	"and": AND,
	"AND": AND,
	"or":  OR,
	"OR":  OR,
	"not": NOT,
	"NOT": NOT,
}

func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return LABEL
}
