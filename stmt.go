package golox

import "fmt"

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

func (b *Block) String() string {
	return fmt.Sprintf("(Block): { stmts:  %v }", b.Stmts)
}

// ================ Expression ================

type Expression struct {
	Expr Expr
}

func (e *Expression) Accept(v StmtVisitor) (interface{}, *LoxError) {
	return v.AcceptExpressionStmt(e)
}

func (e *Expression) String() string {
	return fmt.Sprintf("(Expression): { expr:  %v }", e.Expr)
}

// ================ Print ================

type Print struct {
	Expr Expr
}

func (p *Print) Accept(v StmtVisitor) (interface{}, *LoxError) {
	return v.AcceptPrintStmt(p)
}

func (p *Print) String() string {
	return fmt.Sprintf("(Print): { expr:  %v }", p.Expr)
}

// ================ Var ================

type Var struct {
	Name        Token
	Initializer Expr
}

func (va *Var) Accept(v StmtVisitor) (interface{}, *LoxError) {
	return v.AcceptVarStmt(va)
}

func (va *Var) String() string {
	return fmt.Sprintf("(Var): { name:  %v; initializer:  %v }", va.Name, va.Initializer)
}

// ================ Class ================

type Class struct {
	Name    Token
	Methods []Func
}

func (c *Class) Accept(v StmtVisitor) (interface{}, *LoxError) {
	return v.AcceptClassStmt(c)
}

func (c *Class) String() string {
	return fmt.Sprintf("(Class): { name:  %v; methods:  %v }", c.Name, c.Methods)
}

// ================ Func ================

type Func struct {
	Name   Token
	Params []Token
	Body   []Stmt
}

func (f *Func) Accept(v StmtVisitor) (interface{}, *LoxError) {
	return v.AcceptFuncStmt(f)
}

func (f *Func) String() string {
	return fmt.Sprintf("(Func): { name:  %v; params:  %v; body:  %v }", f.Name, f.Params, f.Body)
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

func (i *If) String() string {
	return fmt.Sprintf("(If): { condition:  %v; body:  %v; elseBody:  %v }", i.Condition, i.Body, i.ElseBody)
}

// ================ While ================

type While struct {
	Condition Expr
	Body      Stmt
}

func (w *While) Accept(v StmtVisitor) (interface{}, *LoxError) {
	return v.AcceptWhileStmt(w)
}

func (w *While) String() string {
	return fmt.Sprintf("(While): { condition:  %v; body:  %v }", w.Condition, w.Body)
}

// ================ Return ================

type Return struct {
	Keyword Token
	Value   Expr
}

func (r *Return) Accept(v StmtVisitor) (interface{}, *LoxError) {
	return v.AcceptReturnStmt(r)
}

func (r *Return) String() string {
	return fmt.Sprintf("(Return): { keyword:  %v; value:  %v }", r.Keyword, r.Value)
}

// ================ StmtVisitor ================

type StmtVisitor interface {
	AcceptBlockStmt(*Block) (interface{}, *LoxError)
	AcceptExpressionStmt(*Expression) (interface{}, *LoxError)
	AcceptPrintStmt(*Print) (interface{}, *LoxError)
	AcceptVarStmt(*Var) (interface{}, *LoxError)
	AcceptClassStmt(*Class) (interface{}, *LoxError)
	AcceptFuncStmt(*Func) (interface{}, *LoxError)
	AcceptIfStmt(*If) (interface{}, *LoxError)
	AcceptWhileStmt(*While) (interface{}, *LoxError)
	AcceptReturnStmt(*Return) (interface{}, *LoxError)
}
