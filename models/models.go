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

type Function struct {
	Name string
	Args []Symbol
	Body SExpression
}

func (*Function) Exp() {}
