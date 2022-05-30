package golox

import "fmt"

type LoxInstance struct {
	c *LoxClass
}

func NewLoxInstance(c *LoxClass) *LoxInstance {
	return &LoxInstance{c: c}
}

func (i *LoxInstance) String() string {
	return fmt.Sprintf("%s instance", i.c.Name)
}
