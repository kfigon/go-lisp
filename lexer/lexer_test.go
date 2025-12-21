package lexer

import (
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLexer(t *testing.T) {
	testCases := []struct {
		code     string
		expected []Token
	}{
		{
			code: `foo (12 23 x "asdf x" (x) (1 2))`,
			expected: []Token{
				{SymbolTok, "foo"},
				{Open, ""},
				{NumberTok, "12"},
				{NumberTok, "23"},
				{SymbolTok, "x"},
				{StringTok, "asdf x"},
				{Open, ""},
				{SymbolTok, "x"},
				{Close, ""},
				{Open, ""},
				{NumberTok, "1"},
				{NumberTok, "2"},
				{Close, ""},
				{Close, ""},
			},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.code, func(t *testing.T) {
			got := Lex(tC.code)
			assert.Equal(t, tC.expected, slices.Collect(got))
		})
	}
}
