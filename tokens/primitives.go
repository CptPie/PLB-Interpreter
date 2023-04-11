package tokens

type TokenType string // Type of token

const (
	// Boilerplate
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	// Whitespace
	BLANK      = "BLANK"      // Space
	WHITESPACE = "WHITESPACE" // Space or multiple spaces
	NEWLINE    = "NEWLINE"    // \r or \n or \r\n

	// characters
	ALPHACHAR   = "ALPHACHAR"   // A-Z, a-z
	CURRENCY    = "CURRENCY"    // $ or universal currency symbol
	FORCING     = "FORCING"     // # or £
	COMMA       = "COMMA"       // , or :
	SEMICOLON   = "SEMICOLON"   // ;
	PREP        = "PREP"        // either COMMA or a PREPOSITION enclosed by WHITESPACE
	PREPOSITION = "PREPOSITION" // a word that is a preposition (see prepList)

	// Lines and strings
	ANYCHAR   = "ANYCHAR"   // Any character in the character set
	ANYSTRING = "ANYSTRING" // zero or more characters or WHITESPACE, or both, delimited by end of line
	NULLLINE  = "NULLLINE"  // A line with no characters, only WHITESPACE, delimited by end of line, subset of COMMENT
	COMMENT   = "COMMENT"   // A line with a comment, delimited by end of line, indicated by leading . * or +

	// Labels
	LABEL          = "LABEL"          // one or more ALPHACHAR, literal $, DIGIT, or _
	DATALABEL      = "DATALABEL"      // see LABEL, data label space
	EQUATELABEL    = "EQUATELABEL"    // see LABEL, data label space
	VARLABEL       = "VARLABEL"       // see LABEL, data label space
	EXECUTIONLABEL = "EXECUTIONLABEL" // see LABEL, execution label space

	// Constants
	DNUM       = "DNUM"       // a decimal number or label pointing to a decimal number, unsigned, 16 bit minimum
	SIGNEDDNUM = "SIGNEDDNUM" // a decimal number or label pointing to a decimal number, signed, 16 bit minimum
	ONUM       = "ONUM"       // indicated by a leading 0, an octal number or label pointing to an octal number, unsigned, minimum 0 through 0177777
	XNUM       = "XNUM"       // indicated by a leading 0x or 0X, a hexadecimal number or label pointing to a hexadecimal number, unsigned, minimum 0 through FFFF
	DOXNUM     = "DOXNUM"     // either DNUM, ONUM or XNUM, see above for details

	NUMERICCONSTANT   = "NUMERICCONSTANT"   // either DIGITS or DIGITS . DIGITS (where in the . form either DIGITS can be ommitted)
	LITERAL           = "LITERAL"           // a string literal, indicated by leading and trailing ", can contain 0 characters or a LITERALVALUE
	LITERALVALUE      = "LITERALVALUE"      // any number of LITERALCHAR or FORCING followed by ANYCHAR or two " in a row representing a literal "
	LITERALCHAR       = "LITERALCHAR"       // any character except FORCING or "
	SINGLECHARLITERAL = "SINGLECHARLITERAL" // a single character literal, indicated by leading and trailing ", can contain 0 characters or a consist of ANYCHAR

)

var prepList = []string{
	"FROM",
	"TO",
	"INTO",
	"IN",
	"BY",
	"OF",
	"WITH",
	"USING",
}

// Token is a token returned by the lexer
type Token struct {
	Type    TokenType
	Literal string
}
