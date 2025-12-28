package main

import (
	"fmt"
	"net/http"
)

func main() {
	fmt.Println("hello")
	code := `
(set port 8080)
(set host "localhost")
(lambda foo (param) (
    (if (= param 1) true false)
))
(lambda fibo (x)(
    (if (= x 0)
        0 
        (if (= x 1)
            1 
            (+ (fibo (- x 1)) (fibo (- x 2)))))	
))`
	srv, err := NewServer(code)
	if err != nil {
		fmt.Println("error starting server", err)
		return
	}
	http.ListenAndServe(":8080", srv)
}
