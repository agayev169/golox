package golox

type Stmt interface {
	Accept(v StmtVisitor) (interface{}, *LoxError)
}

// ================ Block ================

type Block struct {
	Stmts []Stmt
}

func (b *Block) Accept(v StmtVisitor) (interface{}, *LoxError) {
	return v.AcceptBlockStmt(b)
}

// ================ Expression ================

type Expression struct {
	Expr Expr
}

func (e *Expression) Accept(v StmtVisitor) (interface{}, *LoxError) {
	return v.AcceptExpressionStmt(e)
}

// ================ Print ================

type Print struct {
	Expr Expr
}

func (p *Print) Accept(v StmtVisitor) (interface{}, *LoxError) {
	return v.AcceptPrintStmt(p)
}

// ================ Var ================

type Var struct {
	Name        Token
	Initializer Expr
}

func (va *Var) Accept(v StmtVisitor) (interface{}, *LoxError) {
	return v.AcceptVarStmt(va)
}

// ================ If ================

type If struct {
	Condition Expr
	Body      Stmt
	ElseBody  Stmt
}

func (i *If) Accept(v StmtVisitor) (interface{}, *LoxError) {
	return v.AcceptIfStmt(i)
}

// ================ While ================

type While struct {
	Condition Expr
	Body      Stmt
}

func (w *While) Accept(v StmtVisitor) (interface{}, *LoxError) {
	return v.AcceptWhileStmt(w)
}

// ================ StmtVisitor ================

type StmtVisitor interface {
	AcceptBlockStmt(*Block) (interface{}, *LoxError)
	AcceptExpressionStmt(*Expression) (interface{}, *LoxError)
	AcceptPrintStmt(*Print) (interface{}, *LoxError)
	AcceptVarStmt(*Var) (interface{}, *LoxError)
	AcceptIfStmt(*If) (interface{}, *LoxError)
	AcceptWhileStmt(*While) (interface{}, *LoxError)
}
