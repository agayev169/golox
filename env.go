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

func (e *Env) DefineThis(val interface{}) {
	e.vars["this"] = val
}

func (e *Env) GetAt(name Token, depth int) (interface{}, *LoxError) {
	env := e.ancestor(depth)

	return env.Get(name)
}

func (e *Env) Get(name Token) (interface{}, *LoxError) {
	if val, ok := e.vars[name.Lexeme]; ok {
		return val, nil
	}

	return nil, genUndefVarError(name)
}

func (e *Env) AssignAt(name Token, value interface{}, depth int) *LoxError {
	env := e.ancestor(depth)

	return env.Assign(name, value)
}

func (e *Env) Assign(name Token, value interface{}) *LoxError {
	if _, ok := e.vars[name.Lexeme]; !ok {
		return genUndefVarError(name)
	}

	e.vars[name.Lexeme] = value

	return nil
}

func (e *Env) ancestor(depth int) *Env {
	res := e

	for i := 0; i < depth; i++ {
		if res.enclosing == nil {
			panic(fmt.Sprintf("Ancestor with the depth %d does not exist since the depth of the env is %d.", depth, i))
		}

		res = res.enclosing
	}

	return res
}
