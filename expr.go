package golox

import "fmt"

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

func (a *Assign) String() string {
	return fmt.Sprintf("(Assign): {name:  %v; value:  %v}", a.Name, a.Value)
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

func (b *Binary) String() string {
	return fmt.Sprintf("(Binary): {left:  %v; operator:  %v; right:  %v}", b.Left, b.Operator, b.Right)
}

// ================ Grouping ================

type Grouping struct {
	Expr Expr
}

func (g *Grouping) Accept(v ExprVisitor) (interface{}, *LoxError) {
	return v.AcceptGroupingExpr(g)
}

func (g *Grouping) String() string {
	return fmt.Sprintf("(Grouping): {expr:  %v}", g.Expr)
}

// ================ Literal ================

type Literal struct {
	Value interface{}
}

func (l *Literal) Accept(v ExprVisitor) (interface{}, *LoxError) {
	return v.AcceptLiteralExpr(l)
}

func (l *Literal) String() string {
	return fmt.Sprintf("(Literal): {value:  %v}", l.Value)
}

// ================ Unary ================

type Unary struct {
	Operator Token
	Right    Expr
}

func (u *Unary) Accept(v ExprVisitor) (interface{}, *LoxError) {
	return v.AcceptUnaryExpr(u)
}

func (u *Unary) String() string {
	return fmt.Sprintf("(Unary): {operator:  %v; right:  %v}", u.Operator, u.Right)
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

func (c *Call) String() string {
	return fmt.Sprintf("(Call): {callee:  %v; paren:  %v; args:  %v}", c.Callee, c.Paren, c.Args)
}

// ================ Get ================

type Get struct {
	Obj  Expr
	Name Token
}

func (g *Get) Accept(v ExprVisitor) (interface{}, *LoxError) {
	return v.AcceptGetExpr(g)
}

func (g *Get) String() string {
	return fmt.Sprintf("(Get): {obj:  %v; name:  %v}", g.Obj, g.Name)
}

// ================ Set ================

type Set struct {
	Obj   Expr
	Name  Token
	Value Expr
}

func (s *Set) Accept(v ExprVisitor) (interface{}, *LoxError) {
	return v.AcceptSetExpr(s)
}

func (s *Set) String() string {
	return fmt.Sprintf("(Set): {obj:  %v; name:  %v; value:  %v}", s.Obj, s.Name, s.Value)
}

// ================ This ================

type This struct {
	Token Token
}

func (t *This) Accept(v ExprVisitor) (interface{}, *LoxError) {
	return v.AcceptThisExpr(t)
}

func (t *This) String() string {
	return fmt.Sprintf("(This): {token:  %v}", t.Token)
}

// ================ Variable ================

type Variable struct {
	Name Token
}

func (va *Variable) Accept(v ExprVisitor) (interface{}, *LoxError) {
	return v.AcceptVariableExpr(va)
}

func (va *Variable) String() string {
	return fmt.Sprintf("(Variable): {name:  %v}", va.Name)
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

func (l *Logical) String() string {
	return fmt.Sprintf("(Logical): {left:  %v; operator:  %v; right:  %v}", l.Left, l.Operator, l.Right)
}

// ================ ExprVisitor ================

type ExprVisitor interface {
	AcceptAssignExpr(*Assign) (interface{}, *LoxError)
	AcceptBinaryExpr(*Binary) (interface{}, *LoxError)
	AcceptGroupingExpr(*Grouping) (interface{}, *LoxError)
	AcceptLiteralExpr(*Literal) (interface{}, *LoxError)
	AcceptUnaryExpr(*Unary) (interface{}, *LoxError)
	AcceptCallExpr(*Call) (interface{}, *LoxError)
	AcceptGetExpr(*Get) (interface{}, *LoxError)
	AcceptSetExpr(*Set) (interface{}, *LoxError)
	AcceptThisExpr(*This) (interface{}, *LoxError)
	AcceptVariableExpr(*Variable) (interface{}, *LoxError)
	AcceptLogicalExpr(*Logical) (interface{}, *LoxError)
}
