package golox

import (
	"fmt"
	"log"
)

func fatal(err error) {
	if err == nil {
		return
	}

	// log.Printf("[ERR] %s\n", err)
	panic(fmt.Sprintf("[ERR] %s\n", err))
}

func warn(err error) {
	if err == nil {
		return
	}

	log.Printf("[WARN]: %s\n", err)
}

func isDigit(b byte) bool {
	return b >= '0' && b <= '9'
}

func isAlpha(b byte) bool {
	return (b >= 'a' && b <= 'z') ||
		(b >= 'A' && b <= 'Z') ||
		b == '_'
}

func isAlphaNumeric(b byte) bool {
	return (b >= 'a' && b <= 'z') ||
		(b >= 'A' && b <= 'Z') ||
		(b >= '0' && b <= '9') ||
		b == '_'
}
