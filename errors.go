package golox

import "fmt"

type LoxErrorNumber int

const (
	UnexpectedChar LoxErrorNumber = iota
	UnterminatedString
	UnfinishedExpression
	UndefinedVariable
	UnassignedVariable
	InvalidAssignment
	ArgumentLimitExceeded
	InvalidCall
	InvalidArity
	ParamLimitExceeded
	InvalidParamName
	NameAlreadyDefined
	ReturnOutsideFunc
	SelfInitialization
)

var errorNames = map[LoxErrorNumber]string{
	UnexpectedChar:        "Unexpected character",
	UnterminatedString:    "Unterminated string",
	UnfinishedExpression:  "Unfinished expression",
	UndefinedVariable:     "Undefined variable",
	UnassignedVariable:    "Unassigned variable",
	InvalidAssignment:     "Invalid assignment",
	ArgumentLimitExceeded: "Argument limit exceeded",
	InvalidCall:           "Invalid call",
	InvalidArity:          "Invalid arity",
	ParamLimitExceeded:    "Parameter limit exceeded",
	InvalidParamName:      "Invalid parameter name",
	NameAlreadyDefined:    "Name already defined",
	ReturnOutsideFunc:     "Return outside function",
	SelfInitialization:    "Variable self initialization",
}

type LoxError struct {
	File   string
	Line   int
	Col    int
	Msg    string
	Number LoxErrorNumber
}

func (e *LoxError) Error() string {
	return fmt.Sprintf("ERR '%s': %s:%d:%d: %s", errorNames[e.Number], e.File, e.Line, e.Col, e.Msg)
}

func genUndefVarError(t Token) *LoxError {
	return genError(t, UndefinedVariable, fmt.Sprintf("Undefined variable %s", t.Lexeme))
}

func genError(t Token, num LoxErrorNumber, msg string) *LoxError {
	return &LoxError{
		File:   t.File,
		Line:   t.Line,
		Col:    t.Col,
		Number: num,
		Msg:    msg,
	}
}
