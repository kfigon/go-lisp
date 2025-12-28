package main

import (
	"fmt"
	"net/http"
)

func main() {
	fmt.Println("hello")
	srv, err := NewServer(`(set port 8080)
	(set host "localhost")
	(lambda foo (param) (
		(if (= param 1) true false)
	))`)
	if err != nil {
		fmt.Println("error starting server", err)
		return
	}
	http.ListenAndServe(":8080", srv)
}
