package golox_test

import (
	"fmt"
	"log"
	"strings"

	. "github.com/agayev169/golox"
)

type AstPrinter struct {
}

func (ap *AstPrinter) AcceptBinaryExpr(b *Binary) (interface{}, error) {
	return ap.parenthesize(b.Operator.Lexeme, b.Left, b.Right), nil
}

func (ap *AstPrinter) AcceptGroupingExpr(g *Grouping) (interface{}, error) {
	return ap.parenthesize("group", g.Expr), nil
}

func (ap *AstPrinter) AcceptLiteralExpr(l *Literal) (interface{}, error) {
	return fmt.Sprintf("%v", l.Value), nil
}

func (ap *AstPrinter) AcceptUnaryExpr(u *Unary) (interface{}, error) {
	return ap.parenthesize(u.Operator.Lexeme, u.Right), nil
}

func (ap *AstPrinter) parenthesize(name string, exprs ...Expr) string {
	sb := strings.Builder{}

	sb.WriteString(fmt.Sprintf("(%s", name))

	for _, expr := range exprs {
		sb.WriteString(" ")

		res, _ := expr.Accept(ap)

		switch v := res.(type) {
		case string:
			sb.WriteString(v)
		default:
			log.Panicf("Expected the output of expr.Accept(ap) to be string, received %v\n", v)
		}
	}

	sb.WriteString(")")

	return sb.String()
}

// Comparators

func areEqualExprs(e1, e2 Expr) bool {
    if e1 == nil && e2 == nil {
        return true
    }

    if (e1 != nil && e2 == nil) || (e1 == nil && e2 != nil) {
        return false
    }

	ap := &AstPrinter{}

    l, _ := e1.Accept(ap)
    r, _ := e2.Accept(ap)

	return l.(string) == r.(string)
}

func areEqualLoxErrors(e1, e2 *LoxError) bool {
    if (e1 != nil && e2 == nil) || (e1 == nil && e2 != nil) {
        return false
    }

    if e1 == nil && e2 == nil {
        return true
    }

    return e1.Number == e2.Number
}