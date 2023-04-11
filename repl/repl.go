package repl

import (
	"bufio"
	"fmt"
	"io"
	"lang/eval"
	"lang/lexer"
	"lang/object"
	"lang/parser"
	"lang/token"
)

var Stage map[string]int64 = map[string]int64{
	"lexer":  0,
	"parser": 1,
	"tree":   2,
	"eval":   3,
}

const PROMPT = ">> "

func Start(in io.Reader, out io.Writer, mode int64) {
	scanner := bufio.NewScanner(in)
	scope := object.NewScope()
	for {
		fmt.Printf(PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}
		line := scanner.Text()
		l := lexer.NewLexer(line)

		if mode == 0 {
			for tok := l.NextToken(); tok.Type != token.EOF; tok = l.NextToken() {
				fmt.Printf("%#v\n", tok)
			}
		}

		if mode == 1 || mode == 2 || mode == 3 {
			p := parser.NewParser(l)
			program := p.Parse()
			if len(p.Errors()) != 0 {
				printParserErrors(out, p.Errors())
				continue
			}

			if mode == 1 {
				io.WriteString(out, program.String())
				io.WriteString(out, "\n")
			}

			if mode == 2 {
				parser.DrawTree(program)
			}

			if mode == 3 {
				evaluated := eval.Eval(program, scope)
				io.WriteString(out, evaluated.Inspect())
				io.WriteString(out, "\n")
			}
		}
	}
}

func printParserErrors(out io.Writer, errors []string) {
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"\n")
	}
}
