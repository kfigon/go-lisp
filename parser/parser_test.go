package parser

import (
	"go-lisp/lexer"
	"go-lisp/models"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test(t *testing.T) {
	testCases := []struct {
		desc string
		code string

		exp []models.SExpression
		err error
	}{
		{
			desc: "basic num",
			code: `123`,
			exp:  []models.SExpression{models.Number(123)},
		},
		{
			desc: "sum",
			code: `(+ 123 1 2)`,
			exp: []models.SExpression{
				models.List{
					models.Symbol("+"),
					models.Number(123),
					models.Number(1),
					models.Number(2),
				},
			},
		},
		{
			desc: "couple of lists",
			code: `(+ 123 1 2)
					(* 2 3)`,
			exp: []models.SExpression{
				models.List{
					models.Symbol("+"),
					models.Number(123),
					models.Number(1),
					models.Number(2),
				},
				models.List{
					models.Symbol("*"),
					models.Number(2),
					models.Number(3),
				},
			},
		},
		{
			desc: "nested",
			code: `(+ 123 (+ 1 2) (- 3 4))`,
			exp: []models.SExpression{
				models.List{
					models.Symbol("+"),
					models.Number(123),
					models.List{
						models.Symbol("+"),
						models.Number(1),
						models.Number(2),
					},
					models.List{
						models.Symbol("-"),
						models.Number(3),
						models.Number(4),
					},
				},
			},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			got, err := Parse(lexer.Lex(tC.code))

			if tC.err != nil {
				assert.ErrorContains(t, err, tC.err.Error())
			} else {
				assert.Equal(t, tC.exp, got)
			}
		})
	}
}
