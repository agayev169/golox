package golox

type LoxClass struct {
	Name string
}

func NewLoxClass(name string) *LoxClass {
	return &LoxClass{Name: name}
}

func (c *LoxClass) String() string {
	return c.Name
}

func (c *LoxClass) GetArity() int {
	return 0
}

func (c *LoxClass) Call(*Interpreter, []interface{}) (interface{}, *LoxError) {
	return NewLoxInstance(c), nil
}
