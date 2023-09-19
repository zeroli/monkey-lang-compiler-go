package main

import (
	"fmt"
	"monkey/evaluator"
	"monkey/lexer"
	"monkey/object"
	"monkey/parser"
	"monkey/repl"
	"os"
	"os/user"
)

func main() {
	user, err := user.Current()
	if err != nil {
		panic(err)
	}
	if len(os.Args) > 1 {
		line := os.Args[1]
		l := lexer.New(line)
		p := parser.New(l)
		program := p.ParseProgram()
		fmt.Println(program.String())

		env := object.NewEnvironment()
		result := evaluator.Eval(program, env)
		if result != nil {
			fmt.Println("=>", result.Inspect())
		} else {
			fmt.Println("=>NIL")
		}
	} else {
		fmt.Printf("Hello %s! This is the monkey programming language!\n",
			user.Username)
		fmt.Printf("Feel free to type in commands\n")
		repl.Start(os.Stdin, os.Stdout)
	}
}
