package main

import (
	"PLB-Interpreter/lexer"
	"PLB-Interpreter/parser"
	"bufio"
	"fmt"
	"os"
)

func main() {
	path := "examples/strings.plb"
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	lex := lexer.New(reader, path)
	pars := parser.New(lex)

	prog := pars.ParseProgram()

	fmt.Println(prog)

	//for {
	//	tok, err := lex.NextToken()
	//
	//	if tok.Type == tokens.COMMENT || tok.Type == tokens.NEWLINE || tok.Type == tokens.NULLLINE {
	//		fmt.Println(tok)
	//	} else {
	//		fmt.Print(tok)
	//	}
	//
	//	if err != nil {
	//		panic(err)
	//	}
	//	if tok.Type == tokens.EOF {
	//		break
	//	}
	//}

}
