package golox

import (
	"testing"
)

type printerTestDto struct {
	Expression Expr
	Expected   string
}

var data = map[string]printerTestDto{
	"simple": printerTestDto{
		Expression: &Binary{
			Left: &Unary{
				Operator: Token{
					Type:   MINUS,
					Lexeme: "-",
				},
				Right: &Literal{
					Value: 123,
				},
			},
			Operator: Token{
				Type:   STAR,
				Lexeme: "*",
			},
			Right: &Grouping{
				Expr: &Literal{
					Value: 45.67,
				},
			},
		},
		Expected: "(* (- 123) (group 45.67))",
	},
}

func TestPrinter(t *testing.T) {
    ap := &AstPrinter{}
    for k, tv := range data {
        actual := tv.Expression.Accept(ap)

        switch v := actual.(type) {
        case string:
            if tv.Expected != v {
                t.Fatalf("Failed on test %s. Expected: %s, got: %s\n", k, tv.Expected, actual)
            }
        default:
            t.Fatalf("Expected the output of expr.Accept(ap) to be string, received %v\n", v)
        }
    }
}
