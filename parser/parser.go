package parser

import (
	"PLB-Interpreter/ast"
	"PLB-Interpreter/lexer"
	"PLB-Interpreter/plbErrors"
	"PLB-Interpreter/tokens"
	"fmt"
)

type prefixParseFn func() ast.Expression
type infixParseFn func(ast.Expression) ast.Expression

type Parser struct {
	l *lexer.Lexer

	errors []error

	curToken   tokens.Token
	peekToken  tokens.Token
	peekToken2 tokens.Token

	prefixParseFns map[tokens.TokenType]prefixParseFn
	infixParseFns  map[tokens.TokenType]infixParseFn
}

// Advances the parser by one token, setting the current token to the peek token
func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	peek2, err := p.l.NextToken()
	if err != nil {
		p.errors = append(p.errors, err)
		return
	}
	p.peekToken = p.peekToken2
	p.peekToken2 = peek2
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{l: l}
	p.nextToken()
	p.nextToken()
	p.nextToken()
	return p
}

// Errors returns true if there are any errors indicated in the parser
// If the boolean is true, the slice of errors will be non-empty
// If the boolean is false, the slice of errors will be empty
func (p *Parser) Errors() (bool, []error) {
	return len(p.errors) > 0, p.errors
}

func (p *Parser) addError(code, msg string) {
	newErr := plbErrors.NewPLBError(
		code,
		msg,
		p.curToken.FileName,
		p.curToken.Line,
		p.curToken.Col,
		p.curToken.LineTxt,
	)
	p.errors = append(p.errors, newErr)
}

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	for p.curToken.Type != tokens.EOF {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}

	return program
}

func (p *Parser) parseStatement() ast.Statement {
	if p.isValidStatement() {
		fmt.Printf("Line: %d  is a statement line: %s", p.curToken.Line, p.curToken.LineTxt)
		p.consumeTillNewline()
	}
	return nil
}

func (p *Parser) isValidStatement() bool {
	if (p.curToken.Type == tokens.IDENT && p.peekToken.Type == tokens.NEWLINE) ||
		(p.curToken.Type == tokens.IDENT && p.peekToken.Type == tokens.WHITESPACE && p.peekToken2.Type == tokens.NEWLINE) {
		// This is a label line
		return true
	}
	if (p.curToken.Type == tokens.IDENT && p.peekToken.Type == tokens.WHITESPACE && p.peekToken2.Type == tokens.IDENT) ||
		(p.curToken.Type == tokens.WHITESPACE && p.peekToken.Type == tokens.IDENT) {
		// This is a statement line
		return true
	}
	return false
}

func (p *Parser) consumeTillNewline() {
	for p.curToken.Type != tokens.NEWLINE {
		p.nextToken()
	}
}
