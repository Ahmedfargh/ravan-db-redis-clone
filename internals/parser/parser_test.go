package parser

import (
	"testing"
)

func TestLexerQuotedStringsAndParens(t *testing.T) {
	input := `SET key "Hello World" (GET other_key)`
	tokens, err := Tokenize(input)
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

func TestParserNestedExpressions(t *testing.T) {
	input := `SET greeting (GET default_greeting)`
	cmdExpr, err := ParseQuery(input)
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
