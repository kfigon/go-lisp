package eval

import (
	"go-lisp/lexer"
	"go-lisp/models"
	"go-lisp/parser"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEval(t *testing.T) {
	testCases := []struct {
		desc string
		code string

		exp models.SExpression
	}{
		{
			desc: "basic num",
			code: `123`,
			exp:  models.Number(123),
		},
		{
			desc: "sum",
			code: `(+ 123 1 2)`,
			exp:  models.Number(126),
		},
		{
			desc: "couple of lists",
			code: `(+ 123 1 2)
					(* 2 3)`,
			exp: models.Number(6),
		},
		{
			desc: "nested",
			code: `(+ 123 (+ 1 2) (- 3 4))`,
			exp:  models.Number(1),
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			ast, err := parser.Parse(lexer.Lex(tC.code))
			assert.NoError(t, err)

			got, err := Eval(ast)
			assert.NoError(t, err)
			assert.Equal(t, tC.exp, got)
		})
	}
}
