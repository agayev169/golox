package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/agayev169/golox"
)

func main() {
	log.SetLevel(log.InfoLevel)
	log.SetFormatter(&log.JSONFormatter{})

	args := os.Args
	if len(args) > 2 {
		log.Printf("Usage: %s [script]\n", args[0])
		os.Exit(64)
	} else if len(args) == 2 {
		runFile(args[1])
	} else if len(args) == 1 {
		runPrompt()
	}
}

func runFile(path string) {
	f, err := os.Open(path)
	fatal(err)
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			panic(err)
		}
	}(f)

	r := bufio.NewReader(f)

	interp := golox.NewInterpreter()

	_, err = run(r, interp)
	if err != nil {
		fatal(err)
	}
}

func runPrompt() {
	reader := bufio.NewReader(os.Stdin)

	interp := golox.NewInterpreter()

	for {
		fmt.Print("> ")
		line, err := reader.ReadString('\n')
		fatal(err)

		res, err := run(bufio.NewReader(strings.NewReader(line)), interp)
		if !warn(err) && res != nil {
			if _, ok := res.(golox.Nil); ok {
				fmt.Println("nil")
			} else {
				fmt.Println(res)
			}
		}
	}
}

func run(r *bufio.Reader, interp *golox.Interpreter) (interface{}, error) {
	var lerr *golox.LoxError
	bs, err := io.ReadAll(r)

	if err != nil {
		return nil, err
	}

	s := golox.NewScanner(bytes.NewReader(bs))

	tokens, err := s.ScanTokens()

	if err != nil {
		return nil, err
	}

	log.Info("Scanned the following tokens: ", tokens)

	p := golox.NewParser(tokens)

	stmts, lerr := p.Parse()
	if lerr != nil {
		return nil, lerr
	}

	log.Info("Parsed the following statements: ", stmts)

	resolver := golox.NewResolver(interp)
	if lerr = resolver.Resolve(stmts); lerr != nil {
		return nil, lerr
	}

	res, lerr := interp.Interpret(stmts)

	if lerr != nil {
		return nil, lerr
	}

	return res, nil
}

func fatal(err error) {
	if err == nil {
		return
	}

	fmt.Printf("Error happened: %s\n", err.Error())

	os.Exit(1)
}

func warn(err error) bool {
	if err == nil {
		return false
	}

	fmt.Printf("Error happened: %s\n", err.Error())

	return true
}
