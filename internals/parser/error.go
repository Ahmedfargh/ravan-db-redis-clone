package parser

import (
	"fmt"
	"strings"
)

// RqlSyntaxError represents a syntax or lexical error during RQL parsing with exact location context.
type RqlSyntaxError struct {
	Phase   string // e.g. "Lexer" or "Parser"
	Message string // Explanation of what went wrong
	Line    int    // 1-based line number
	Column  int    // 1-based column number
	Query   string // Original raw query string
	Hint    string // Optional helpful suggestion
}

func (e *RqlSyntaxError) Error() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("[RQL %s Error] at line %d, col %d: %s\n", e.Phase, e.Line, e.Column, e.Message))
	sb.WriteString(FormatErrorSnippet(e.Query, e.Line, e.Column))
	if e.Hint != "" {
		sb.WriteString(fmt.Sprintf("\n  Hint: %s", e.Hint))
	}
	return sb.String()
}

// FormatErrorSnippet generates a visual pointer (caret ^) at the specified column for context.
func FormatErrorSnippet(query string, line int, col int) string {
	if query == "" {
		return ""
	}

	lines := strings.Split(query, "\n")
	targetLineIdx := line - 1
	if targetLineIdx < 0 || targetLineIdx >= len(lines) {
		targetLineIdx = 0
	}

	queryLine := lines[targetLineIdx]
	if col < 1 {
		col = 1
	}

	// Indent 4 spaces for readability
	codeLine := fmt.Sprintf("    %s", queryLine)
	
	// Create caret line
	caretOffset := col - 1
	if caretOffset > len(queryLine) {
		caretOffset = len(queryLine)
	}
	caretLine := fmt.Sprintf("    %s^", strings.Repeat(" ", caretOffset))

	return fmt.Sprintf("%s\n%s", codeLine, caretLine)
}
