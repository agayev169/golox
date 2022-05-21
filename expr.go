package golox

type Expr interface {
	Accept(v ExprVisitor) (interface{}, *LoxError)
}

// ================ Assign ================

type Assign struct {
	Name  Token
	Value Expr
}

func (a *Assign) Accept(v ExprVisitor) (interface{}, *LoxError) {
	return v.AcceptAssignExpr(a)
}

// ================ Binary ================

type Binary struct {
	Left     Expr
	Operator Token
	Right    Expr
}

func (b *Binary) Accept(v ExprVisitor) (interface{}, *LoxError) {
	return v.AcceptBinaryExpr(b)
}

// ================ Grouping ================

type Grouping struct {
	Expr Expr
}

func (g *Grouping) Accept(v ExprVisitor) (interface{}, *LoxError) {
	return v.AcceptGroupingExpr(g)
}

// ================ Literal ================

type Literal struct {
	Value interface{}
}

func (l *Literal) Accept(v ExprVisitor) (interface{}, *LoxError) {
	return v.AcceptLiteralExpr(l)
}

// ================ Unary ================

type Unary struct {
	Operator Token
	Right    Expr
}

func (u *Unary) Accept(v ExprVisitor) (interface{}, *LoxError) {
	return v.AcceptUnaryExpr(u)
}

// ================ Call ================

type Call struct {
	Callee Expr
	Paren  Token
	Args   []Expr
}

func (c *Call) Accept(v ExprVisitor) (interface{}, *LoxError) {
	return v.AcceptCallExpr(c)
}

// ================ Variable ================

type Variable struct {
	Name Token
}

func (va *Variable) Accept(v ExprVisitor) (interface{}, *LoxError) {
	return v.AcceptVariableExpr(va)
}

// ================ Logical ================

type Logical struct {
	Left     Expr
	Operator Token
	Right    Expr
}

func (l *Logical) Accept(v ExprVisitor) (interface{}, *LoxError) {
	return v.AcceptLogicalExpr(l)
}

// ================ ExprVisitor ================

type ExprVisitor interface {
	AcceptAssignExpr(*Assign) (interface{}, *LoxError)
	AcceptBinaryExpr(*Binary) (interface{}, *LoxError)
	AcceptGroupingExpr(*Grouping) (interface{}, *LoxError)
	AcceptLiteralExpr(*Literal) (interface{}, *LoxError)
	AcceptUnaryExpr(*Unary) (interface{}, *LoxError)
	AcceptCallExpr(*Call) (interface{}, *LoxError)
	AcceptVariableExpr(*Variable) (interface{}, *LoxError)
	AcceptLogicalExpr(*Logical) (interface{}, *LoxError)
}
