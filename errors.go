package golox

import "fmt"

type LoxErrorNumber int

const (
	UnexpectedChar LoxErrorNumber = iota
	UnterminatedString
	UnfinishedExpression
	RuntimeError
	UndefinedVariable
    UnassignedVariable
	InvalidAssignment
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

func genUndefVarError(t Token) *LoxError {
	return &LoxError{
		File:   t.File,
		Line:   t.Line,
		Col:    t.Col,
		Number: UndefinedVariable,
		Msg:    fmt.Sprintf("Undefined variable %s", t.Lexeme),
	}
}
