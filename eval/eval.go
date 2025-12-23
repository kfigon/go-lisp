package eval

import (
	"fmt"
	"go-lisp/models"
)

type Evaluator struct {
	RootEnv *models.Env
}

func NewEvaluator() *Evaluator {
	rootEnv := models.NewEnv(nil)
	initStdLib(rootEnv)
	return &Evaluator{
		RootEnv: rootEnv,
	}
}

func initStdLib(e *models.Env) {
	// todo: init special forms, basic operators etc
}

func (e *Evaluator) Eval(ast []models.SExpression) (models.SExpression, error) {
	var lastEx models.SExpression
	for _, v := range ast {
		var err error
		lastEx, err = e.evalSingle(v)
		if err != nil {
			return nil, fmt.Errorf("error evaluating %v: %w", v, err)
		}
	}
	return lastEx, nil
}

func (e *Evaluator) evalSingle(ex models.SExpression) (models.SExpression, error) {
	switch ex.(type) {
	case models.Bool, models.Number, models.String, models.Symbol:
		return ex, nil
	case models.List:
		panic("todo")
	default:
		return nil, fmt.Errorf("invalid expression: %T", ex)
	}
}
