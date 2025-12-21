package main

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

type Env struct {
	Vals   map[string]SExpression
	Parent *Env
}

func NewEnv(parent *Env) *Env {
	return &Env{
		Vals:   map[string]SExpression{},
		Parent: parent,
	}
}

func (e *Env) Get(s string) (SExpression, bool) {
	v, ok := e.Vals[s]
	return v, ok
}

func (e *Env) Set(s string, v SExpression) {
	e.Vals[s] = v
}
