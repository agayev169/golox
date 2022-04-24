package golox

import "fmt"

type Interpreter struct {
    env *Env
}

func NewInterpreter() *Interpreter {
    return &Interpreter{env: NewEnv(nil)}
}

func (interp *Interpreter) Interpret(stmts []Stmt) (interface{}, error) {
	for _, stmt := range stmts {
		_, err := stmt.Accept(interp)
		if err != nil {
			return nil, err
		}
	}

	return nil, nil
}

func (interp *Interpreter) AcceptExpressionStmt(expr *Expression) (interface{}, error) {
	_, err := interp.evaluate(expr.Expr)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (interp *Interpreter) AcceptPrintStmt(expr *Print) (interface{}, error) {
	val, err := interp.evaluate(expr.Expr)
	if err != nil {
		return nil, err
	}

	fmt.Println(val)

	return nil, nil
}

func (interp *Interpreter) AcceptBlockStmt(b *Block) (interface{}, error) {
    interp.env = NewEnv(interp.env);
    defer func() {
        interp.env = interp.env.enclosing
    }()

    for _, s := range b.Stmts {
        _, err := s.Accept(interp)
        if err != nil {
            return nil, err
        }
    }

    return nil, nil
}

func (interp *Interpreter) AcceptVarStmt(v *Var) (interface{}, error) {
    var init interface{} = nil

    if v.Initializer != nil {
        val, err := interp.evaluate(v.Initializer)
        if err != nil {
            return nil, err
        }

        init = val
    }

    interp.env.Define(v.Name.Lexeme, init)

    return nil, nil
}

func (interp *Interpreter) AcceptAssignExpr(a *Assign) (interface{}, error) {
    val, err := interp.evaluate(a.Value)
    if err != nil {
        return nil, err
    }

    err2 := interp.env.Assign(a.Name, val)
    if err2 != nil {
        return nil, err2
    }

    return nil, nil
}

func (interp *Interpreter) AcceptLiteralExpr(l *Literal) (interface{}, error) {
	return l.Value, nil
}

func (interp *Interpreter) AcceptGroupingExpr(g *Grouping) (interface{}, error) {
	return g.Expr.Accept(interp)
}

func (interp *Interpreter) AcceptUnaryExpr(u *Unary) (interface{}, error) {
	v, err := u.Right.Accept(interp)

	if err != nil {
		return nil, err
	}

	switch u.Operator.Type {
	case MINUS:
		err = interp.checkNumberOperand(u.Operator, v)

		if err != nil {
			return nil, err
		}

		return -(v.(float64)), nil
	case BANG:
		return !interp.isTruthy(v), nil
	}

	return nil, &LoxError{
		Number: UnexpectedChar, File: u.Operator.File, Line: u.Operator.Line, Col: u.Operator.Col,
		Msg: fmt.Sprintf("Unsupported operator type `%s`.", u.Operator.Lexeme),
	}
}

func (interp *Interpreter) AcceptBinaryExpr(b *Binary) (interface{}, error) {
	lv, err := b.Left.Accept(interp)

	if err != nil {
		return nil, err
	}

	rv, err := b.Right.Accept(interp)

	if err != nil {
		return nil, err
	}

	switch b.Operator.Type {
	case MINUS:
		err := interp.checkNumberOperands(b.Operator, lv, rv)

		if err != nil {
			return nil, err
		}

		return lv.(float64) - rv.(float64), nil
	case STAR:
		err := interp.checkNumberOperands(b.Operator, lv, rv)

		if err != nil {
			return nil, err
		}

		return lv.(float64) * rv.(float64), nil
	case SLASH:
		err := interp.checkNumberOperands(b.Operator, lv, rv)

		if err != nil {
			return nil, err
		}

		return lv.(float64) / rv.(float64), nil
	case PLUS:
		if lvv, ok := lv.(float64); ok {
			if rvv, ok := rv.(float64); ok {
				return lvv + rvv, nil
			}
		}

		if lvv, ok := lv.(string); ok {
			if rvv, ok := rv.(string); ok {
				return lvv + rvv, nil
			}
		}

		return nil, &LoxError{
			Number: UnexpectedChar, File: b.Operator.File, Line: b.Operator.Line, Col: b.Operator.Col,
			Msg: "Operands must be two numbers or two strings.",
		}
	case GREATER:
		err := interp.checkNumberOperands(b.Operator, lv, rv)

		if err != nil {
			return nil, err
		}

		return lv.(float64) > rv.(float64), nil
	case GREATER_EQUAL:
		err := interp.checkNumberOperands(b.Operator, lv, rv)

		if err != nil {
			return nil, err
		}

		return lv.(float64) >= rv.(float64), nil
	case LESS:
		err := interp.checkNumberOperands(b.Operator, lv, rv)

		if err != nil {
			return nil, err
		}

		return lv.(float64) < rv.(float64), nil
	case LESS_EQUAL:
		err := interp.checkNumberOperands(b.Operator, lv, rv)

		if err != nil {
			return nil, err
		}

		return lv.(float64) <= rv.(float64), nil
	case BANG_EQUAL:
		return !interp.isEqual(lv, rv), nil
	case EQUAL_EQUAL:
		return interp.isEqual(lv, rv), nil
	}

	return nil, &LoxError{
		Number: UnexpectedChar, File: b.Operator.File, Line: b.Operator.Line, Col: b.Operator.Col,
		Msg: fmt.Sprintf("Unsupported operator type `%s`.", b.Operator.Lexeme),
	}
}

func (interp *Interpreter) AcceptVariableExpr(v *Variable) (interface{}, error) {
    val, err := interp.env.Get(v.Name)
    if err != nil {
        return nil, err
    }

    return val, nil
}

func (interp *Interpreter) checkNumberOperand(op Token, r interface{}) error {
	if _, ok := r.(float64); !ok {
		return &LoxError{
			Number: UnexpectedChar, File: op.File, Line: op.Line, Col: op.Col,
			Msg: fmt.Sprintf("Expected a number but found `%s`.", op.Lexeme),
		}
	}

	return nil
}

func (interp *Interpreter) checkNumberOperands(op Token, l, r interface{}) error {
	_, ok1 := l.(float64)
	_, ok2 := r.(float64)

	if !ok1 || !ok2 {
		return &LoxError{
			Number: UnexpectedChar, File: op.File, Line: op.Line, Col: op.Col,
			Msg: "Operands must be numbers.",
		}
	}

	return nil
}

func (interp *Interpreter) isTruthy(v interface{}) bool {
	if v == nil {
		return false
	}

	switch v := v.(type) {
	case bool:
		return v
	default:
		return true
	}
}

func (interp *Interpreter) isEqual(lv, rv interface{}) bool {
	if lv == nil && rv == nil {
		return true
	}

	if lv == nil || rv == nil {
		return false
	}

	return lv == rv
}

func (interp *Interpreter) evaluate(expr Expr) (interface{}, error) {
	return expr.Accept(interp)
}
