package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/agayev169/golox"
)

func main() {
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

	err = run(r)
	if err != nil {
		fatal(err)
	}
}

func runPrompt() {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("> ")
		line, err := reader.ReadString('\n')
		fatal(err)

		err = run(bufio.NewReader(strings.NewReader(line)))
		warn(err)
	}
}

func run(r *bufio.Reader) error {
	bs, err := io.ReadAll(r)
	fatal(err)

	s := golox.NewScanner(bytes.NewReader(bs))

	tokens, err := s.ScanTokens()
	fatal(err)

	log.Printf("[INFO] Scanned the following tokens: %v\n", tokens)

	return nil
}

func fatal(err error) {
	if err == nil {
		return
	}

	log.Printf("[ERR] %s\n", err)
	os.Exit(1)
}

func warn(err error) {
	if err == nil {
		return
	}

	log.Printf("[WARN]: %s\n", err)
}
