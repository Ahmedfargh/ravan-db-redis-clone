package parser

import (
	"fmt"
	"strings"
)

type Parser struct {
	query  string
	tokens []Token
	pos    int
}

func NewParser(tokens []Token) *Parser {
	return &Parser{tokens: tokens, pos: 0}
}

func NewParserWithQuery(query string, tokens []Token) *Parser {
	return &Parser{query: query, tokens: tokens, pos: 0}
}

func NewParserFromInput(input string) (*Parser, error) {
	tokens, err := NewLexer(input).Tokenize()
	if err != nil {
		return nil, err
	}
	return NewParserWithQuery(input, tokens), nil
}

func (p *Parser) curToken() Token {
	if p.pos >= len(p.tokens) {
		lastCol := 1
		lastLine := 1
		if len(p.tokens) > 0 {
			lastTok := p.tokens[len(p.tokens)-1]
			lastLine = lastTok.Line
			lastCol = lastTok.Column + len(lastTok.Literal)
		}
		return Token{Type: TOKEN_EOF, Literal: "", Line: lastLine, Column: lastCol}
	}
	return p.tokens[p.pos]
}

func (p *Parser) nextToken() {
	if p.pos < len(p.tokens) {
		p.pos++
	}
}

func (p *Parser) Parse() (*CommandExpr, error) {
	cmd, err := p.ParseCommand()
	if err != nil {
		return nil, err
	}

	// Verify no trailing extra tokens remain after top-level command
	tok := p.curToken()
	if tok.Type != TOKEN_EOF {
		return nil, &RqlSyntaxError{
			Phase:   "Parser",
			Message: fmt.Sprintf("unexpected trailing token '%s' (%s) after complete command expression", tok.Literal, tok.Type),
			Line:    tok.Line,
			Column:  tok.Column,
			Query:   p.query,
			Hint:    "verify argument formatting or enclose nested sub-queries in parentheses (...)",
		}
	}

	return cmd, nil
}

func ParseQuery(input string) (*CommandExpr, error) {
	p, err := NewParserFromInput(input)
	if err != nil {
		return nil, err
	}
	return p.Parse()
}

func (p *Parser) ParseCommand() (*CommandExpr, error) {
	cur := p.curToken()
	if cur.Type == TOKEN_EOF {
		return nil, &RqlSyntaxError{
			Phase:   "Parser",
			Message: "unexpected end of input; expected command name",
			Line:    cur.Line,
			Column:  cur.Column,
			Query:   p.query,
			Hint:    "provide a valid command such as PING, GET, SET, etc.",
		}
	}

	isEnclosedInParen := false
	openParenTok := cur

	if cur.Type == TOKEN_LPAREN {
		isEnclosedInParen = true
		p.nextToken()
		cur = p.curToken()
	}

	if cur.Type != TOKEN_IDENT {
		return nil, &RqlSyntaxError{
			Phase:   "Parser",
			Message: fmt.Sprintf("expected command name, got %s ('%s')", cur.Type, cur.Literal),
			Line:    cur.Line,
			Column:  cur.Column,
			Query:   p.query,
			Hint:    "commands must begin with an identifier (e.g. SET, GET, DEL)",
		}
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
			if isEnclosedInParen {
				return nil, &RqlSyntaxError{
					Phase:   "Parser",
					Message: fmt.Sprintf("unclosed parenthesis '(' opened at col %d (missing matching ')')", openParenTok.Column),
					Line:    tok.Line,
					Column:  tok.Column,
					Query:   p.query,
					Hint:    "ensure every opening '(' has a corresponding closing ')'",
				}
			}
			break
		}
		if isEnclosedInParen && tok.Type == TOKEN_RPAREN {
			p.nextToken()
			break
		}
		if !isEnclosedInParen && tok.Type == TOKEN_RPAREN {
			return nil, &RqlSyntaxError{
				Phase:   "Parser",
				Message: "unexpected closing parenthesis ')' without matching '('",
				Line:    tok.Line,
				Column:  tok.Column,
				Query:   p.query,
				Hint:    "remove unmatched closing parenthesis",
			}
		}
		if tok.Type == TOKEN_SQR_PAR_CLS {
			return nil, &RqlSyntaxError{
				Phase:   "Parser",
				Message: "unexpected closing bracket ']' without matching '['",
				Line:    tok.Line,
				Column:  tok.Column,
				Query:   p.query,
				Hint:    "remove unmatched closing bracket",
			}
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

	case TOKEN_SQR_PAR_OPEN:
		openBracketTok := tok
		p.nextToken()
		list := &ListLiteral{}
		for {
			cur := p.curToken()
			if cur.Type == TOKEN_EOF {
				return nil, &RqlSyntaxError{
					Phase:   "Parser",
					Message: fmt.Sprintf("unclosed list bracket '[' opened at col %d (missing closing ']')", openBracketTok.Column),
					Line:    cur.Line,
					Column:  cur.Column,
					Query:   p.query,
					Hint:    "close list literal with ']'",
				}
			}
			if cur.Type == TOKEN_SQR_PAR_CLS {
				p.nextToken()
				break
			}
			if cur.Type == TOKEN_RPAREN {
				return nil, &RqlSyntaxError{
					Phase:   "Parser",
					Message: "unexpected ')' inside list literal; expected list element or ']'",
					Line:    cur.Line,
					Column:  cur.Column,
					Query:   p.query,
					Hint:    "close nested sub-queries before closing the list, or use ']' to terminate the list",
				}
			}

			node, err := p.parseExpression()
			if err != nil {
				return nil, err
			}
			list.Values = append(list.Values, node)
		}
		return list, nil

	case TOKEN_RPAREN, TOKEN_EOF:
		return nil, nil

	default:
		return nil, &RqlSyntaxError{
			Phase:   "Parser",
			Message: fmt.Sprintf("unexpected token in expression: %s ('%s')", tok.Type, tok.Literal),
			Line:    tok.Line,
			Column:  tok.Column,
			Query:   p.query,
			Hint:    "expected literal (string, number), nested command (...), or list [...]",
		}
	}
}
