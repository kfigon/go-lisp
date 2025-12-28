package eval

import (
	"fmt"
	"go-lisp/models"
)

type EnvFun func(*Evaluator, ...models.SExpression) (models.SExpression, error)
type Env struct {
	Vals   map[string]EnvFun
	Parent *Env
}

func NewEnv(parent *Env) *Env {
	return &Env{
		Vals:   map[string]EnvFun{},
		Parent: parent,
	}
}

func (e *Env) Get(s string) (EnvFun, bool) {
	v, ok := e.Vals[s]
	if ok {
		return v, ok
	} else if e.Parent != nil {
		return e.Parent.Get(s)
	}
	return nil, false
}

func (e *Env) Set(s string, v models.SExpression) {
	e.Vals[s] = func(_ *Evaluator, _ ...models.SExpression) (models.SExpression, error) {
		return v, nil
	}
}

type Evaluator struct {
	RootEnv *Env
}

func NewEvaluator(rootEnv *Env) *Evaluator {
	e := &Evaluator{}
	if rootEnv == nil {
		e.RootEnv = NewEnv(nil)
		e.initStdLib()
	} else {
		e.RootEnv = rootEnv
	}

	return e
}

func (e *Evaluator) initStdLib() {
	// e.RootEnv.Vals["print"] = models.Nil{}
	e.RootEnv.Vals["+"] = createArithmeticOp(func(a, b int) int { return a + b })
	e.RootEnv.Vals["-"] = createArithmeticOp(func(a, b int) int { return a - b })
	e.RootEnv.Vals["*"] = createArithmeticOp(func(a, b int) int { return a * b })
	e.RootEnv.Vals["/"] = createArithmeticOp(func(a, b int) int { return a / b })
	e.RootEnv.Vals["="] = func(ev *Evaluator, args ...models.SExpression) (models.SExpression, error) {
		if len(args) != 2 {
			return nil, fmt.Errorf("'=' expects 2 args, got %d", len(args))
		}
		a, err := ev.evalSingle(args[0])
		if err != nil {
			return nil, err
		}

		b, err := ev.evalSingle(args[1])
		if err != nil {
			return nil, err
		}

		aN, aOk := a.(models.Number)
		bN, bOk := b.(models.Number)
		if aOk && bOk {
			return models.Bool(int(aN) == int(bN)), nil
		}

		aB, aOk := a.(models.Bool)
		bB, bOk := b.(models.Bool)
		if aOk && bOk {
			return models.Bool(bool(aB) == bool(bB)), nil
		}
		aS, aOk := a.(models.String)
		bS, bOk := b.(models.String)
		if aOk && bOk {
			return models.Bool(string(aS) == string(bS)), nil
		}

		aS2, aOk := a.(models.Symbol)
		bS2, bOk := b.(models.Symbol)
		if aOk && bOk {
			return models.Bool(string(aS2) == string(bS2)), nil
		}

		return nil, fmt.Errorf("type mismatch: %T, %T", a, b)
	}
	e.RootEnv.Vals["!="] = func(ev *Evaluator, args ...models.SExpression) (models.SExpression, error) {
		eqFn, _ := e.RootEnv.Get("=")
		got, err := eqFn(ev, args...)
		if err != nil {
			return nil, err
		}
		b, ok := got.(models.Bool)
		if !ok {
			return nil, fmt.Errorf("non boolean arg: %T", got)
		}
		return models.Bool(!bool(b)), nil
	}
	e.RootEnv.Vals["set"] = func(ev *Evaluator, args ...models.SExpression) (models.SExpression, error) {
		if len(args) != 2 {
			return nil, fmt.Errorf("set expects 2 args, got %d", len(args))
		}
		key, ok := args[0].(models.Symbol)
		if !ok {
			return nil, fmt.Errorf("first arg to set should be a symbol, got %T", args[0])
		}
		res, err := ev.evalSingle(args[1])
		if err != nil {
			return nil, err
		}
		ev.RootEnv.Set(string(key), res)
		return models.Nil{}, nil
	}
	e.RootEnv.Vals["and"] = func(ev *Evaluator, args ...models.SExpression) (models.SExpression, error) {
		if len(args) != 2 {
			return nil, fmt.Errorf("and expects 2 args, got %d", len(args))
		}
		a, err := ev.evalSingle(args[0])
		if err != nil {
			return nil, err
		}
		aB, ok := a.(models.Bool)
		if !ok {
			return nil, fmt.Errorf("and operand A should be bool, got %T", a)
		}
		if !aB {
			return models.Bool(false), nil
		}

		b, err := ev.evalSingle(args[1])
		if err != nil {
			return nil, err
		}
		bB, ok := b.(models.Bool)
		if !ok {
			return nil, fmt.Errorf("and operand B should be bool, got %T", b)
		}
		return models.Bool(aB && bB), nil
	}

	e.RootEnv.Vals["or"] = func(ev *Evaluator, args ...models.SExpression) (models.SExpression, error) {
		if len(args) != 2 {
			return nil, fmt.Errorf("or expects 2 args, got %d", len(args))
		}
		a, err := ev.evalSingle(args[0])
		if err != nil {
			return nil, err
		}
		aB, ok := a.(models.Bool)
		if !ok {
			return nil, fmt.Errorf("or operand A should be bool, got %T", a)
		}
		if aB {
			return models.Bool(true), nil
		}

		b, err := ev.evalSingle(args[1])
		if err != nil {
			return nil, err
		}
		bB, ok := b.(models.Bool)
		if !ok {
			return nil, fmt.Errorf("or operand B should be bool, got %T", b)
		}
		return models.Bool(aB || bB), nil
	}

	e.RootEnv.Vals["if"] = func(ev *Evaluator, args ...models.SExpression) (models.SExpression, error) {
		if len(args) != 3 {
			return nil, fmt.Errorf("if expects 3 args, got %d", len(args))
		}
		condition, err := ev.evalSingle(args[0])
		if err != nil {
			return nil, err
		}
		pred, ok := condition.(models.Bool)
		if !ok {
			return nil, fmt.Errorf("if expects boolean as predicate result, got %T", condition)
		}

		if pred {
			return ev.evalSingle(args[1])
		}
		return ev.evalSingle(args[2])
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
	case models.Bool, models.Number, models.String:
		return ex, nil
	case models.Symbol:
		fn, ok := e.RootEnv.Get(string(v))
		if !ok {
			return nil, fmt.Errorf("unknown symbol %s", v)
		}
		return fn(e)
	case models.List:
		return e.evalList(v)
	case *models.Function:
		return e.evalFunctionDeclaration(v)
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
	return val(e, v[1:]...)
}

func createArithmeticOp(op func(int, int) int) EnvFun {
	return func(ev *Evaluator, args ...models.SExpression) (models.SExpression, error) {
		var res *int
		for _, v := range args {
			got, err := ev.evalSingle(v)
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

func (e *Evaluator) evalFunctionDeclaration(v *models.Function) (models.SExpression, error) {
	e.RootEnv.Vals[v.Name] = func(ev *Evaluator, s ...models.SExpression) (models.SExpression, error) {
		if len(s) != len(v.Args) {
			return nil, fmt.Errorf("invalid numer of args to function, exp %d, got %d", len(v.Args), len(s))
		}
		evaluatedArgs := []models.SExpression{}
		for i, a := range s {
			got, err := ev.evalSingle(a)
			if err != nil {
				return nil, fmt.Errorf("error evaluating %dth argument to %s: %w", i, v.Name, err)
			}
			evaluatedArgs = append(evaluatedArgs, got)
		}

		newEnv := NewEnv(e.RootEnv)
		for i, arg := range evaluatedArgs {
			newEnv.Set(string(v.Args[i]), arg)
		}
		newE := NewEvaluator(newEnv)
		l, _ := v.Body.(models.List)
		return newE.evalSingle(l[0])
	}
	return models.Nil{}, nil
}

func ptr[T any](v T) *T {
	return &v
}
