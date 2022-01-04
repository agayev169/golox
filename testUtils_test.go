package golox_test

import (
	"fmt"
	"log"
	"strings"

	. "github.com/agayev169/golox"
)

type AstPrinter struct {
}

func (ap *AstPrinter) AcceptBinaryExpr(b *Binary) interface{} {
	return ap.parenthesize(b.Operator.Lexeme, b.Left, b.Right)
}

func (ap *AstPrinter) AcceptGroupingExpr(g *Grouping) interface{} {
	return ap.parenthesize("group", g.Expr)
}

func (ap *AstPrinter) AcceptLiteralExpr(l *Literal) interface{} {
	return fmt.Sprintf("%v", l.Value)
}

func (ap *AstPrinter) AcceptUnaryExpr(u *Unary) interface{} {
	return ap.parenthesize(u.Operator.Lexeme, u.Right)
}

func (ap *AstPrinter) parenthesize(name string, exprs ...Expr) string {
	sb := strings.Builder{}

	sb.WriteString(fmt.Sprintf("(%s", name))

	for _, expr := range exprs {
		sb.WriteString(" ")

		res := expr.Accept(ap)

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

	return e1.Accept(ap).(string) == e2.Accept(ap).(string)
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