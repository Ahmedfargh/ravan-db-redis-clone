package parser

import (
	"strings"
	"testing"
)

func TestLexerQuotedStringsAndParens(t *testing.T) {
	input := `SET key "Hello World" (GET other_key)`
	lexer := NewLexer(input)
	tokens, err := lexer.Tokenize()
	if err != nil {
		t.Fatalf("Unexpected lexer error: %v", err)
	}

	expectedLiterals := []string{"SET", "key", "Hello World", "(", "GET", "other_key", ")", ""}
	if len(tokens) != len(expectedLiterals) {
		t.Fatalf("Expected %d tokens, got %d", len(expectedLiterals), len(tokens))
	}

	for i, tok := range tokens {
		if tok.Literal != expectedLiterals[i] {
			t.Errorf("Token %d: expected %q, got %q", i, expectedLiterals[i], tok.Literal)
		}
	}
}

func TestLexerSquareBrackets(t *testing.T) {
	input := `RPUSH list [item1 "item 2" 100]`
	lexer := NewLexer(input)
	tokens, err := lexer.Tokenize()
	if err != nil {
		t.Fatalf("Unexpected lexer error: %v", err)
	}

	expectedTypes := []TokenType{
		TOKEN_IDENT, TOKEN_IDENT, TOKEN_SQR_PAR_OPEN,
		TOKEN_IDENT, TOKEN_STRING, TOKEN_NUMBER, TOKEN_SQR_PAR_CLS, TOKEN_EOF,
	}

	if len(tokens) != len(expectedTypes) {
		t.Fatalf("Expected %d tokens, got %d", len(expectedTypes), len(tokens))
	}

	for i, tok := range tokens {
		if tok.Type != expectedTypes[i] {
			t.Errorf("Token %d: expected type %s, got %s", i, expectedTypes[i], tok.Type)
		}
	}
}

func TestParserNestedExpressions(t *testing.T) {
	input := `SET greeting (GET default_greeting)`
	p, err := NewParserFromInput(input)
	if err != nil {
		t.Fatalf("Unexpected parser creation error: %v", err)
	}

	cmdExpr, err := p.Parse()
	if err != nil {
		t.Fatalf("Unexpected parser error: %v", err)
	}

	if cmdExpr.CommandName != "SET" {
		t.Fatalf("Expected command name SET, got %s", cmdExpr.CommandName)
	}

	if len(cmdExpr.Args) != 2 {
		t.Fatalf("Expected 2 args, got %d", len(cmdExpr.Args))
	}

	if cmdExpr.Args[0].String() != "greeting" {
		t.Fatalf("Expected first arg 'greeting', got %s", cmdExpr.Args[0].String())
	}

	subCmd, ok := cmdExpr.Args[1].(*CommandExpr)
	if !ok {
		t.Fatalf("Expected second arg to be *CommandExpr, got %T", cmdExpr.Args[1])
	}

	if subCmd.CommandName != "GET" || subCmd.Args[0].String() != "default_greeting" {
		t.Fatalf("Expected sub-command (GET default_greeting), got %s", subCmd.String())
	}
}

func TestParserListLiteral(t *testing.T) {
	input := `SET list [val1 "val 2" (GET k1)]`
	p, err := NewParserFromInput(input)
	if err != nil {
		t.Fatalf("Unexpected parser creation error: %v", err)
	}

	cmdExpr, err := p.Parse()
	if err != nil {
		t.Fatalf("Unexpected parser error: %v", err)
	}

	if len(cmdExpr.Args) != 2 {
		t.Fatalf("Expected 2 args, got %d", len(cmdExpr.Args))
	}

	listNode, ok := cmdExpr.Args[1].(*ListLiteral)
	if !ok {
		t.Fatalf("Expected second arg to be *ListLiteral, got %T", cmdExpr.Args[1])
	}

	if len(listNode.Values) != 3 {
		t.Fatalf("Expected 3 list elements, got %d", len(listNode.Values))
	}

	if listNode.Values[0].String() != "val1" {
		t.Errorf("Expected first element 'val1', got %s", listNode.Values[0].String())
	}

	if listNode.Values[1].String() != "val 2" {
		t.Errorf("Expected second element 'val 2', got %s", listNode.Values[1].String())
	}

	subCmd, ok := listNode.Values[2].(*CommandExpr)
	if !ok {
		t.Fatalf("Expected third element to be nested *CommandExpr, got %T", listNode.Values[2])
	}
	if subCmd.CommandName != "GET" || subCmd.Args[0].String() != "k1" {
		t.Errorf("Expected sub-command (GET k1), got %s", subCmd.String())
	}
}

