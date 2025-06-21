package msg

import (
	_ "embed"
	"strings"
)

//go:embed templates/tool_previous_error_message.tmpl
var toolPreviousErrorMessage string

// RenderToolPreviousErrorMessage returns the static error message as a string.
func RenderToolPreviousErrorMessage() string {
	return strings.TrimSpace(toolPreviousErrorMessage)
}
