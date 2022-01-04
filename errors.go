package golox

import "fmt"

type LoxErrorNumber int

const (
	UnexpectedChar LoxErrorNumber = iota
	UnterminatedString
	UnfinishedExpression
)

type LoxError struct {
	File   string
	Line   int
	Col    int
	Msg    string
	Number LoxErrorNumber
}

func (e *LoxError) Error() string {
	return fmt.Sprintf("err #%d, %s:%d:%d: %s", e.Number, e.File, e.Line, e.Col, e.Msg)
}
