package golox

import "fmt"

type Env struct {
	vars      map[string]interface{}
	enclosing *Env
}

func NewEnv(enclosing *Env) *Env {
	return &Env{vars: make(map[string]interface{}), enclosing: enclosing}
}

func (e *Env) Define(name Token, val interface{}) *LoxError {
	if _, ok := e.vars[name.Lexeme]; ok {
		return genError(name, NameAlreadyDefined, fmt.Sprintf("Cannot redefine '%s'\n", name.Lexeme))
	}

	e.vars[name.Lexeme] = val

	return nil
}

func (e *Env) Get(name Token) (interface{}, *LoxError) {
	if val, ok := e.vars[name.Lexeme]; ok {
		return val, nil
	}

	if e.enclosing == nil {
		return nil, genUndefVarError(name)
	}

	return e.enclosing.Get(name)
}

func (e *Env) Assign(name Token, value interface{}) *LoxError {
	if _, ok := e.vars[name.Lexeme]; !ok {
		if e.enclosing == nil {
			return genUndefVarError(name)
		}

		return e.enclosing.Assign(name, value)
	}

	e.vars[name.Lexeme] = value

	return nil
}