func TestParserEmptyListLiteral(t *testing.T) {
	input := `SET empty_list []`
	p, err := NewParserFromInput(input)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	cmdExpr, err := p.Parse()
	if err != nil {
		t.Fatalf("Unexpected parser error: %v", err)
	}

	listNode, ok := cmdExpr.Args[1].(*ListLiteral)
	if !ok {
		t.Fatalf("Expected *ListLiteral, got %T", cmdExpr.Args[1])
	}

	if len(listNode.Values) != 0 {
		t.Fatalf("Expected 0 elements, got %d", len(listNode.Values))
	}
}

func TestLexerUnterminatedStringError(t *testing.T) {
	input := `SET key "unterminated string`
	_, err := NewLexer(input).Tokenize()
	if err == nil {
		t.Fatalf("Expected error for unterminated string, got nil")
	}

	syntaxErr, ok := err.(*RqlSyntaxError)
	if !ok {
		t.Fatalf("Expected *RqlSyntaxError, got %T: %v", err, err)
	}

	if syntaxErr.Phase != "Lexer" {
		t.Errorf("Expected Lexer phase, got %s", syntaxErr.Phase)
	}

	errStr := err.Error()
	if !strings.Contains(errStr, "[RQL Lexer Error]") {
		t.Errorf("Expected [RQL Lexer Error] banner, got:\n%s", errStr)
	}
	if !strings.Contains(errStr, "^") {
		t.Errorf("Expected caret pointer in error output, got:\n%s", errStr)
	}
}

func TestParserUnclosedParenthesisError(t *testing.T) {
	input := `SET role (GET default_role`
	_, err := ParseQuery(input)
	if err == nil {
		t.Fatalf("Expected error for unclosed parenthesis, got nil")
	}

	syntaxErr, ok := err.(*RqlSyntaxError)
	if !ok {
		t.Fatalf("Expected *RqlSyntaxError, got %T: %v", err, err)
	}

	if !strings.Contains(syntaxErr.Message, "unclosed parenthesis") {
		t.Errorf("Expected 'unclosed parenthesis' message, got: %s", syntaxErr.Message)
	}

	errStr := err.Error()
	if !strings.Contains(errStr, "[RQL Parser Error]") {
		t.Errorf("Expected [RQL Parser Error] banner, got:\n%s", errStr)
	}
}

func TestParserUnclosedBracketError(t *testing.T) {
	input := `SET tags [ "server" 8080`
	_, err := ParseQuery(input)
	if err == nil {
		t.Fatalf("Expected error for unclosed list bracket, got nil")
	}

	syntaxErr, ok := err.(*RqlSyntaxError)
	if !ok {
		t.Fatalf("Expected *RqlSyntaxError, got %T: %v", err, err)
	}

	if !strings.Contains(syntaxErr.Message, "unclosed list bracket") {
		t.Errorf("Expected 'unclosed list bracket' message, got: %s", syntaxErr.Message)
	}
}

func TestParserTrailingTokensError(t *testing.T) {
	input := `(GET key) extra_unexpected_token`
	_, err := ParseQuery(input)
	if err == nil {
		t.Fatalf("Expected error for trailing token, got nil")
	}

	syntaxErr, ok := err.(*RqlSyntaxError)
	if !ok {
		t.Fatalf("Expected *RqlSyntaxError, got %T: %v", err, err)
	}

	if !strings.Contains(syntaxErr.Message, "unexpected trailing token") {
		t.Errorf("Expected trailing token error, got: %s", syntaxErr.Message)
	}
}

