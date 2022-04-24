package golox

import (
	"fmt"
)

type Env struct {
	vars map[string]interface{}
}

func NewEnv() *Env {
	return &Env{vars: make(map[string]interface{})}
}

func (e *Env) Define(name string, val interface{}) *LoxError {
	e.vars[name] = val

	return nil
}

func (e *Env) Get(name Token) (interface{}, *LoxError) {
	if val, ok := e.vars[name.Lexeme]; ok {
		return val, nil
	}

	return nil,
		&LoxError{
			File:   name.File,
			Line:   name.Line,
			Col:    name.Col,
			Number: UndefinedVariable,
			Msg:    fmt.Sprintf("Undefined variable %s", name.Lexeme),
		}
}
