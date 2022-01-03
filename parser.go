package golox

import "fmt"

type Parser struct {
	tokens  []Token
	current int
}

func NewParser(tokens []Token) *Parser {
	return &Parser{tokens: tokens, current: 0}
}

func (p *Parser) Parse() (Expr, error) {
	p.current = 0

	return p.parseExpression()
}

func (p *Parser) parseExpression() (Expr, error) {
	return p.parseEquality()
}

func (p *Parser) parseEquality() (Expr, error) {
	expr, err := p.parseComparison()
	if err != nil {
		return nil, err
	}

	for p.peek(BANG_EQUAL, EQUAL_EQUAL) {
		t := p.getNextToken()

		right, err := p.parseComparison()
		if err != nil {
			return nil, err
		}

		expr = &Binary{Left: expr, Operator: t, Right: right}
	}

	return expr, nil
}

func (p *Parser) parseComparison() (Expr, error) {
	expr, err := p.parseTerm()
	if err != nil {
		return nil, err
	}

	for p.peek(LESS, LESS_EQUAL, GREATER, GREATER_EQUAL) {
		t := p.getNextToken()

		right, err := p.parseTerm()
		if err != nil {
			return nil, err
		}

		expr = &Binary{Left: expr, Operator: t, Right: right}
	}

	return expr, nil
}

func (p *Parser) parseTerm() (Expr, error) {
	expr, err := p.parseFactor()
	if err != nil {
		return nil, err
	}

	for p.peek(PLUS, MINUS) {
		t := p.getNextToken()

		right, err := p.parseFactor()
		if err != nil {
			return nil, err
		}

		expr = &Binary{Left: expr, Operator: t, Right: right}
	}

	return expr, nil
}

func (p *Parser) parseFactor() (Expr, error) {
	expr, err := p.parseUnary()
	if err != nil {
		return nil, err
	}

	for p.peek(STAR, SLASH) {
		t := p.getNextToken()

		right, err := p.parseUnary()
		if err != nil {
			return nil, err
		}

		expr = &Binary{Left: expr, Operator: t, Right: right}
	}

	return expr, nil
}

func (p *Parser) parseUnary() (Expr, error) {
	if p.peek(MINUS, BANG) {
		t := p.getNextToken()

		right, err := p.parseUnary()
		if err != nil {
			return nil, err
		}

		return &Unary{Operator: t, Right: right}, nil
	}

	return p.parsePrimary()
}

func (p *Parser) parsePrimary() (Expr, error) {
	t := p.getNextToken()

	if t.Type == NUMBER || t.Type == STRING {
		return &Literal{Value: t.Literal}, nil
	}

	if t.Type == TRUE {
		return &Literal{Value: true}, nil
	}

	if t.Type == FALSE {
		return &Literal{Value: false}, nil
	}

	if t.Type == NIL {
		return &Literal{Value: nil}, nil
	}

	if t.Type == LEFT_PAREN {
		expr, err := p.parseExpression()
		if err != nil {
			return nil, err
		}

		_, err = p.consume(RIGHT_PAREN)
		if err != nil {
			return nil, err
		}

		return &Grouping{Expr: expr}, nil
	}

	return nil, &ScanError{File: t.File, Line: t.Line, Col: t.Col, Msg: fmt.Sprintf("Expected `(` but found `%s`.", t.Lexeme)}
}

func (p *Parser) consume(tt TokenType) (*Token, error) {
	if p.isAtEnd() {
		lt := p.tokens[len(p.tokens)-1]
		return nil, &ScanError{File: lt.File, Line: lt.Line, Col: lt.Col, Msg: "Unfinished expression. Expected `)` but found EOF."}
	}

	t := p.getNextToken()

	if t.Type != tt {
		return nil, &ScanError{File: t.File, Line: t.Line, Col: t.Col, Msg: fmt.Sprintf("Unfinished expression. Expected `)` but found `%s`.", t.Lexeme)}
	}

	return &t, nil
}

func (p *Parser) getNextToken() Token {
	p.current += 1

	return p.tokens[p.current-1]
}

func (p *Parser) peek(ts ...TokenType) bool {
	if p.isAtEnd() {
		return false
	}

	for _, t := range ts {
		if p.tokens[p.current].Type == t {
			return true
		}
	}

	return false
}

func (p *Parser) isAtEnd() bool {
	return p.current >= len(p.tokens)
}
