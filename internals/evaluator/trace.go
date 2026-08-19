package evaluator

import (
	"fmt"
	"strings"
)

// TraceFrame represents one stack frame in the AST evaluation path.
type TraceFrame struct {
	Step        int
	Description string
	Expression  string
}

// EvalError captures the complete hierarchical evaluation trace and root cause error.
type EvalError struct {
	Query  string
	Frames []TraceFrame
	Cause  error
	Hint   string
}

func (e *EvalError) Error() string {
	var sb strings.Builder
	sb.WriteString("[RQL Runtime Error] Execution failed during query evaluation:\n")

	totalFrames := len(e.Frames)
	for i, frame := range e.Frames {
		prefix := "  ├─"
		if i == totalFrames-1 {
			prefix = "  └─"
		}
		if frame.Expression != "" {
			sb.WriteString(fmt.Sprintf("%s Step %d: %s -> %s\n", prefix, i+1, frame.Description, frame.Expression))
		} else {
			sb.WriteString(fmt.Sprintf("%s Step %d: %s\n", prefix, i+1, frame.Description))
		}
	}

	causeMsg := "unknown runtime failure"
	if e.Cause != nil {
		causeMsg = e.Cause.Error()
	}
	sb.WriteString(fmt.Sprintf("  Cause: %s", causeMsg))

	if e.Hint != "" {
		sb.WriteString(fmt.Sprintf("\n  Hint: %s", e.Hint))
	}

	return sb.String()
}

// Unwrap returns the underlying cause error.
func (e *EvalError) Unwrap() error {
	return e.Cause
}

// PushFrame prepends or appends a trace frame to an existing EvalError, or creates a new one.
func wrapWithFrame(err error, description string, expression string) *EvalError {
	if err == nil {
		return nil
	}

	frame := TraceFrame{
		Description: description,
		Expression:  expression,
	}

	if evalErr, ok := err.(*EvalError); ok {
		// Prepend frame so that outermost frames are listed first
		evalErr.Frames = append([]TraceFrame{frame}, evalErr.Frames...)
		return evalErr
	}

	hint := generateHintFromCause(err)
	return &EvalError{
		Frames: []TraceFrame{frame},
		Cause:  err,
		Hint:   hint,
	}
}

func generateHintFromCause(err error) string {
	if err == nil {
		return ""
	}
	msg := strings.ToLower(err.Error())
	if strings.Contains(msg, "unknown command") {
		return "verify command name spelling or check registered custom commands"
	}
	if strings.Contains(msg, "wrong number of arguments") {
		return "check the required argument count for this command"
	}
	if strings.Contains(msg, "out of bounds") || strings.Contains(msg, "out of range") {
		return "ensure index is within valid array or string bounds"
	}
	if strings.Contains(msg, "not found") || strings.Contains(msg, "don't exists") {
		return "verify the key exists in data store before operating on it"
	}
	return ""
}
