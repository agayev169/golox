package golox

import (
	"bytes"
	"io"
	"log"
	"strconv"
	"strings"
)

var keywords = map[string]TokenType{
	"and":    AND,
	"class":  CLASS,
	"else":   ELSE,
	"false":  FALSE,
	"for":    FOR,
	"fun":    FUN,
	"if":     IF,
	"nil":    NIL,
	"or":     OR,
	"print":  PRINT,
	"return": RETURN,
	"super":  SUPER,
	"this":   THIS,
	"true":   TRUE,
	"var":    VAR,
	"while":  WHILE,
}

type Scanner struct {
	source     *bytes.Reader
	tokens     []Token
	line       int
	col        int
	curTokenSb strings.Builder
}

func NewScanner(r *bytes.Reader) *Scanner {
	return &Scanner{source: r, tokens: make([]Token, 0), line: 0, col: 0, curTokenSb: strings.Builder{}}
}

func (s *Scanner) ScanTokens() ([]Token, error) {
	for !s.isAtEnd() {
		err := s.scanToken()
		if err != nil {
			return nil, err
		}
	}

	s.addToken(EOF, nil)
	return s.tokens, nil
}

func (s *Scanner) isAtEnd() bool {
	log.Printf("[TRACE] s.source.Len() == %d\n", s.source.Len())
	return s.source.Len() <= 0
}

func (s *Scanner) scanToken() error {
	s.curTokenSb.Reset()
	b := s.readNext()

	var typ TokenType = NONE
	var literal interface{} = nil

	switch b {
	case '(':
		typ = LEFT_PAREN
	case ')':
		typ = RIGHT_PAREN
	case '{':
		typ = LEFT_BRACE
	case '}':
		typ = RIGHT_BRACE
	case ',':
		typ = COMMA
	case '.':
		typ = DOT
	case '-':
		typ = MINUS
	case '+':
		typ = PLUS
	case ';':
		typ = SEMICOLON
	case '*':
		typ = STAR
	case '!':
		if s.match("=") {
			typ = BANG_EQUAL
		} else {
			typ = BANG
		}
	case '=':
		if s.match("=") {
			typ = EQUAL_EQUAL
		} else {
			typ = EQUAL
		}
	case '<':
		if s.match("=") {
			typ = LESS_EQUAL
		} else {
			typ = LESS
		}
	case '>':
		if s.match("=") {
			typ = GREATER_EQUAL
		} else {
			typ = GREATER
		}
	case '/':
		if s.match("/") {
			for !s.isAtEnd() {
				b := s.peek()
				log.Printf("[TRACE] b == %b, (b == 10) == %v\n", b, b == 10)

				if b == '\n' || s.isAtEnd() {
					break
				}

				s.readNext()
			}
		} else {
			typ = SLASH
		}
	case ' ':
	case '\r':
	case '\t':
	case '\n':
		// Ignore whitespace.
	case '"':
		s.parseString()
	default:
		if isDigit(b) {
			s.parseNumber()
		} else if isAlpha(b) {
			s.parseIdentifier()
		} else {
			return &LoxError{Number: UnexpectedChar, File: "", Line: s.line, Col: s.col, Msg: "Unexpected character."}
		}
	}

	if typ != NONE {
		s.addToken(typ, &literal)
	}

	return nil
}

func (s *Scanner) parseString() error {
	for !s.isAtEnd() {
		b := s.peek()

		if b == '"' {
			break
		}

		s.readNext()
	}

	if s.isAtEnd() {
		return &LoxError{Number: UnterminatedString, File: "", Line: s.line, Col: s.col, Msg: "Unterminated string. Expected \""}
	}

	s.readNext()

	val := s.curTokenSb.String()

	s.addToken(STRING, val[1:len(val)-1])

	return nil
}

func (s *Scanner) parseNumber() error {
	for !s.isAtEnd() {
		b := s.peek()

		if !isDigit(b) {
			break
		}

		s.readNext()
	}

	if !s.isAtEnd() && s.peek() == '.' && isDigit(s.peekNext()) {
		s.readNext()

		for !s.isAtEnd() {
			b := s.peek()

			if !isDigit(b) {
				break
			}

			s.readNext()
		}
	}

	val, err := strconv.ParseFloat(s.curTokenSb.String(), 64)
	fatal(err)

	s.addToken(NUMBER, val)

	return nil
}

func (s *Scanner) parseIdentifier() {
	for !s.isAtEnd() && isAlphaNumeric(s.peek()) {
		s.readNext()
	}

	text := s.curTokenSb.String()

	if t, ok := keywords[text]; ok {
		s.addToken(t, nil)
	} else {
		s.addToken(IDENTIFIER, nil)
	}
}

func (s *Scanner) readNext() byte {
	b, err := s.source.ReadByte()

	if err == io.EOF {
		err = nil
	} else {
		fatal(err)
	}

	s.curTokenSb.WriteByte(b)

	log.Printf("[TRACE] Read %s\n", string(b))

	s.col += 1

	if b == '\n' {
		s.col = 0
		s.line++
	}

	return b
}

func (s *Scanner) peek() byte {
	b, err := s.source.ReadByte()
	fatal(err)
	err = s.source.UnreadByte()
	fatal(err)
	return b
}

func (s *Scanner) peekNext() byte {
	if s.source.Len() < 2 {
		return 0
	}

	_, err := s.source.ReadByte()
	fatal(err)

	b, err := s.source.ReadByte()
	fatal(err)

	err = s.source.UnreadByte()
	fatal(err)

	err = s.source.UnreadByte()
	fatal(err)
	return b
}

func (s *Scanner) match(m string) bool {
	sb := strings.Builder{}

	line := s.line
	col := s.col

	for sb.Len() < len(m) && !s.isAtEnd() {
		b := s.readNext()

		sb.WriteByte(b)
	}

	if sb.String() != m {
		s.unread(sb.Len(), line, col)

		return false
	}

	return true
}

func (s *Scanner) unread(n, line, col int) {
	s.source.Seek(-int64(n), io.SeekCurrent)
	s.line = line
	s.col = col
}

func (s *Scanner) addToken(t TokenType, literal interface{}) {
	token := Token{Type: t, Lexeme: s.curTokenSb.String(), Literal: literal, File: "file.lox", Line: s.line, Col: s.col}

	log.Printf("[DEBUG] Adding token: %v\n", token)

	s.tokens = append(s.tokens, token)
	s.curTokenSb.Reset()
}
