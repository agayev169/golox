package golox

type Expr interface {
	Accept(v Visitor) interface{}
}

// ================ Binary ================

type Binary struct {
	Left     Expr
	Operator Token
	Right    Expr
}

func (b *Binary) Accept(v Visitor) interface{} {
	return v.AcceptBinaryExpr(b)
}

// ================ Grouping ================

type Grouping struct {
	Expr Expr
}

func (g *Grouping) Accept(v Visitor) interface{} {
	return v.AcceptGroupingExpr(g)
}

// ================ Literal ================

type Literal struct {
	Value interface{}
}

func (l *Literal) Accept(v Visitor) interface{} {
	return v.AcceptLiteralExpr(l)
}

// ================ Unary ================

type Unary struct {
	Operator Token
	Right    Expr
}

func (u *Unary) Accept(v Visitor) interface{} {
	return v.AcceptUnaryExpr(u)
}

// ================ Visitor ================

type Visitor interface {
	AcceptBinaryExpr(*Binary) interface{}
	AcceptGroupingExpr(*Grouping) interface{}
	AcceptLiteralExpr(*Literal) interface{}
	AcceptUnaryExpr(*Unary) interface{}
}
