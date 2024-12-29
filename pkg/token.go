package gox

import "fmt"

type TokenType int

const (
	TokenEOF TokenType = iota

	// Single-character tokens
	LeftParen
	RightParen
	LeftBrace
	RightBrace
	Comma
	Dot
	Minus
	Plus
	Semicolon
	Slash
	Star

	// One or two character tokens
	Bang
	BangEqual
	Equal
	EqualEqual
	Greater
	GreaterEqual
	Less
	LessEqual

	// Literals
	Identifier
	String
	Number

	// Keywords
	And
	Class
	Else
	False
	Fun
	For
	If
	Nil
	Or
	Print
	Return
	Super
	This
	True
	Var
	While
)

type Token struct {
	Type    TokenType
	Lexeme  string
	Literal interface{}
	Line    int
}

func NewToken(t TokenType, lexeme string, literal interface{}, line int) *Token {
	return &Token{Type: t, Lexeme: lexeme, Literal: literal, Line: line}
}

func (t Token) String() string {
	return fmt.Sprintf("%d %s", t.Type, t.Lexeme)
}
