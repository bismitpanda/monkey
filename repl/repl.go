package repl

import (
	"bufio"
	"fmt"
	"io"
	"monkey/evaluator"
	"monkey/lexer"
	"monkey/object"
	"monkey/parser"
)

const PROMPT = ">> "
const MONKEY_FACE = `           __,__
   .--.  .-"   "-.  .--.
 /  .. \/ .-. .-. \/ ..  \
| |   '| /   Y   \ |'   | |
| \    \ \ 0 | 0 / /    / |
 \ '- ,\.-"""""""-./, -' /
  ''-' /_   ^ ^   _\ '-''
      |  \._   _./  |
       \  \ '~' /  /
        '._'-=-'_.'
          '-----'
`

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	env := object.NewEnvironment()

	for {
		fmt.Fprint(out, PROMPT)
		if scanned := scanner.Scan(); !scanned {
			return
		}

		line := scanner.Text()
		l := lexer.New(line)
		p := parser.New(l)

		program := p.ParseProgram()
		if len(p.Errors()) != 0 {
			printParseErrors(out, p.Errors())
			continue
		}

		evaluated := evaluator.Eval(program, env)
		if evaluated != nil {
			fmt.Fprintln(out, evaluated.Inspect())
		}
	}
}

func printParseErrors(out io.Writer, errors []string) {
	fmt.Fprint(out, MONKEY_FACE)
	fmt.Fprintln(out, "Woops! We ran into some monkey business here!")
	fmt.Fprintln(out, " parse errors:")
	for _, msg := range errors {
		fmt.Fprintln(out, "\t", msg)
	}
}
