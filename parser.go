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
		stmt, err := p.parseDeclaration()
		if err != nil {
			return nil, err
		}

		stmts = append(stmts, stmt)
	}

	return stmts, nil
}

func (p *Parser) parseDeclaration() (Stmt, *LoxError) {
	if p.peek(VAR) {
		s, err := p.parseVarDeclaration()
		if err != nil {
			p.sync()

			return nil, err
		}

		return s, nil
	}

	return p.parseStmt()
}

func (p *Parser) parseVarDeclaration() (Stmt, *LoxError) {
	_, err := p.consume(VAR)
	if err != nil {
		return nil, err
	}

	name, err2 := p.consume(IDENTIFIER)
	if err2 != nil {
		return nil, err2
	}

	var initializer Expr

	if p.peek(EQUAL) {
		_, err3 := p.consume(EQUAL)
		if err3 != nil {
			return nil, err3
		}

		init, err4 := p.parseExpression()
		if err4 != nil {
			return nil, err4
		}

		initializer = init
	}

	_, err3 := p.consume(SEMICOLON)
	if err3 != nil {
		return nil, err3
	}

	return &Var{Name: *name, Initializer: initializer}, nil
}

func (p *Parser) parseStmt() (Stmt, *LoxError) {
	if p.peek(PRINT) {
		return p.parsePrintStmt()
	} else if p.peek(LEFT_BRACE) {
        stmts, err := p.parseBlock()
        if err != nil {
            return nil, err
        }

        return &Block{Stmts: stmts}, nil
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

func (p *Parser) parseBlock() ([]Stmt, *LoxError) {
    _, err := p.consume(LEFT_BRACE)
    if err != nil {
        return nil, err
    }

    stmts := make([]Stmt, 0)

    for !p.isAtEnd() && !p.peek(RIGHT_BRACE) {
        stmt, err2 := p.parseDeclaration()
        if err2 != nil {
            return nil, err2
        }

        stmts = append(stmts, stmt)
    }

    _, err2 := p.consume(RIGHT_BRACE)
    if err2 != nil {
        return nil, err2
    }

    return stmts, nil
}

func (p *Parser) parseExpression() (Expr, *LoxError) {
	return p.parseAssignment()
}

func (p *Parser) parseAssignment() (Expr, *LoxError) {
    expr, err := p.parseEquality()
    if err != nil {
        return nil, err
    }

    if p.peek(EQUAL) {
        equals, err2 := p.consume(EQUAL)
        if err2 != nil {
            return nil, err2
        }

        v, ok := expr.(*Variable)
        if !ok {
            return nil, &LoxError{File: equals.File, Line: equals.Line, Col: equals.Col, Number: InvalidAssignment, Msg: "Invalid assignment target."}
        }

        assignment, err3 := p.parseAssignment()
        if err3 != nil {
            return nil, err3
        }

        return &Assign{Name: v.Name, Value: assignment}, nil
    }
    
    return expr, nil
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
		return &Literal{Value: Nil{}}, nil
	}

	if t.Type == IDENTIFIER {
		return &Variable{Name: t}, nil
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

	return nil, &LoxError{Number: UnexpectedChar, File: t.File, Line: t.Line, Col: t.Col, Msg: fmt.Sprintf("Expected one of (number, string, `true`, `false`, `nil`, identifier, `(`}) but found `%s`.", t.Lexeme)}
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
