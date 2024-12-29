package gox

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

type Lexer struct {
	Source  string
	Reader  *bufio.Reader
	Tokens  []*Token
	Start   int
	Current int
	Line    int
}

func NewLexer(source string) *Lexer {
	return &Lexer{Source: source, Reader: bufio.NewReader(strings.NewReader(source)), Tokens: []*Token{}, Start: 0, Current: 0, Line: 1}
}

func (l *Lexer) Lex() {
	for {
		r, _, err := l.Reader.ReadRune()
		if err != nil {
			if err == io.EOF {
				l.Tokens = append(l.Tokens, NewToken(TokenEOF, "EOF", nil, l.Line))
				return
			}

			panic(err)
		}

		l.Start = l.Current
		l.Current++

		switch r {

		// Single-character tokens
		case '(':
			l.addToken(LeftParen, "(", nil)
		case ')':
			l.addToken(RightParen, ")", nil)
		case '{':
			l.addToken(LeftBrace, "{", nil)
		case '}':
			l.addToken(RightBrace, "}", nil)
		case ',':
			l.addToken(Comma, ",", nil)
		case '.':
			l.addToken(Dot, ".", nil)
		case '-':
			l.addToken(Minus, "-", nil)
		case '+':
			l.addToken(Plus, "+", nil)
		case ';':
			l.addToken(Semicolon, ";", nil)
		case '*':
			l.addToken(Star, "*", nil)

		// One or two character tokens
		case '!':
			l.addTokenWithPeek(Bang, BangEqual, '=')
		case '=':
			l.addTokenWithPeek(Equal, EqualEqual, '=')
		case '<':
			l.addTokenWithPeek(Less, LessEqual, '=')
		case '>':
			l.addTokenWithPeek(Greater, GreaterEqual, '=')

		// Misc characters
		case ' ':
		case '\r':
		case '\t':
		case '\n':
			l.Line++

		case '"':
			l.string()

		case '/':
			// Handle comments as the division operator and comments share the same start character
			if l.peek() == '/' {
				l.advance()
				for l.peek() != '\n' && l.peek() != 0 {
					l.advance()
				}
			} else {
				l.addToken(Slash, "/", nil)
			}

		default:
			fmt.Println("unknown character: ", string(r))
		}
	}
}

func (l *Lexer) addToken(tokenType TokenType, lexeme string, literal interface{}) {
	l.Tokens = append(l.Tokens, NewToken(tokenType, lexeme, literal, l.Line))
}

func (l *Lexer) addTokenWithPeek(tokenType TokenType, otherTokenType TokenType, expectedRune rune) {
	r := l.peek()

	if r == expectedRune {
		l.advance()
		l.addToken(otherTokenType, l.Source[l.Start:l.Current], nil)
	} else {
		l.addToken(tokenType, l.Source[l.Start:l.Current], nil)
	}
}

func (l *Lexer) advance() rune {
	r, _, err := l.Reader.ReadRune()
	if err != nil {
		panic(err)
	}

	l.Current++
	return r
}

func (l *Lexer) peek() rune {
	r, _, err := l.Reader.ReadRune()
	if err != nil {
		if err == io.EOF {
			return 0
		}

		panic(fmt.Sprintf("failed to peek: %v", err))
	}

	err = l.Reader.UnreadRune()
	if err != nil {
		panic(fmt.Sprintf("failed to peek unread: %v", err))
	}

	return r
}

func (l *Lexer) string() {
	for l.peek() != '"' && l.peek() != 0 {
		if l.peek() == '\n' {
			l.Line++
		}

		l.advance()
	}

	if l.peek() == 0 {
		fmt.Println("error: unterminated string")
		return
	}

	l.advance() // handle the closing quote

	content := l.Source[l.Start+1 : l.Current-1]
	l.addToken(String, content, content)
}
