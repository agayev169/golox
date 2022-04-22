package golox

import "fmt"

type Parser struct {
	tokens  []Token
	current int
}

func NewParser(tokens []Token) *Parser {
	return &Parser{tokens: tokens, current: 0}
}

func (p *Parser) Parse() ([]Stmt, *LoxError) {
	p.current = 0
	stmts := make([]Stmt, 0)

	for !p.isAtEnd() {
		stmt, err := p.parseStmt()
		if err != nil {
			return nil, err
		}

		stmts = append(stmts, stmt)
	}

	return stmts, nil
}

func (p *Parser) parseStmt() (Stmt, *LoxError) {
	if p.peek(PRINT) {
		return p.parsePrintStmt()
	}

	return p.parseExprStmt()
}

func (p *Parser) parsePrintStmt() (Stmt, *LoxError) {
	_, err := p.consume(PRINT)
	if err != nil {
		return nil, err
	}

	expr, err2 := p.parseExpression()
	if err2 != nil {
		return nil, err2
	}

	_, err3 := p.consume(SEMICOLON)
	if err3 != nil {
		return nil, err3
	}

	return &Print{Expr: expr}, nil
}

func (p *Parser) parseExprStmt() (Stmt, *LoxError) {
	expr, err := p.parseExpression()
	if err != nil {
		return nil, err
	}

	_, err2 := p.consume(SEMICOLON)
	if err2 != nil {
		return nil, err2
	}

	return &Expression{Expr: expr}, nil
}

func (p *Parser) parseExpression() (Expr, *LoxError) {
	return p.parseEquality()
}

func (p *Parser) parseEquality() (Expr, *LoxError) {
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

func (p *Parser) parseComparison() (Expr, *LoxError) {
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

func (p *Parser) parseTerm() (Expr, *LoxError) {
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

func (p *Parser) parseFactor() (Expr, *LoxError) {
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

func (p *Parser) parseUnary() (Expr, *LoxError) {
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

func (p *Parser) parsePrimary() (Expr, *LoxError) {
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

	return nil, &LoxError{Number: UnexpectedChar, File: t.File, Line: t.Line, Col: t.Col, Msg: fmt.Sprintf("Expected one of (number, string, `true`, `false`, `nil`, `(`) but found `%s`.", t.Lexeme)}
}

func (p *Parser) sync() {
	for !p.isAtEnd() {
		t := p.getNextToken()
		if t.Type == SEMICOLON {
			return
		}

		if p.peek(CLASS, FUN, VAR, FOR, IF, WHILE, PRINT, RETURN) {
			return
		}
	}
}

func (p *Parser) consume(tt TokenType) (*Token, *LoxError) {
	if p.isAtEnd() {
		lt := p.tokens[len(p.tokens)-1]
		return nil, &LoxError{Number: UnfinishedExpression, File: lt.File, Line: lt.Line, Col: lt.Col, Msg: fmt.Sprintf("Unfinished expression. Expected `%s` but found EOF.", tt)}
	}

	t := p.getNextToken()

	if t.Type != tt {
		return nil, &LoxError{Number: UnfinishedExpression, File: t.File, Line: t.Line, Col: t.Col, Msg: fmt.Sprintf("Unfinished expression. Expected `%s` but found `%s`.", tt, t.Lexeme)}
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
    return p.current >= len(p.tokens) || p.tokens[p.current].Type == EOF
}
