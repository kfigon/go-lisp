package parser

import (
	"errors"
	"fmt"
	"go-lisp/lexer"
	"go-lisp/models"
	"iter"
	"strconv"
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
	for p.tokenOk {
		got, err := p.parseNode()
		if err != nil {
			p.emitError(err)
		} else {
			p.emitNode(got)
		}

		p.advance()
	}
}

func (p *parser) parseNode() (models.SExpression, error) {
	if p.currentTok.TokType == lexer.NumberTok {
		return p.parseNum()
	} else if p.currentTok.TokType == lexer.StringTok {
		return p.parseString(), nil
	} else if p.currentTok.TokType == lexer.SymbolTok {
		return p.parseSymbol(), nil
	} else if p.currentTok.TokType == lexer.Open {
		return p.parseList()
	}
	panic("unreachable " + p.currentTok.String())
}

func (p *parser) parseNum() (models.Number, error) {
	num, err := strconv.Atoi(p.currentTok.Lexeme)
	if err != nil {
		return 0, fmt.Errorf("error parsing number %q: %w", p.currentTok.Lexeme, err)
	}
	return models.Number(num), nil
}

func (p *parser) parseString() models.String {
	return models.String(p.currentTok.Lexeme)
}

func (p *parser) parseSymbol() models.SExpression {
	if p.currentTok.Lexeme == "nil" {
		return models.Nil{}
	} else if p.currentTok.Lexeme == "true" {
		return models.Bool(true)
	} else if p.currentTok.Lexeme == "false" {
		return models.Bool(false)
	}
	return models.Symbol(p.currentTok.Lexeme)
}

func (p *parser) parseList() (models.List, error) {
	nodes := models.List{}
	p.advance()
	for p.tokenOk && p.currentTok.TokType != lexer.Close {
		got, err := p.parseNode()
		if err != nil {
			return nil, err
		}
		nodes = append(nodes, got)
		p.advance()
	}
	return nodes, nil
}

func (p *parser) emitError(err error) {
	p.errors = append(p.errors, err)
}

func (p *parser) emitNode(m models.SExpression) {
	p.out = append(p.out, m)
}
