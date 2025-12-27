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
			desc: "basic if",
			code: `(if true 1 2)`,
			exp:  models.Number(1),
		},
		{
			desc: "basic if2",
			code: `(if false 1 2)`,
			exp:  models.Number(2),
		},
		{
			desc: "if1",
			code: `(if (= 1 2) "foo" (+ 1 2))`,
			exp:  models.Number(3),
		},
		{
			desc: "if2",
			code: `(if (and (= 1 1) true) "foo" (+ 1 2))`,
			exp:  models.String("foo"),
		},
		{
			desc: "logic",
			code: `(and 
						(or 
							(or (= 1 2) true)
							true)
						(and true false))`,
			exp: models.Bool(false),
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
		{
			desc: "function declaration",
			code: `(lambda foobar (x y)(
				(+ 5 x y)))`,
			exp: models.Nil{},
			envAssertion: func(t *testing.T, e *models.Env) {
				_, ok := e.Get("foobar")
				assert.True(t, ok)

				_, ok = e.Get("x")
				assert.False(t, ok)

				_, ok = e.Get("y")
				assert.False(t, ok)
			},
		},
		{
			desc: "function invocation",
			code: `(lambda foobar (x)(
				(+ 5 x)))
				(foobar 10)`,
			exp: models.Number(15),
			envAssertion: func(t *testing.T, e *models.Env) {
				_, ok := e.Get("foobar")
				assert.True(t, ok)

				_, ok = e.Get("x")
				assert.False(t, ok)
			},
		},
		{
			desc: "function invocation2",
			code: `(lambda incr (x)((+ 1 x)))
				(incr (- 10 1))`,
			exp: models.Number(10),
		},
		{
			desc: "fibonacci",
			code: `
			(lambda fibo (x)(
				(if (= x 0)
					0 
					(if (= x 1)
						1 
						(+ (fibo (- x 1)) (fibo (- x 2)))))	
			))

			(fibo 10)`,
			exp: models.Number(55),
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			ast, err := parser.Parse(lexer.Lex(tC.code))
			assert.NoError(t, err)

			e := NewEvaluator(nil)
			got, err := e.Eval(ast)
			assert.NoError(t, err)
			assert.Equal(t, tC.exp, got)

			if tC.envAssertion != nil {
				tC.envAssertion(t, e.RootEnv)
			}
		})
	}
}
