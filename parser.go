package golox

import (
	"fmt"
)

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
	if p.peek(FUN) {
		s, err := p.parseFunDeclaration()
		if err != nil {
			return nil, err
		}

		return s, nil
	} else if p.peek(VAR) {
		s, err := p.parseVarDeclaration()
		if err != nil {
			p.sync()

			return nil, err
		}

		return s, nil
	} else if p.peek(CLASS) {
		return p.parseClassDeclaration()
	}

	return p.parseStmt()
}

func (p *Parser) parseClassDeclaration() (Stmt, *LoxError) {
	if _, err := p.consume(CLASS); err != nil {
		return nil, err
	}

	name, err := p.consume(IDENTIFIER)
	if err != nil {
		return nil, err
	}

	if _, err := p.consume(LEFT_BRACE); err != nil {
		return nil, err
	}

	ms := make([]Func, 0)

	for !p.peek(RIGHT_BRACE) {
		f, err := p.parseFunction()
		if err != nil {
			return nil, err
		}

		ms = append(ms, *f)
	}

	if _, err := p.consume(RIGHT_BRACE); err != nil {
		return nil, err
	}

	return &Class{
		Name:    *name,
		Methods: ms,
	}, nil
}

func (p *Parser) parseFunDeclaration() (Stmt, *LoxError) {
	if _, err := p.consume(FUN); err != nil {
		return nil, err
	}

	return p.parseFunction()
}

func (p *Parser) parseFunction() (*Func, *LoxError) {
	name, err := p.consume(IDENTIFIER)
	if err != nil {
		return nil, err
	}

	_, err = p.consume(LEFT_PAREN)
	if err != nil {
		return nil, err
	}

	params, err := p.parseParams()
	if err != nil {
		return nil, err
	}

	_, err = p.consume(RIGHT_PAREN)
	if err != nil {
		return nil, err
	}

	b, err := p.parseBlock()
	if err != nil {
		return nil, err
	}

	return &Func{Name: *name, Params: params, Body: b}, nil
}

func (p *Parser) parseParams() ([]Token, *LoxError) {
	res := make([]Token, 0)

	firstArg := true

	for !p.peek(RIGHT_PAREN) {
		if !firstArg {
			if _, err := p.consume(COMMA); err != nil {
				return nil, err
			}
		}
		t := p.getNextToken()

		if len(res) >= 255 {
			return nil, genError(t, ParamLimitExceeded, "Can't have more than 255 parameters.")
		}

		if t.Type != IDENTIFIER {
			return nil, genError(t, InvalidParamName, fmt.Sprintf("Expect parameter Name, got %s.", t.Lexeme))
		}

		res = append(res, t)

		firstArg = false
	}

	return res, nil
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
	} else if p.peek(IF) {
		return p.parseIfStmt()
	} else if p.peek(WHILE) {
		return p.parseWhileStmt()
	} else if p.peek(FOR) {
		return p.parseForStmt()
	} else if p.peek(RETURN) {
		return p.parseReturnStmt()
	}

	return p.parseExprStmt()
}

