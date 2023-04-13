package main

import (
	"PLB-Interpreter/lexer"
	"PLB-Interpreter/tokens"
	"bufio"
	"fmt"
	"os"
)

func main() {
	path := "examples/test.plb"
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	lex := lexer.New(reader, path)

	var toks []tokens.Token

	for {
		tok, err := lex.NextToken()

		toks = append(toks, tok)

		if err != nil {
			panic(err)
		}
		if tok.Type == tokens.EOF {
			break
		}
	}

	fmt.Println(toks)
}
