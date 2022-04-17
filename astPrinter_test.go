package golox_test

import (
	"testing"

	"github.com/agayev169/golox"
)

type printerTestDto struct {
	Expression golox.Expr
	Expected   string
}

var astPrinterTestData = map[string]printerTestDto{
	"simple": printerTestDto{
		Expression: &golox.Binary{
			Left: &golox.Unary{
				Operator: golox.Token{
					Type:   golox.MINUS,
					Lexeme: "-",
				},
				Right: &golox.Literal{
					Value: 123,
				},
			},
			Operator: golox.Token{
				Type:   golox.STAR,
				Lexeme: "*",
			},
			Right: &golox.Grouping{
				Expr: &golox.Literal{
					Value: 45.67,
				},
			},
		},
		Expected: "(* (- 123) (group 45.67))",
	},
}

func TestPrinter(t *testing.T) {
	ap := &AstPrinter{}
	for k, tv := range astPrinterTestData {
		actual, err := tv.Expression.Accept(ap)

		if err != nil {
			t.Fatalf("Failed on test %s. Expression.Accept return non-nil error: %v\n", k, err.Error())
		}

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
