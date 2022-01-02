package golox

import "fmt"

type ScanError struct {
	File string
	Line int
	Col  int
	Msg  string
}

func (e *ScanError) Error() string {
	return fmt.Sprintf("%s:%d:%d: %s", e.File, e.Line, e.Col, e.Msg)
}
