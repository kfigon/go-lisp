package parser

import (
	"errors"
	"go-lisp/lexer"
	"go-lisp/models"
	"iter"
)

func Parse(toks iter.Seq[lexer.Token]) ([]models.SExpression, error) {
	nextFn, stop := iter.Pull(toks)
	defer stop()
	p := &parser{nextFn: nextFn}
	p.advance()
	p.parse()

	return p.out, errors.Join(p.errors...)
}

type parser struct {
	nextFn     func() (lexer.Token, bool)
	currentTok lexer.Token
	tokenOk    bool

	out    []models.SExpression
	errors []error
}

func (p *parser) advance() {
	p.currentTok, p.tokenOk = p.nextFn()
}

func (p *parser) parse() {
}
