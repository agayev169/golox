package golox

import (
	"time"
)

type Callable interface {
	GetArity() int
	Call(i *Interpreter, args []interface{}) (interface{}, *LoxError)
}

// Custom

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

func (f *LoxFunction) Call(i *Interpreter, args []interface{}) (interface{}, *LoxError) {
	env := NewEnv(f.closure)

	for i, param := range f.decl.Params {
		arg := args[i]
		if err := env.Define(param, arg); err != nil {
			return nil, err
		}
	}

	if err := i.ExecuteBlock(f.decl.Body, env); err != nil {
		return nil, err
	}

	return nil, nil
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
