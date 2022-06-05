package golox

import "fmt"

type LoxInstance struct {
	c     *LoxClass
	props map[string]interface{}
}

func NewLoxInstance(c *LoxClass) *LoxInstance {
	return &LoxInstance{c: c, props: make(map[string]interface{})}
}

func (i *LoxInstance) Get(t Token) (interface{}, *LoxError) {
	p, ok := i.props[t.Lexeme]
	if !ok {
		if p, ok = i.c.GetMethod(t.Lexeme); !ok {
			return nil, genError(t, UndefinedProperty, fmt.Sprintf("Undefined property '%s'.", t.Lexeme))
		}
	}

	return p, nil
}

func (i *LoxInstance) Set(t Token, v interface{}) {
	i.props[t.Lexeme] = v
}

func (i *LoxInstance) String() string {
	return fmt.Sprintf("%s instance", i.c.Name)
}
