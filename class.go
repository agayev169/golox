package golox

type LoxClass struct {
	Name    string
	methods map[string]*LoxFunction
}

func NewLoxClass(name string, methods map[string]*LoxFunction) *LoxClass {
	return &LoxClass{Name: name, methods: methods}
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

func (c *LoxClass) GetMethod(name string) (*LoxFunction, bool) {
	m, ok := c.methods[name]
	return m, ok
}
