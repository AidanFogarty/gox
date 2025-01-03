package gox

import "fmt"

type AstPrinter struct{}

func NewAstPrinter() *AstPrinter {
	return &AstPrinter{}
}

func (a *AstPrinter) Print(expr Expr) string {
	return expr.Accept(a).(string)
}

func (a *AstPrinter) parenthesize(name string, exprs ...Expr) string {
	ast := ""

	ast += "(" + name

	for _, expr := range exprs {
		ast += " " + expr.Accept(a).(string)
	}

	ast += ")"

	return ast
}

func (a *AstPrinter) VisitBinary(expr *Binary) interface{} {
	return a.parenthesize(expr.Operator.Lexeme, expr.Left, expr.Right)
}

func (a *AstPrinter) VisitGrouping(expr *Grouping) interface{} {
	return a.parenthesize("group", expr.Expression)
}

func (a *AstPrinter) VisitLiteral(expr *Literal) interface{} {
	return fmt.Sprintf("%v", expr.Value)
}

func (a *AstPrinter) VisitUnary(expr *Unary) interface{} {
	return a.parenthesize(expr.Operator.Lexeme, expr.Right)
}
