package golox

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

	return nil, genUndefVarError(name)
}

func (e *Env) Assign(name Token, value interface{}) *LoxError {
	if _, ok := e.vars[name.Lexeme]; !ok {
		return genUndefVarError(name)
	}

	e.vars[name.Lexeme] = value

	return nil
}