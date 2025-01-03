package gox

// equality → comparison ( ( "!=" | "==" ) comparison )* ;

type Parser struct {
	Tokens  []*Token
	Current int
}

func NewParser(tokens []*Token) *Parser {
	return &Parser{Tokens: tokens, Current: 0}
}

func (p *Parser) match(types ...TokenType) bool {
	return false
}

func (p *Parser) expression() Expr {
	return nil

}
