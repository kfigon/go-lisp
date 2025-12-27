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
	switch p.currentTok.TokType {
	case lexer.NumberTok:
		return p.parseNum()
	case lexer.StringTok:
		return p.parseString(), nil
	case lexer.SymbolTok:
		return p.parseSymbol(), nil
	case lexer.Open:
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
	switch p.currentTok.Lexeme {
	case "nil":
		return models.Nil{}
	case "true":
		return models.Bool(true)
	case "false":
		return models.Bool(false)
	}
	return models.Symbol(p.currentTok.Lexeme)
}

func (p *parser) parseList() (models.SExpression, error) {
	p.advance() // (
	if p.currentTok.Lexeme == "lambda" {
		return p.parseLambda()
	}

	nodes := models.List{}
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

func (p *parser) parseLambda() (*models.Function, error) {
	p.advance() // lambda
	if p.currentTok.TokType != lexer.SymbolTok {
		return nil, fmt.Errorf("invalid lambda name, expected symbol, got %s", p.currentTok.TokType)
	}
	name := p.currentTok.Lexeme
	p.advance() // symbol
	if p.currentTok.TokType != lexer.Open {
		return nil, fmt.Errorf("invalid lambda arg list, expected list, got %s", p.currentTok.TokType)
	}
	p.advance() // (
	argList := []models.Symbol{}
	for p.tokenOk && p.currentTok.TokType != lexer.Close {
		if p.currentTok.TokType != lexer.SymbolTok {
			return nil, fmt.Errorf("invalid lambda declaration, expected list of symbols as args, got %s", p.currentTok.TokType)
		}
		argList = append(argList, models.Symbol(p.currentTok.Lexeme))
		p.advance()
	}
	p.advance() // )
	body, err := p.parseNode()
	if err != nil {
		return nil, fmt.Errorf("erro parsing body of lambda: %w", err)
	}
	p.advance() // )
	return &models.Function{
		Name: name,
		Args: argList,
		Body: body,
	}, nil
}

func (p *parser) emitError(err error) {
	p.errors = append(p.errors, err)
}

func (p *parser) emitNode(m models.SExpression) {
	p.out = append(p.out, m)
}
