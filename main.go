package main

import (
	"fmt"
	"go-lisp/lexer"
	"slices"
)

func main() {
	fmt.Println("hello")
	fmt.Println(slices.Collect(lexer.Lex("123")))
}
