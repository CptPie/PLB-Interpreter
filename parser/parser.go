package parser

import (
	"PLB-Interpreter/ast"
	"PLB-Interpreter/lexer"
	"PLB-Interpreter/plbErrors"
	"PLB-Interpreter/tokens"
)

type prefixParseFn func() ast.Expression
type infixParseFn func(ast.Expression) ast.Expression

type Parser struct {
	l *lexer.Lexer

	errors []error

	curToken  tokens.Token
	peekToken tokens.Token

	prefixParseFns map[tokens.TokenType]prefixParseFn
	infixParseFns  map[tokens.TokenType]infixParseFn
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	peek, err := p.l.NextToken()
	if err != nil {
		p.errors = append(p.errors, err)
		return
	}
	p.peekToken = peek
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{l: l}
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
	//switch p.curToken.Type {
	//case tokens.LET:
	//	return p.parseLetStatement()
	//case tokens.RETURN:
	//	return p.parseReturnStatement()
	//default:
	//	return p.parseExpressionStatement()
	//}
	return nil
}
