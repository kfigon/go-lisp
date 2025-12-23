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
	e := &Evaluator{
		RootEnv: rootEnv,
	}
	e.initStdLib()
	return e
}

func (e *Evaluator) initStdLib() {
	// e.RootEnv.Vals["if"] = models.Nil{}
	// e.RootEnv.Vals["set"] = models.Nil{}
	// e.RootEnv.Vals["lambda"] = models.Nil{}
	// e.RootEnv.Vals["print"] = models.Nil{}
	e.RootEnv.Vals["+"] = createArithmeticOp(e, func(a, b int) int { return a + b })
	e.RootEnv.Vals["-"] = createArithmeticOp(e, func(a, b int) int { return a - b })
	e.RootEnv.Vals["*"] = createArithmeticOp(e, func(a, b int) int { return a * b })
	e.RootEnv.Vals["/"] = createArithmeticOp(e, func(a, b int) int { return a / b })
	// e.RootEnv.Vals["="] = models.Nil{}
	// e.RootEnv.Vals["!="] = models.Nil{}
	e.RootEnv.Vals["set"] = func(in ...models.SExpression) (models.SExpression, error) {
		args := in[0].(models.List)
		if len(args) != 2 {
			return nil, fmt.Errorf("set expects 2 args, got %d", len(args))
		}
		key, ok := args[0].(models.Symbol)
		if !ok {
			return nil, fmt.Errorf("first arg to set should be a symbol, got %T", args[0])
		}
		res, err := e.evalSingle(args[1])
		if err != nil {
			return nil, err
		}
		e.RootEnv.Set(string(key), res)
		return models.Nil{}, nil
	}
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

	val, ok := e.RootEnv.Get(string(symbol))
	if !ok {
		return nil, fmt.Errorf("%s not found", symbol)
	}
	return val(v[1:])
}

func createArithmeticOp(e *Evaluator, op func(int, int) int) models.EnvFun {
	return func(args ...models.SExpression) (models.SExpression, error) {
		list := args[0].(models.List)
		var res *int
		for _, v := range list {
			got, err := e.evalSingle(v)
			if err != nil {
				return nil, err
			}
			num, ok := got.(models.Number)
			if !ok {
				return nil, fmt.Errorf("not a number: %v", got)
			}

			if res != nil {
				res = ptr(op(*res, int(num)))
			} else {
				res = ptr(int(num))
			}
		}
		return models.Number(*res), nil
	}
}

func ptr[T any](v T) *T {
	return &v
}
