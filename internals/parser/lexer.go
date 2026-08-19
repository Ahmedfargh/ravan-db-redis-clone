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
	Type     TokenType
	Literal  string
	Line     int
	Column   int
	Position int
}

type Lexer struct {
	input        string
	position     int  // current position in input (points to current char)
	readPosition int  // current reading position in input (after current char)
	ch           byte // current char under examination
	line         int  // current line (1-based)
	col          int  // current column (1-based)
	lastErr      string
}

func NewLexer(input string) *Lexer {
	l := &Lexer{
		input: input,
		line:  1,
		col:   0,
	}
	l.readChar()
	return l
}

func (l *Lexer) readChar() {
	if l.ch == '\n' {
		l.line++
		l.col = 0
	}

	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition++
	l.col++
}

func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	}
	return l.input[l.readPosition]
}

func (l *Lexer) NextToken() Token {
	l.skipWhitespace()

	startPos := l.position
	startLine := l.line
	startCol := l.col

	var tok Token
	tok.Line = startLine
	tok.Column = startCol
	tok.Position = startPos

	switch l.ch {
	case '(':
		tok.Type = TOKEN_LPAREN
		tok.Literal = "("
	case ')':
		tok.Type = TOKEN_RPAREN
		tok.Literal = ")"
	case '[':
		tok.Type = TOKEN_SQR_PAR_OPEN
		tok.Literal = "["
	case ']':
		tok.Type = TOKEN_SQR_PAR_CLS
		tok.Literal = "]"
	case '"', '\'':
		quoteChar := l.ch
		str, terminated := l.readString(quoteChar)
		if !terminated {
			tok.Type = TOKEN_ILLEGAL
			tok.Literal = str
			l.lastErr = fmt.Sprintf("unterminated string literal starting with quote %c at col %d (missing closing quote)", quoteChar, startCol)
			return tok
		}
		tok.Type = TOKEN_STRING
		tok.Literal = str
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
			tok.Type = TOKEN_ILLEGAL
			tok.Literal = string(l.ch)
			l.lastErr = fmt.Sprintf("illegal character encountered: '%c' (0x%02X)", l.ch, l.ch)
		}
	}

	l.readChar()
	return tok
}

func (l *Lexer) readString(quote byte) (string, bool) {
	l.readChar() // skip starting quote
	start := l.position
	for l.ch != quote && l.ch != 0 {
		if l.ch == '\\' && (l.peekChar() == quote || l.peekChar() == '\\') {
			l.readChar()
		}
		l.readChar()
	}
	if l.ch != quote {
		// Unterminated string
		return l.input[start:l.position], false
	}
	str := l.input[start:l.position]
	l.readChar() // skip ending quote
	return str, true
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
			errMsg := l.lastErr
			if errMsg == "" {
				errMsg = fmt.Sprintf("illegal character encountered: '%s'", tok.Literal)
			}
			return nil, &RqlSyntaxError{
				Phase:   "Lexer",
				Message: errMsg,
				Line:    tok.Line,
				Column:  tok.Column,
				Query:   l.input,
				Hint:    "check for unescaped quotes or invalid special characters",
			}
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
