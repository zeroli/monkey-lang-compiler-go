package main

import (
	"fmt"
	"monkey/compiler"
	"monkey/lexer"
	"monkey/parser"
	"monkey/repl"
	"monkey/vm"
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

		compiler := compiler.New()
		err := compiler.Compile(program)
		if err != nil {
			fmt.Println("=>NIL")
		}

		machine := vm.New(compiler.Bytecode())
		err = machine.Run()
		if err != nil {
			fmt.Println("=>NIL")
		}
		stackTop := machine.StackTop()
		fmt.Println("=>", stackTop.Inspect())
	} else {
		fmt.Printf("Hello %s! This is the monkey programming language!\n",
			user.Username)
		fmt.Printf("Feel free to type in commands\n")
		repl.Start(os.Stdin, os.Stdout)
	}
}
