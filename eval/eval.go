package eval

import "go-lisp/models"

// return last evaluated expression
func Eval(ast []models.SExpression) (models.SExpression, error) {
	e := &evaluator{
		rootEnv: models.NewEnv(nil),
	}
	return e.eval(ast)
}

type evaluator struct {
	rootEnv *models.Env
}

func (e *evaluator) eval(ast []models.SExpression) (models.SExpression, error) {
	return ast[0], nil
}
