package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	log "github.com/sirupsen/logrus"
	"os"
	"strings"

	"github.com/agayev169/golox"
)

func main() {
    log.SetLevel(log.WarnLevel)

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
	defer f.Close()

	r := bufio.NewReader(f)

	interp := &golox.Interpreter{}

	err = run(r, interp)
	if err != nil {
		fatal(err)
	}
}

func runPrompt() {
	reader := bufio.NewReader(os.Stdin)

	interp := &golox.Interpreter{}

	for {
		fmt.Print("> ")
		line, err := reader.ReadString('\n')
		fatal(err)

		err = run(bufio.NewReader(strings.NewReader(line)), interp)
		warn(err)
	}
}

func run(r *bufio.Reader, interp *golox.Interpreter) error {
	bs, err2 := io.ReadAll(r)

	if err2 != nil {
		return err2
	}

	s := golox.NewScanner(bytes.NewReader(bs))

	tokens, err2 := s.ScanTokens()

	if err2 != nil {
		return err2
	}

	log.Info("Scanned the following tokens: ", tokens)

	p := golox.NewParser(tokens)

	expr, err3 := p.Parse()

	if err3 != nil {
		return err3
	}

	str, err4 := interp.Interpret(expr)

	if err4 != nil {
		return err4
	}

	fmt.Println(str)

	return nil
}

func fatal(err error) {
	if err == nil {
		return
	}

	log.WithField("error", err).Error("Error happened")
	os.Exit(1)
}

func warn(err error) {
	if err == nil {
		return
	}

	log.WithField("error", err).Warn("Warning")
}
