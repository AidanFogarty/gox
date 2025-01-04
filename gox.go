package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"

	gox "github.com/AidanFogarty/gox/pkg"
)

func main() {
	if len(os.Args) > 2 {
		fmt.Println("Usage: gox [script]")
		os.Exit(1)
	} else if len(os.Args) == 2 {
		if os.Args[1] == "expr" {
			expr := gox.NewBinary(
				gox.NewUnary(
					gox.NewToken(gox.Minus, "-", nil, 1),
					gox.NewLiteral(123),
				),
				gox.NewToken(gox.Star, "*", nil, 1),
				gox.NewGrouping(
					gox.NewLiteral(45.67),
				),
			)

			printer := gox.NewAstPrinter()
			fmt.Println(expr.Accept(printer))
		} else {
			scriptPath, err := filepath.Abs(os.Args[1])
			if err != nil {
				fmt.Println("error: unable to get absolute path")
				os.Exit(1)
			}
			runFile(scriptPath)
		}
	} else {
		runPrompt()
	}
}

func runFile(path string) {
	data, err := os.ReadFile(path)
	if err != nil {
		fmt.Println("error: unable to read file")
		os.Exit(1)
	}
	run(string(data))
}

func runPrompt() {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("> ")
		code, _, err := reader.ReadLine()
		if err != nil {
			fmt.Println("error: unable to read line")
			os.Exit(1)
		}
		run(string(code))
	}
}

func run(code string) {
	lexer := gox.NewLexer(code)
	lexer.Lex()

	parser := gox.NewParser(lexer.Tokens)
	expr, err := parser.Parse()
	if err != nil {
		fmt.Println(err)
	}

	printer := gox.NewAstPrinter()
	fmt.Println(expr.Accept(printer))
}
