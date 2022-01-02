package golox

type Expr interface {
	Accept(v *Visitor)
}

// ================ Binary ================

type Binary struct {
	Left     Expr
	Operator Token
	Right    Expr
}

func (b *Binary) Accept(v Visitor) {
	v.AcceptBinaryExpr(b)
}

// ================ Grouping ================

type Grouping struct {
	Expr Expr
}

func (g *Grouping) Accept(v Visitor) {
	v.AcceptGroupingExpr(g)
}

// ================ Literal ================

type Literal struct {
	Value interface{}
}

func (l *Literal) Accept(v Visitor) {
	v.AcceptLiteralExpr(l)
}

// ================ Unary ================

type Unary struct {
	Operator Token
	Right    Expr
}

func (u *Unary) Accept(v Visitor) {
	v.AcceptUnaryExpr(u)
}

// ================ Visitor ================

type Visitor interface {
	AcceptBinaryExpr(*Binary)
	AcceptGroupingExpr(*Grouping)
	AcceptLiteralExpr(*Literal)
	AcceptUnaryExpr(*Unary)
}
