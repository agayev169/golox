package golox

type Stmt interface {
	Accept(v StmtVisitor) (interface{}, error)
}

// ================ Expression ================

type Expression struct {
	Expr Expr
}

func (e *Expression) Accept(v StmtVisitor) (interface{}, error) {
	return v.AcceptExpressionStmt(e)
}

// ================ Print ================

type Print struct {
	Expr Expr
}

func (p *Print) Accept(v StmtVisitor) (interface{}, error) {
	return v.AcceptPrintStmt(p)
}

// ================ Var ================

type Var struct {
	Name        Token
	Initializer Expr
}

func (va *Var) Accept(v StmtVisitor) (interface{}, error) {
	return v.AcceptVarStmt(va)
}

// ================ StmtVisitor ================

type StmtVisitor interface {
	AcceptExpressionStmt(*Expression) (interface{}, error)
	AcceptPrintStmt(*Print) (interface{}, error)
	AcceptVarStmt(*Var) (interface{}, error)
}
