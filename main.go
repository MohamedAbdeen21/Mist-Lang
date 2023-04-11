package main

import (
	"fmt"
	"lang/eval"
	"lang/lexer"
	"lang/object"
	"lang/parser"
	"os"
	"strconv"
	"strings"
)

func printError(code string, msg string) {
	split := strings.SplitAfterN(msg, " ", 2)
	pos := strings.Split(split[0], ",")
	row, _ := strconv.Atoi(strings.Split(pos[0], "[")[1])
	col, _ := strconv.Atoi(strings.Split(pos[1], "]")[0])

	println()
	rowNumber := strconv.Itoa(row) + ": "
	println(rowNumber, strings.Split(code, "\n")[row-1])
	println(strings.Repeat(" ", col+len(rowNumber)) + ("^ ") + split[1])
}

func main() {
	file := os.Args[1]
	bytes, err := os.ReadFile(file)
	if err != nil {
		fmt.Printf("file not found %s\n", file)
		return
	}

	code := string(bytes) + "\nmain();"
	l := lexer.NewLexer(string(code))
	// emit tokens, need to comment everyting else

	// for l.NextToken().Type != token.EOF {
	// 	fmt.Printf("%#v\n", l.NextToken())
	// }
	p := parser.NewParser(l)
	program := p.Parse()
	if len(p.Errors()) != 0 {
		for _, err := range p.Errors() {
			printError(code, err)
			return
		}
	}

	// uncomment to see AST, redirect to file if tree is too wide
	// parser.DrawTree(program)

	// prevents output if an error occurs, since this is currently an interpreter
	// not a compiler
	f := eval.SetupStdout()
	defer f.Close()

	evaluated := eval.Eval(program, object.NewScope())
	if evaluated.Type() == object.ERROR_OBJ {
		printError(code, evaluated.Inspect())
		return
	}
	stdout, _ := os.ReadFile(f.Name())
	print(string(stdout))
	print(evaluated.Inspect())
}
