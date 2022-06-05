package golox

import "fmt"

type Interpreter struct {
	globEnv *Env
	env     *Env
	locals  map[Expr]int
}

func NewInterpreter() *Interpreter {
	env := NewEnv(nil)

	if err := addFunc(env, Token{Lexeme: "clock", Type: IDENTIFIER}, &LoxClock{}); err != nil {
		panic(err)
	}

	return &Interpreter{globEnv: env, env: env, locals: make(map[Expr]int)}
}

func addFunc(env *Env, name Token, f Callable) *LoxError {
	return env.Define(name, f)
}

func (interp *Interpreter) Interpret(stmts []Stmt) (interface{}, *LoxError) {
	var res interface{}

	for _, stmt := range stmts {
		v, err := stmt.Accept(interp)
		if err != nil {
			return nil, err
		}

		res = v
	}

	return res, nil
}

func (interp *Interpreter) AcceptClassStmt(c *Class) (interface{}, *LoxError) {
	if err := interp.env.Define(c.Name, nil); err != nil {
		return nil, err
	}

	cl := NewLoxClass(c.Name.Lexeme)

	return nil, interp.env.Assign(c.Name, cl)
}

func (interp *Interpreter) AcceptFuncStmt(f *Func) (interface{}, *LoxError) {
	if err := addFunc(interp.env, f.Name, NewLoxFunction(f, interp.env)); err != nil {
		return nil, err
	}

	return nil, nil
}

