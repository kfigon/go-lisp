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

		exp          models.SExpression
		envAssertion func(*testing.T, *models.Env)
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
			exp:  models.Number(123 + (1 + 2) + (3 - 4)),
		},
		{
			desc: "variable declaration",
			code: `(set x 123)`,
			exp:  models.Nil{},
			envAssertion: func(t *testing.T, e *models.Env) {
				got, ok := e.Get("x")
				assert.True(t, ok, "x not found")
				val, err := got()
				assert.NoError(t, err)
				assert.Equal(t, models.Number(123), val)
			},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			ast, err := parser.Parse(lexer.Lex(tC.code))
			assert.NoError(t, err)

			e := NewEvaluator()
			got, err := e.Eval(ast)
			assert.NoError(t, err)
			assert.Equal(t, tC.exp, got)

			if tC.envAssertion != nil {
				tC.envAssertion(t, e.RootEnv)
			}
		})
	}
}
