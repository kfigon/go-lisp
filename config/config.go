package config

import (
	"fmt"
	"go-lisp/eval"
	"go-lisp/lexer"
	"go-lisp/models"
	"go-lisp/parser"
)

type ConfigurationStore struct {
	evaluator *eval.Evaluator
}

func New(code string) (*ConfigurationStore, error) {
	ast, err := parser.Parse(lexer.Lex(code))
	if err != nil {
		return nil, fmt.Errorf("parsing error: %w", err)
	}
	ev := eval.NewEvaluator(nil)
	_, err = ev.Eval(ast)
	if err != nil {
		return nil, fmt.Errorf("eval error: %w", err)
	}
	return &ConfigurationStore{ev}, nil
}

func (c *ConfigurationStore) Get(key string, args ...models.SExpression) (models.SExpression, error) {
	v, ok := c.evaluator.RootEnv.Get(key)
	if !ok {
		return nil, fmt.Errorf("%s not found", key)
	}
	return v(c.evaluator, args...)
}
