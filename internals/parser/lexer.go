package parser

import (
	"fmt"
)

type TokenType string

const (
	TOKEN_EOF          TokenType = "EOF"
	TOKEN_ILLEGAL      TokenType = "ILLEGAL"
	TOKEN_IDENT        TokenType = "IDENT"
	TOKEN_STRING       TokenType = "STRING"
	TOKEN_NUMBER       TokenType = "NUMBER"
	TOKEN_LPAREN       TokenType = "("
	TOKEN_RPAREN       TokenType = ")"
	TOKEN_SQR_PAR_OPEN TokenType = "["
	TOKEN_SQR_PAR_CLS  TokenType = "]"
)

type Token struct {
	Type    TokenType
	Literal string
}

type Lexer struct {
	input        string
	position     int  // current position in input (points to current char)
	readPosition int  // current reading position in input (after current char)
	ch           byte // current char under examination
}

func NewLexer(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition++
}

func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	}
	return l.input[l.readPosition]
}

func (l *Lexer) NextToken() Token {
	var tok Token

	l.skipWhitespace()

	switch l.ch {
	case '(':
		tok = Token{Type: TOKEN_LPAREN, Literal: "("}
	case ')':
		tok = Token{Type: TOKEN_RPAREN, Literal: ")"}
	case '[':
		tok = Token{Type: TOKEN_SQR_PAR_OPEN, Literal: "["}
	case ']':
		tok = Token{Type: TOKEN_SQR_PAR_CLS, Literal: "]"}
	case '"', '\'':
		tok.Type = TOKEN_STRING
		tok.Literal = l.readString(l.ch)
		return tok
	case 0:
		tok.Literal = ""
		tok.Type = TOKEN_EOF
		return tok
	default:
		if isLetter(l.ch) || l.ch == '_' || l.ch == ':' || l.ch == '*' || l.ch == '-' {
			tok.Literal = l.readIdentifier()
			tok.Type = TOKEN_IDENT
			return tok
		} else if isDigit(l.ch) {
			tok.Literal = l.readNumber()
			tok.Type = TOKEN_NUMBER
			return tok
		} else {
			tok = Token{Type: TOKEN_ILLEGAL, Literal: string(l.ch)}
		}
	}

	l.readChar()
	return tok
}

func (l *Lexer) readString(quote byte) string {
	l.readChar() // skip starting quote
	start := l.position
	for l.ch != quote && l.ch != 0 {
		if l.ch == '\\' && (l.peekChar() == quote || l.peekChar() == '\\') {
			l.readChar()
		}
		l.readChar()
	}
	str := l.input[start:l.position]
	if l.ch == quote {
		l.readChar() // skip ending quote
	}
	return str
}

func (l *Lexer) readIdentifier() string {
	start := l.position
	for isLetter(l.ch) || isDigit(l.ch) || l.ch == '_' || l.ch == ':' || l.ch == '-' || l.ch == '.' || l.ch == '*' {
		l.readChar()
	}
	return l.input[start:l.position]
}

func (l *Lexer) readNumber() string {
	start := l.position
	for isDigit(l.ch) || l.ch == '.' {
		l.readChar()
	}
	return l.input[start:l.position]
}

func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z'
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func (l *Lexer) Tokenize() ([]Token, error) {
	var tokens []Token
	for {
		tok := l.NextToken()
		if tok.Type == TOKEN_ILLEGAL {
			return nil, fmt.Errorf("illegal character encountered: %s", tok.Literal)
		}
		tokens = append(tokens, tok)
		if tok.Type == TOKEN_EOF {
			break
		}
	}
	return tokens, nil
}

func Tokenize(input string) ([]Token, error) {
	return NewLexer(input).Tokenize()
}
