package msg

import (
	_ "embed"
	"strings"
	"text/template"
)

//go:embed templates/tool_error_message.tmpl
var toolErrorMessage string

var toolErrorMessageTmpl = template.Must(template.New("tool_error_message").Parse(toolErrorMessage))

type ToolErrorMessageParams struct {
	// Error message that was returned
	Error string
}

func RenderToolErrorMessage(params ToolErrorMessageParams) (string, error) {
	var sb strings.Builder
	if err := toolErrorMessageTmpl.Execute(&sb, params); err != nil {
		return "", err
	}

	return strings.TrimSpace(sb.String()), nil
}
