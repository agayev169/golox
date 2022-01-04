package golox_test

import (
	"testing"

	"github.com/agayev169/golox"
)

type parserTestDto struct {
	Tokens   []golox.Token
	Expected golox.Expr
}

var parserTestData = map[string]parserTestDto{
	"simple": {
		Tokens: []golox.Token{
			{Type: golox.MINUS, Lexeme: "-"},
			{Type: golox.NUMBER, Lexeme: "123", Literal: 123},
			{Type: golox.STAR, Lexeme: "*"},
			{Type: golox.LEFT_PAREN, Lexeme: "("},
			{Type: golox.NUMBER, Lexeme: "45.67", Literal: 45.67},
			{Type: golox.RIGHT_PAREN, Lexeme: ")"},
		},
		Expected: &golox.Binary{
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
	},
}

func TestParser(t *testing.T) {
	for k, tv := range parserTestData {
		p := golox.NewParser(tv.Tokens)
		actual, err := p.Parse()
		if err != nil {
			t.Fatalf("Failed on test %s. Got error on p.Parse(): %s", k, err.Error())
		}

		if !areEqualExprs(actual, tv.Expected) {
			t.Fatalf("Failed on test %s\n", k)
		}
	}
}