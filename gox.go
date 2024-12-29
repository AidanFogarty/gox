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
		scriptPath, err := filepath.Abs(os.Args[1])
		if err != nil {
			fmt.Println("error: unable to get absolute path")
			os.Exit(1)
		}
		runFile(scriptPath)
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

	for _, token := range lexer.Tokens {
		fmt.Println(token)
	}
}
