package golox

type Expr interface {
	Accept(v ExprVisitor) (interface{}, error)
}

// ================ Binary ================

type Binary struct {
	Left     Expr
	Operator Token
	Right    Expr
}

func (b *Binary) Accept(v ExprVisitor) (interface{}, error) {
	return v.AcceptBinaryExpr(b)
}

// ================ Grouping ================

type Grouping struct {
	Expr Expr
}

func (g *Grouping) Accept(v ExprVisitor) (interface{}, error) {
	return v.AcceptGroupingExpr(g)
}

// ================ Literal ================

type Literal struct {
	Value interface{}
}

func (l *Literal) Accept(v ExprVisitor) (interface{}, error) {
	return v.AcceptLiteralExpr(l)
}

// ================ Unary ================

type Unary struct {
	Operator Token
	Right    Expr
}

func (u *Unary) Accept(v ExprVisitor) (interface{}, error) {
	return v.AcceptUnaryExpr(u)
}

// ================ Variable ================

type Variable struct {
	Name Token
}

func (va *Variable) Accept(v ExprVisitor) (interface{}, error) {
	return v.AcceptVariableExpr(va)
}

// ================ ExprVisitor ================

type ExprVisitor interface {
	AcceptBinaryExpr(*Binary) (interface{}, error)
	AcceptGroupingExpr(*Grouping) (interface{}, error)
	AcceptLiteralExpr(*Literal) (interface{}, error)
	AcceptUnaryExpr(*Unary) (interface{}, error)
	AcceptVariableExpr(*Variable) (interface{}, error)
}