func (p *Parser) parseReturnStmt() (Stmt, *LoxError) {
	ret, err := p.consume(RETURN)
	if err != nil {
		return nil, err
	}

	var val Expr
	if !p.peek(SEMICOLON) {
		val, err = p.parseExpression()
		if err != nil {
			return nil, err
		}
	}

	if _, err = p.consume(SEMICOLON); err != nil {
		return nil, err
	}

	return &Return{Keyword: *ret, Value: val}, nil
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

func (p *Parser) parseIfStmt() (Stmt, *LoxError) {
	if _, err := p.consume(IF); err != nil {
		return nil, err
	}

	if _, err := p.consume(LEFT_PAREN); err != nil {
		return nil, err
	}

	expr, err := p.parseExpression()
	if err != nil {
		return nil, err
	}

	if _, err2 := p.consume(RIGHT_PAREN); err2 != nil {
		return nil, err2
	}

	body, err2 := p.parseStmt()
	if err2 != nil {
		return nil, err2
	}

	var elseBody Stmt
	if p.peek(ELSE) {
		if _, err3 := p.consume(ELSE); err3 != nil {
			return nil, err3
		}

		stmt, err3 := p.parseStmt()
		if err3 != nil {
			return nil, err3
		}

		elseBody = stmt
	}

	return &If{Condition: expr, Body: body, ElseBody: elseBody}, nil
}

func (p *Parser) parseWhileStmt() (Stmt, *LoxError) {
	if _, err := p.consume(WHILE); err != nil {
		return nil, err
	}

	if _, err := p.consume(LEFT_PAREN); err != nil {
		return nil, err
	}

	expr, err := p.parseExpression()
	if err != nil {
		return nil, err
	}

	if _, err2 := p.consume(RIGHT_PAREN); err2 != nil {
		return nil, err2
	}

	body, err2 := p.parseStmt()
	if err2 != nil {
		return nil, err2
	}

	return &While{Condition: expr, Body: body}, nil
}

func (p *Parser) parseForStmt() (Stmt, *LoxError) {
	if _, err := p.consume(FOR); err != nil {
		return nil, err
	}

	if _, err := p.consume(LEFT_PAREN); err != nil {
		return nil, err
	}

	init, err := p.parseForInit()
	if err != nil {
		return nil, err
	}

	cond, err2 := p.parseForCond()
	if err2 != nil {
		return nil, err2
	}

	increment, err3 := p.parseForIncrement()
	if err3 != nil {
		return nil, err3
	}

	body, err4 := p.parseStmt()
	if err4 != nil {
		return nil, err4
	}

	if increment != nil {
		body = &Block{Stmts: []Stmt{body, &Expression{Expr: increment}}}
	}

	if cond != nil {
		body = &While{Condition: cond, Body: body}
	}

	if init != nil {
		body = &Block{Stmts: []Stmt{init, body}}
	}

	return body, nil
}

func (p *Parser) parseForInit() (Stmt, *LoxError) {
	var init Stmt
	if p.peek(SEMICOLON) {
		if _, err := p.consume(SEMICOLON); err != nil {
			return nil, err
		}

		init = nil
	} else if p.peek(VAR) {
		if s, err := p.parseVarDeclaration(); err != nil {
			return nil, err
		} else {
			init = s
		}
	} else {
		if s, err := p.parseExprStmt(); err != nil {
			return nil, err
		} else {
			init = s
		}
	}

	return init, nil
}

func (p *Parser) parseForCond() (Expr, *LoxError) {
	var cond Expr
	if !p.peek(SEMICOLON) {
		if e, err := p.parseExpression(); err != nil {
			return nil, err
		} else {
			cond = e
		}
	}

	if _, err := p.consume(SEMICOLON); err != nil {
		return nil, err
	}

	return cond, nil
}

func (p *Parser) parseForIncrement() (Expr, *LoxError) {
	var increment Expr
	if !p.peek(RIGHT_PAREN) {
		if e, err := p.parseExpression(); err != nil {
			return nil, err
		} else {
			increment = e
		}
	}

	if _, err := p.consume(RIGHT_PAREN); err != nil {
		return nil, err
	}

	return increment, nil
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
	expr, err := p.parseOr()
	if err != nil {
		return nil, err
	}

	if p.peek(EQUAL) {
		equals, err2 := p.consume(EQUAL)
		if err2 != nil {
			return nil, err2
		}

		assignment, err3 := p.parseAssignment()
		if err3 != nil {
			return nil, err3
		}

		if v, ok := expr.(*Variable); ok {
			return &Assign{Name: v.Name, Value: assignment}, nil
		} else if g, ok := expr.(*Get); ok {
			return &Set{Obj: g.Obj, Name: g.Name, Value: assignment}, nil
		} else {
			return nil, genError(*equals, InvalidAssignment, "Invalid assignment target.")
		}
	}

	return expr, nil
}

func (p *Parser) parseOr() (Expr, *LoxError) {
	left, err := p.parseAnd()
	if err != nil {
		return nil, err
	}

	if !p.peek(OR) {
		return left, nil
	}

	op, err2 := p.consume(OR)
	if err2 != nil {
		return nil, err2
	}

	right, err3 := p.parseAnd()
	if err3 != nil {
		return nil, err3
	}

	return &Logical{Left: left, Operator: *op, Right: right}, nil
}

func (p *Parser) parseAnd() (Expr, *LoxError) {
	left, err := p.parseEquality()
	if err != nil {
		return nil, err
	}

	if !p.peek(AND) {
		return left, nil
	}

	op, err2 := p.consume(AND)
	if err2 != nil {
		return nil, err2
	}

	right, err3 := p.parseEquality()
	if err3 != nil {
		return nil, err3
	}

	return &Logical{Left: left, Operator: *op, Right: right}, nil
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

	return p.parseCall()
}

func (p *Parser) parseCall() (Expr, *LoxError) {
	pr, err := p.parsePrimary()
	if err != nil {
		return nil, err
	}

	res := pr

	for true {
		if p.peek(LEFT_PAREN) {
			lp, err2 := p.consume(LEFT_PAREN)
			if err2 != nil {
				return nil, err2
			}

			args, err3 := p.parseArguments()
			if err3 != nil {
				return nil, err3
			}

			_, err4 := p.consume(RIGHT_PAREN)
			if err4 != nil {
				return nil, err4
			}

			res = &Call{Callee: res, Paren: *lp, Args: args}
		} else if p.peek(DOT) {
			if _, err := p.consume(DOT); err != nil {
				return nil, err
			}

			name, err := p.consume(IDENTIFIER)
			if err != nil {
				return nil, err
			}

			res = &Get{Obj: res, Name: *name}
		} else {
			break
		}
	}

	return res, nil
}

func (p *Parser) parseArguments() ([]Expr, *LoxError) {
	res := make([]Expr, 0)

	firstArg := true

	for !p.peek(RIGHT_PAREN) {
		if !firstArg {
			if _, err := p.consume(COMMA); err != nil {
				return nil, err
			}
		}

		if len(res) >= 255 {
			return nil, genError(p.getNextToken(), ArgumentLimitExceeded, "Can't have more than 255 arguments.")
		}

		expr, err := p.parseExpression()
		if err != nil {
			return nil, err
		}

		res = append(res, expr)

		firstArg = false
	}

	return res, nil
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
		return nil, &LoxError{Number: UnfinishedExpression, File: lt.File, Line: lt.Line, Col: lt.Col, Msg: fmt.Sprintf("Unfinished expression. Expected `%s` but found EOF.", tokenNames[tt])}
	}

	t := p.getNextToken()

	if t.Type != tt {
		return nil, &LoxError{Number: UnfinishedExpression, File: t.File, Line: t.Line, Col: t.Col, Msg: fmt.Sprintf("Unfinished expression. Expected `%s` but found `%s`.", tokenNames[tt], t.Lexeme)}
	}

	return &t, nil
}

func (p *Parser) getPrevToken() Token {
	return p.tokens[p.current-1]
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
