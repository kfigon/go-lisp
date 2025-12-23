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
	e.Vals["set"] = models.Nil{}
	e.Vals["lambda"] = models.Nil{}
	e.Vals["print"] = models.Nil{}
	e.Vals["+"] = models.Nil{}
	e.Vals["-"] = models.Nil{}
	e.Vals["*"] = models.Nil{}
	e.Vals["/"] = models.Nil{}
	e.Vals["="] = models.Nil{}
	e.Vals["!="] = models.Nil{}
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
	switch v := ex.(type) {
	case models.Bool, models.Number, models.String, models.Symbol:
		return ex, nil
	case models.List:
		return e.evalList(v)
	default:
		return nil, fmt.Errorf("invalid expression: %T", ex)
	}
}

func (e *Evaluator) evalList(v models.List) (models.SExpression, error) {
	if len(v) < 2 {
		return nil, fmt.Errorf("invalid list of len %d, require at least 2", len(v))
	}

	op := v[0]
	symbol, ok := op.(models.Symbol)
	if !ok {
		return nil, fmt.Errorf("expected first symbol of list to be symbol, got %T", op)
	}

	variable, ok := e.RootEnv.Get(string(symbol))
	if !ok {
		return nil, fmt.Errorf("could not find symbol %s", symbol)
	}
	_ = variable

	panic("todo")
}
