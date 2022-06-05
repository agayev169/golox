package golox

import "fmt"

type Resolver struct {
	scopes      []map[string]bool
	interpreter *Interpreter
	curf        FunctionType
}

func NewResolver(i *Interpreter) *Resolver {
	return &Resolver{scopes: make([]map[string]bool, 0), interpreter: i, curf: None}
}

func (r *Resolver) Resolve(stmts []Stmt) *LoxError {
	for _, stmt := range stmts {
		if err := r.resolveStmt(stmt); err != nil {
			return err
		}
	}

	return nil
}

func (r *Resolver) AcceptAssignExpr(a *Assign) (interface{}, *LoxError) {
	err := r.resolveExpr(a.Value)
	if err != nil {
		return nil, err
	}

	err = r.resolveLocalExpr(a, a.Name)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (r *Resolver) AcceptBinaryExpr(b *Binary) (interface{}, *LoxError) {
	if err := r.resolveExpr(b.Left); err != nil {
		return nil, err
	}

	if err := r.resolveExpr(b.Right); err != nil {
		return nil, err
	}

	return nil, nil
}

func (r *Resolver) AcceptGroupingExpr(gr *Grouping) (interface{}, *LoxError) {
	return nil, r.resolveExpr(gr.Expr)
}

func (r *Resolver) AcceptLiteralExpr(*Literal) (interface{}, *LoxError) {
	return nil, nil
}

func (r *Resolver) AcceptUnaryExpr(u *Unary) (interface{}, *LoxError) {
	return nil, r.resolveExpr(u.Right)
}

func (r *Resolver) AcceptCallExpr(c *Call) (interface{}, *LoxError) {
	if err := r.resolveExpr(c.Callee); err != nil {
		return nil, err
	}

	for _, arg := range c.Args {
		if err := r.resolveExpr(arg); err != nil {
			return nil, err
		}
	}

	return nil, nil
}

func (r *Resolver) AcceptGetExpr(g *Get) (interface{}, *LoxError) {
	return nil, r.resolveExpr(g.Obj)
}

func (r *Resolver) AcceptSetExpr(s *Set) (interface{}, *LoxError) {
	if err := r.resolveExpr(s.Obj); err != nil {
		return nil, err
	}

	return nil, r.resolveExpr(s.Value)
}

func (r *Resolver) AcceptVariableExpr(v *Variable) (interface{}, *LoxError) {
	if len(r.scopes) != 0 {
		if def, ok := r.scopes[len(r.scopes)-1][v.Name.Lexeme]; ok && !def {
			return nil, genError(v.Name, SelfInitialization, "Can't read local variable in its own initializer.")
		}
	}

	return nil, r.resolveLocalExpr(v, v.Name)
}

func (r *Resolver) AcceptLogicalExpr(l *Logical) (interface{}, *LoxError) {
	if err := r.resolveExpr(l.Left); err != nil {
		return nil, err
	}

	if err := r.resolveExpr(l.Right); err != nil {
		return nil, err
	}

	return nil, nil
}

func (r *Resolver) AcceptBlockStmt(b *Block) (interface{}, *LoxError) {
	r.beginScope()
	defer r.endScope()

	return nil, r.resolveBlock(b.Stmts)
}

func (r *Resolver) AcceptExpressionStmt(e *Expression) (interface{}, *LoxError) {
	return nil, r.resolveExpr(e.Expr)
}

func (r *Resolver) AcceptPrintStmt(p *Print) (interface{}, *LoxError) {
	return nil, r.resolveExpr(p.Expr)
}

func (r *Resolver) AcceptVarStmt(v *Var) (interface{}, *LoxError) {
	if err := r.declare(v.Name); err != nil {
		return nil, err
	}

	if v.Initializer != nil {
		if err := r.resolveExpr(v.Initializer); err != nil {
			return nil, err
		}
	}

	r.define(v.Name)

	return nil, nil
}

func (r *Resolver) AcceptFuncStmt(f *Func) (interface{}, *LoxError) {
	if err := r.declare(f.Name); err != nil {
		return nil, err
	}

	r.define(f.Name)

	oldf := r.curf
	r.curf = Function
	defer func() {
		r.curf = oldf
	}()

	r.beginScope()
	defer r.endScope()

	for _, p := range f.Params {
		if err := r.declare(p); err != nil {
			return nil, err
		}

		r.define(p)
	}

	return nil, r.resolveBlock(f.Body)
}

func (r *Resolver) AcceptClassStmt(c *Class) (interface{}, *LoxError) {
	if err := r.declare(c.Name); err != nil {
		return nil, err
	}

	r.define(c.Name)

	return nil, nil
}

func (r *Resolver) AcceptIfStmt(i *If) (interface{}, *LoxError) {
	if err := r.resolveExpr(i.Condition); err != nil {
		return nil, err
	}

	if err := r.resolveStmt(i.Body); err != nil {
		return nil, err
	}

	return nil, r.resolveStmt(i.ElseBody)
}

func (r *Resolver) AcceptWhileStmt(w *While) (interface{}, *LoxError) {
	if err := r.resolveExpr(w.Condition); err != nil {
		return nil, err
	}

	return nil, r.resolveStmt(w.Body)
}

func (r *Resolver) AcceptReturnStmt(ret *Return) (interface{}, *LoxError) {
	if r.curf == None {
		return nil, genError(ret.Keyword, ReturnOutsideFunc, "return statement cannot be used outside function.")
	}

	if ret.Value != nil {
		return nil, r.resolveExpr(ret.Value)
	}

	return nil, nil
}

func (r *Resolver) resolveExpr(expr Expr) *LoxError {
	_, err := expr.Accept(r)

	return err
}

func (r *Resolver) resolveStmt(stmt Stmt) *LoxError {
	_, err := stmt.Accept(r)

	return err
}

func (r *Resolver) resolveLocalExpr(expr Expr, name Token) *LoxError {
	for i := 0; i < len(r.scopes); i++ {
		if _, ok := r.scopes[len(r.scopes)-i-1][name.Lexeme]; ok {
			r.interpreter.Resolve(expr, i)

			return nil
		}
	}

	return nil
}

func (r *Resolver) declare(name Token) *LoxError {
	if len(r.scopes) == 0 {
		return nil
	}

	sc := r.scopes[len(r.scopes)-1]

	if _, ok := sc[name.Lexeme]; ok {
		return genError(name, NameAlreadyDefined, fmt.Sprintf("Cannot redefine '%s'\n", name.Lexeme))
	}

	sc[name.Lexeme] = false

	return nil
}

func (r *Resolver) define(name Token) {
	if len(r.scopes) == 0 {
		return
	}

	r.scopes[len(r.scopes)-1][name.Lexeme] = true
}

func (r *Resolver) resolveBlock(stmts []Stmt) *LoxError {
	for _, stmt := range stmts {
		if err := r.resolveStmt(stmt); err != nil {
			return err
		}
	}
	return nil
}

func (r *Resolver) beginScope() {
	r.scopes = append(r.scopes, make(map[string]bool))
}

func (r *Resolver) endScope() {
	if len(r.scopes) < 1 {
		panic("Cannot end a scope since none has been started.")
	}

	r.scopes = r.scopes[:len(r.scopes)-1]
}
