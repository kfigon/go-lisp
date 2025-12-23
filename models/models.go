package models

type SExpression interface {
	Exp()
}

type Number int

func (Number) Exp() {}

type Symbol string

func (Symbol) Exp() {}

type String string

func (String) Exp() {}

type Bool bool

func (Bool) Exp() {}

type List []SExpression

func (List) Exp() {}

type Nil struct{}

func (Nil) Exp() {}

type EnvFun func(...SExpression) (SExpression, error)
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
	return v, ok
}

func (e *Env) Set(s string, v SExpression) {
	e.Vals[s] = func(s ...SExpression) (SExpression, error) {
		return v, nil
	}
}
