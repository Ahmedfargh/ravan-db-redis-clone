package parser

import (
	"fmt"
	"strings"
)

type Parser struct {
	tokens []Token
	pos    int
}

func NewParser(tokens []Token) *Parser {
	return &Parser{tokens: tokens, pos: 0}
}

func (p *Parser) curToken() Token {
	if p.pos >= len(p.tokens) {
		return Token{Type: TOKEN_EOF, Literal: ""}
	}
	return p.tokens[p.pos]
}

func (p *Parser) nextToken() {
	if p.pos < len(p.tokens) {
		p.pos++
	}
}

func ParseQuery(input string) (*CommandExpr, error) {
	tokens, err := Tokenize(input)
	if err != nil {
		return nil, err
	}
	p := NewParser(tokens)
	return p.ParseCommand()
}

func (p *Parser) ParseCommand() (*CommandExpr, error) {
	cur := p.curToken()
	isEnclosedInParen := false

	if cur.Type == TOKEN_LPAREN {
		isEnclosedInParen = true
		p.nextToken()
		cur = p.curToken()
	}

	if cur.Type != TOKEN_IDENT {
		return nil, fmt.Errorf("expected command name, got %s ('%s')", cur.Type, cur.Literal)
	}

	cmdName := strings.ToUpper(cur.Literal)
	p.nextToken()

	cmdExpr := &CommandExpr{
		CommandName: cmdName,
		Args:        make([]Node, 0),
	}

	for {
		tok := p.curToken()
		if tok.Type == TOKEN_EOF {
			break
		}
		if isEnclosedInParen && tok.Type == TOKEN_RPAREN {
			p.nextToken()
			break
		}

		node, err := p.parseExpression()
		if err != nil {
			return nil, err
		}
		if node == nil {
			break
		}
		cmdExpr.Args = append(cmdExpr.Args, node)
	}

	return cmdExpr, nil
}

func (p *Parser) parseExpression() (Node, error) {
	tok := p.curToken()

	switch tok.Type {
	case TOKEN_STRING:
		p.nextToken()
		return &StringLiteral{Value: tok.Literal}, nil

	case TOKEN_NUMBER:
		p.nextToken()
		return &NumberLiteral{Value: tok.Literal}, nil

	case TOKEN_IDENT:
		p.nextToken()
		return &StringLiteral{Value: tok.Literal}, nil

	case TOKEN_LPAREN:
		// Nested command expression like (GET key)
		return p.ParseCommand()

	case TOKEN_RPAREN, TOKEN_EOF:
		return nil, nil

	default:
		return nil, fmt.Errorf("unexpected token in expression: %s ('%s')", tok.Type, tok.Literal)
	}
}