func (interp *Interpreter) AcceptExpressionStmt(expr *Expression) (interface{}, *LoxError) {
	res, err := interp.evaluate(expr.Expr)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (interp *Interpreter) AcceptPrintStmt(expr *Print) (interface{}, *LoxError) {
	val, err := interp.evaluate(expr.Expr)
	if err != nil {
		return nil, err
	}

	if _, ok := val.(Nil); ok || val == nil {
		fmt.Println("nil")
	} else {
		fmt.Println(val)
	}

	return nil, nil
}

func (interp *Interpreter) AcceptBlockStmt(b *Block) (interface{}, *LoxError) {
	err := interp.ExecuteBlock(b.Stmts, NewEnv(interp.env))
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (interp *Interpreter) ExecuteBlock(ss []Stmt, env *Env) *LoxError {
	oldEnv := interp.env
	interp.env = env
	defer func() {
		interp.env = oldEnv
	}()

	for _, s := range ss {
		_, err := s.Accept(interp)
		if err != nil {
			return err
		}
	}

	return nil
}

func (interp *Interpreter) AcceptIfStmt(iff *If) (interface{}, *LoxError) {
	cond, err := interp.evaluate(iff.Condition)
	if err != nil {
		return nil, err
	}

	if interp.isTruthy(cond) {
		if res, err2 := interp.execute(iff.Body); err2 != nil {
			return nil, err2
		} else {
			return res, nil
		}
	}

	if iff.ElseBody != nil {
		return interp.execute(iff.ElseBody)
	}

	return nil, nil
}

func (interp *Interpreter) AcceptWhileStmt(w *While) (interface{}, *LoxError) {
	for {
		cond, err := interp.evaluate(w.Condition)
		if err != nil {
			return nil, err
		}

		if !interp.isTruthy(cond) {
			break
		}

		_, err2 := interp.execute(w.Body)
		if err2 != nil {
			return nil, err2
		}
	}

	return nil, nil
}

func (interp *Interpreter) AcceptVarStmt(v *Var) (interface{}, *LoxError) {
	var init interface{} = nil

	if v.Initializer != nil {
		val, err := interp.evaluate(v.Initializer)
		if err != nil {
			return nil, err
		}

		init = val
	}

	if err := interp.env.Define(v.Name, init); err != nil {
		return nil, err
	}

	return nil, nil
}

func (interp *Interpreter) AcceptReturnStmt(r *Return) (interface{}, *LoxError) {
	var ret interface{} = nil

	if r.Value != nil {
		r, err := interp.evaluate(r.Value)
		if err != nil {
			return nil, err
		}

		ret = r
	}

	panic(NewReturn(ret))
}

func (interp *Interpreter) AcceptAssignExpr(a *Assign) (interface{}, *LoxError) {
	val, err := interp.evaluate(a.Value)
	if err != nil {
		return nil, err
	}

	if d, ok := interp.locals[a]; ok {
		if err = interp.env.AssignAt(a.Name, val, d); err != nil {
			return nil, err
		}
	} else {
		if err = interp.globEnv.Assign(a.Name, val); err != nil {
			return nil, err
		}
	}

	return nil, nil
}

func (interp *Interpreter) AcceptLiteralExpr(l *Literal) (interface{}, *LoxError) {
	return l.Value, nil
}

func (interp *Interpreter) AcceptGroupingExpr(g *Grouping) (interface{}, *LoxError) {
	return g.Expr.Accept(interp)
}

func (interp *Interpreter) AcceptUnaryExpr(u *Unary) (interface{}, *LoxError) {
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

func (interp *Interpreter) AcceptCallExpr(c *Call) (interface{}, *LoxError) {
	callee, err := interp.evaluate(c.Callee)
	if err != nil {
		return nil, err
	}

	cf, ok := callee.(Callable)
	if !ok {
		return nil, genError(c.Paren, InvalidCall, "Can only call functions and classes.")
	}

	if cf.GetArity() != len(c.Args) {
		return nil, genError(c.Paren, InvalidArity,
			fmt.Sprintf("Expected %d arguments but got %d.", cf.GetArity(), len(c.Args)))
	}

	args := make([]interface{}, 0, len(c.Args))
	for _, arg := range c.Args {
		a, err2 := interp.evaluate(arg)
		if err2 != nil {
			return nil, err2
		}

		args = append(args, a)
	}

	return cf.Call(interp, args)
}

func (interp *Interpreter) AcceptGetExpr(g *Get) (interface{}, *LoxError) {
	obj, err := interp.evaluate(g.Obj)
	if err != nil {
		return nil, err
	}

	i, ok := obj.(*LoxInstance)

	if !ok {
		return nil, genError(g.Name, NonInstanceProperty, "Only instances have properties.")
	}

	return i.Get(g.Name)
}

func (interp *Interpreter) AcceptSetExpr(s *Set) (interface{}, *LoxError) {
	obj, err := interp.evaluate(s.Obj)
	if err != nil {
		return nil, err
	}

	i, ok := obj.(*LoxInstance)

	if !ok {
		return nil, genError(s.Name, NonInstanceProperty, "Only instances have properties.")
	}

	r, err := interp.evaluate(s.Value)
	if err != nil {
		return nil, err
	}

	i.Set(s.Name, r)

	return nil, nil
}

func (interp *Interpreter) AcceptBinaryExpr(b *Binary) (interface{}, *LoxError) {
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

func (interp *Interpreter) AcceptVariableExpr(v *Variable) (interface{}, *LoxError) {
	var val interface{}
	var err *LoxError

	if d, ok := interp.locals[v]; ok {
		val, err = interp.env.GetAt(v.Name, d)
		if err != nil {
			return nil, err
		}
	} else {
		val, err = interp.globEnv.Get(v.Name)
	}

	if val == nil {
		return nil,
			&LoxError{File: v.Name.File,
				Line:   v.Name.Line,
				Col:    v.Name.Col,
				Number: UnassignedVariable,
				Msg:    fmt.Sprintf("Usage of unassigned variable %s", v.Name.Lexeme)}
	}

	return val, nil
}

func (interp *Interpreter) AcceptLogicalExpr(l *Logical) (interface{}, *LoxError) {
	left, err := interp.evaluate(l.Left)
	if err != nil {
		return nil, err
	}

	if interp.isTruthy(left) && l.Operator.Type == OR {
		return left, nil
	} else if !interp.isTruthy(left) && l.Operator.Type == AND {
		return left, nil
	}

	return interp.evaluate(l.Right)
}

func (interp *Interpreter) Resolve(expr Expr, depth int) {
	interp.locals[expr] = depth
}

func (interp *Interpreter) checkNumberOperand(op Token, r interface{}) *LoxError {
	if _, ok := r.(float64); !ok {
		return &LoxError{
			Number: UnexpectedChar, File: op.File, Line: op.Line, Col: op.Col,
			Msg: fmt.Sprintf("Expected a number but found `%s`.", op.Lexeme),
		}
	}

	return nil
}

func (interp *Interpreter) checkNumberOperands(op Token, l, r interface{}) *LoxError {
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
	if _, ok := v.(Nil); ok {
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

func (interp *Interpreter) evaluate(expr Expr) (interface{}, *LoxError) {
	return expr.Accept(interp)
}

func (interp *Interpreter) execute(stmt Stmt) (interface{}, *LoxError) {
	return stmt.Accept(interp)
}
