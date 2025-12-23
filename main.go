package main

import (
	"encoding/json"
	"fmt"
	"strings"
)

// todo: use Lisp as a configuration language with evaluation (think feature flags)

func main() {
	fmt.Println("hello")

	var out map[string]any
	json.Unmarshal([]byte(`{"foo":"bar", "x":1,"asdf":{"xxx":true}}`), &out)

	fmt.Println(extract(out, "foo"))
	fmt.Println(extract(out, "x"))
	fmt.Println(extract(out, "asdf"))
	fmt.Println(extract(out, "asdf.xxx"))
	fmt.Println(extract(out, "asdf.xxa"))
	fmt.Println(extract(out, "asdfgfs"))
}

func extract(m map[string]any, query string) (any, bool) {
	parts := strings.Split(query, ".")

	var ptr *map[string]any = &m
	for i, p := range parts {
		if candidate, ok := (*ptr)[p]; ok {
			if i == len(parts)-1 {
				return candidate, true
			}

			next, ok := (candidate).(map[string]any)
			if !ok {
				return nil, false
			}
			ptr = &next
		} else {
			return nil, false
		}
	}
	return nil, false
}
