package gox

import "fmt"

// equality → comparison ( ( "!=" | "==" ) comparison )* ;

type ParseError struct {
	Token   *Token
	Message string
}

func (e *ParseError) Error() string {
	if e.Token.Type == TokenEOF {
		return fmt.Sprintf("%d at end %s", e.Token.Line, e.Message)
	}

	return fmt.Sprintf("%d at '%s' %s", e.Token.Line, e.Token.Lexeme, e.Message)
}

type Parser struct {
	Tokens  []*Token
	Current int
}

func NewParser(tokens []*Token) *Parser {
	return &Parser{Tokens: tokens, Current: 0}
}

func (p *Parser) peek() *Token {
	return p.Tokens[p.Current]
}

func (p *Parser) isAtEnd() bool {
	return p.peek().Type == TokenEOF
}

func (p *Parser) previous() *Token {
	return p.Tokens[p.Current-1]
}

func (p *Parser) check(tokenType TokenType) bool {
	if p.isAtEnd() {
		return false
	}

	return p.peek().Type == tokenType
}

func (p *Parser) advance() *Token {
	if !p.isAtEnd() {
		p.Current++
	}

	return p.previous()
}

func (p *Parser) match(types ...TokenType) bool {
	for _, tokenType := range types {
		if p.check(tokenType) {
			p.advance()
			return true
		}
	}

	return false
}

func (p *Parser) consume(tokenType TokenType, message string) (*Token, error) {
	if p.check(tokenType) {
		return p.advance(), nil
	}

	return nil, &ParseError{
		Token:   p.peek(),
		Message: message,
	}

}

func (p *Parser) primary() (Expr, error) {
	if p.match(False) {
		return &Literal{
			Value: false,
		}, nil
	}

	if p.match(True) {
		return &Literal{
			Value: true,
		}, nil
	}

	if p.match(Nil) {
		return &Literal{
			Value: nil,
		}, nil
	}

	if p.match(Number, String) {
		return &Literal{
			Value: p.previous().Literal,
		}, nil
	}

	if p.match(LeftParen) {
		expr, err := p.expression()
		if err != nil {
			return nil, err
		}

		p.consume(RightParen, "Expect ')' after expression.")
		return &Grouping{
			Expression: expr,
		}, nil
	}

	return nil, &ParseError{
		Token:   p.peek(),
		Message: "Expect expression.",
	}

}

func (p *Parser) unary() (Expr, error) {
	if p.match(Bang, Minus) {
		operator := p.previous()
		right, err := p.unary()
		if err != nil {
			return nil, err
		}

		return &Unary{
			Operator: operator,
			Right:    right,
		}, nil
	}

	return p.primary()
}

func (p *Parser) factor() (Expr, error) {
	expr, err := p.unary()
	if err != nil {
		return nil, err
	}

	for p.match(Star, Slash) {
		operator := p.previous()
		right, err := p.unary()
		if err != nil {
			return nil, err
		}

		expr = &Binary{
			Left:     expr,
			Operator: operator,
			Right:    right,
		}
	}

	return expr, nil
}

func (p *Parser) term() (Expr, error) {
	expr, err := p.factor()
	if err != nil {
		return nil, err
	}

	for p.match(Minus, Plus) {
		operator := p.previous()
		right, err := p.factor()
		if err != nil {
			return nil, err
		}

		expr = &Binary{
			Left:     expr,
			Operator: operator,
			Right:    right,
		}
	}

	return expr, nil
}

func (p *Parser) comparision() (Expr, error) {
	expr, err := p.term()
	if err != nil {
		return nil, err
	}

	for p.match(Greater, GreaterEqual, Less, LessEqual) {
		operator := p.previous()
		right, err := p.term()
		if err != nil {
			return nil, err
		}

		expr = &Binary{
			Left:     expr,
			Operator: operator,
			Right:    right,
		}
	}

	return expr, nil
}

func (p *Parser) equality() (Expr, error) {
	expr, err := p.comparision()
	if err != nil {
		return nil, err
	}

	for p.match(BangEqual, EqualEqual) {
		operator := p.previous()
		right, err := p.comparision()
		if err != nil {
			return nil, err
		}

		expr = &Binary{
			Left:     expr,
			Operator: operator,
			Right:    right,
		}
	}

	return expr, nil
}

func (p *Parser) expression() (Expr, error) {
	return p.equality()
}

func (p *Parser) synchronize() {
	p.advance()

	for !p.isAtEnd() {
		if p.previous().Type == Semicolon {
			return
		}

		switch p.peek().Type {
		case Class, Fun, Var, For, If, While, Print, Return:
			return
		}

		p.advance()
	}
}

func (p *Parser) Parse() (Expr, error) {
	expr, err := p.expression()
	if err != nil {
		return nil, err
	}

	return expr, nil
}
