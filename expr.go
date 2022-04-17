package golox

type Expr interface {
	Accept(v Visitor) (interface{}, error)
}

// ================ Binary ================

type Binary struct {
	Left     Expr
	Operator Token
	Right    Expr
}

func (b *Binary) Accept(v Visitor) (interface{}, error) {
	return v.AcceptBinaryExpr(b)
}

// ================ Grouping ================

type Grouping struct {
	Expr Expr
}

func (g *Grouping) Accept(v Visitor) (interface{}, error) {
	return v.AcceptGroupingExpr(g)
}

// ================ Literal ================

type Literal struct {
	Value interface{}
}

func (l *Literal) Accept(v Visitor) (interface{}, error) {
	return v.AcceptLiteralExpr(l)
}

// ================ Unary ================

type Unary struct {
	Operator Token
	Right    Expr
}

func (u *Unary) Accept(v Visitor) (interface{}, error) {
	return v.AcceptUnaryExpr(u)
}

// ================ Visitor ================

type Visitor interface {
	AcceptBinaryExpr(*Binary) (interface{}, error)
	AcceptGroupingExpr(*Grouping) (interface{}, error)
	AcceptLiteralExpr(*Literal) (interface{}, error)
	AcceptUnaryExpr(*Unary) (interface{}, error)
}
