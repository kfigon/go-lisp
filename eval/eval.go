package eval

import "go-lisp/models"

// return last evaluated expression
func Eval(ast []models.SExpression) (models.SExpression, error) {

	return ast[0], nil
}
