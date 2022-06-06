package golox

import (
	"fmt"
	"time"
)

type Callable interface {
	GetArity() int
	Call(i *Interpreter, args []interface{}) (interface{}, *LoxError)
}

// Custom

type FunctionType = int

const (
	None FunctionType = iota
	Function
	Method
)

type LoxFunction struct {
	decl    *Func
	closure *Env
}

func NewLoxFunction(decl *Func, closure *Env) *LoxFunction {
	return &LoxFunction{decl: decl, closure: closure}
}

func (f *LoxFunction) GetArity() int {
	return len(f.decl.Params)
}

func (f *LoxFunction) Bind(i *LoxInstance) *LoxFunction {
	env := NewEnv(f.closure)
	env.DefineThis(i)

	return NewLoxFunction(f.decl, env)
}

func (f *LoxFunction) Call(i *Interpreter, args []interface{}) (interface{}, *LoxError) {
	return f.call(i, args)
}

func (f *LoxFunction) call(i *Interpreter, args []interface{}) (ret interface{}, err *LoxError) {
	defer func() {
		if r := recover(); r != nil {
			c, ok := r.(*Control)
			if !ok || c.Type != Ret {
				panic(r)
			}

			ret = c.Val
		}
	}()

	ret, err = nil, nil

	env := NewEnv(f.closure)

	for i, param := range f.decl.Params {
		arg := args[i]
		if err = env.Define(param, arg); err != nil {
			return
		}
	}

	if err = i.ExecuteBlock(f.decl.Body, env); err != nil {
		return
	}

	return
}

func (f *LoxFunction) String() string {
	return fmt.Sprintf("<fn %s>", f.decl.Name.Lexeme)
}

// Clock

type LoxClock struct{}

func (l *LoxClock) GetArity() int {
	return 0
}

func (l *LoxClock) Call(*Interpreter, []interface{}) (interface{}, *LoxError) {
	return float64(time.Now().UnixMilli()) / 1000.0, nil
}

func (l *LoxClock) String() string {
	return "<native fn>"
}
