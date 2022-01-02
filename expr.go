package golox

type Expr interface {
   Accept(v *Visitor)
}

// ================ Binary ================

type Binary struct {
    left Expr
    operator Token
    right Expr
}

func (b *Binary) Accept(v Visitor) {
    v.AcceptBinaryExpr(b)
}

// ================ Grouping ================

type Grouping struct {
    expr Expr
}

func (g *Grouping) Accept(v Visitor) {
    v.AcceptGroupingExpr(g)
}

// ================ Literal ================

type Literal struct {
    value interface{}
}

func (l *Literal) Accept(v Visitor) {
    v.AcceptLiteralExpr(l)
}

// ================ Unary ================

type Unary struct {
    operator Token
    right Expr
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