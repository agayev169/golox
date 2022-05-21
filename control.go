package golox

type ControlType int

const (
	Ret ControlType = iota
)

type Control struct {
	Type ControlType
	Val  interface{}
}

func NewReturn(val interface{}) *Control {
	return &Control{
		Type: Ret,
		Val:  val,
	}
}
