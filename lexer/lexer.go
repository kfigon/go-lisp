package lexer

import (
	"fmt"
	"iter"
	"strings"
	"unicode"
)

type TokenType int

const (
	Open TokenType = iota
	Close
	NumberTok
	SymbolTok
	StringTok
	Keyword
)

func (t TokenType) String() string {
	return [...]string{
		"Open",
		"Close",
		"NumberTok",
		"SymbolTok",
		"StringTok",
		"Keyword",
	}[t]
}

type Token struct {
	TokType TokenType
	Lexeme  string
}

func (t Token) String() string {
	return fmt.Sprintf("(%v; %v)", t.TokType, t.Lexeme)
}

func Lex(s string) iter.Seq[Token] {
	return func(yield func(Token) bool) {
		var fn lexerFn = empty
		for _, c := range s {
			next, tok := fn(c)
			for _, t := range tok {
				if !yield(t) {
					return
				}
			}
			fn = next
		}
	}
}

type lexerFn func(rune) (lexerFn, []Token)

func empty(r rune) (lexerFn, []Token) {
	if unicode.IsSpace(r) {
		return empty, nil
	} else if r == '(' {
		return empty, []Token{{TokType: Open}}
	} else if r == ')' {
		return empty, []Token{{TokType: Close}}
	} else if unicode.IsDigit(r) {
		b := &strings.Builder{}
		b.WriteRune(r)
		return makePendingDigit(b), nil
	} else if r == '"' {
		b := &strings.Builder{}
		return makePendingString(b), nil
	}
	b := &strings.Builder{}
	b.WriteRune(r)
	return makePendingSymbol(b), nil
}

func makePendingDigit(pendingStr *strings.Builder) lexerFn {
	return func(r rune) (lexerFn, []Token) {
		if unicode.IsDigit(r) {
			pendingStr.WriteRune(r)
			return makePendingDigit(pendingStr), nil
		}
		out := Token{TokType: NumberTok, Lexeme: pendingStr.String()}
		if r == ')' {
			return empty, []Token{out, {TokType: Close}}
		} else if r == '(' {
			return empty, []Token{out, {TokType: Open}}
		}
		return empty, []Token{out}
	}
}
func makePendingString(pending *strings.Builder) lexerFn {
	return func(r rune) (lexerFn, []Token) {
		if r == '"' {
			return empty, []Token{{TokType: StringTok, Lexeme: pending.String()}}
		}

		pending.WriteRune(r)
		return makePendingString(pending), nil
	}
}
func makePendingSymbol(pending *strings.Builder) lexerFn {
	return func(r rune) (lexerFn, []Token) {
		if unicode.IsSpace(r) {
			return empty, []Token{{TokType: SymbolTok, Lexeme: pending.String()}}
		} else if r == ')' {
			return empty, []Token{{TokType: SymbolTok, Lexeme: pending.String()}, {TokType: Close}}
		} else if r == '(' {
			return empty, []Token{{TokType: SymbolTok, Lexeme: pending.String()}, {TokType: Open}}
		}

		pending.WriteRune(r)
		return makePendingSymbol(pending), nil
	}
}
