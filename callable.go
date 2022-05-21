package golox

import (
	"time"
)

type Callable interface {
	GetArity() int
	Call(i *Interpreter, args []interface{}) (interface{}, error)
}

// Clock

type LoxClock struct{}

func (l LoxClock) GetArity() int {
	return 0
}

func (l LoxClock) Call(*Interpreter, []interface{}) (interface{}, error) {
	return float64(time.Now().UnixMilli()) / 1000.0, nil
}

func (l LoxClock) String() string {
	return "<native fn>"
}
