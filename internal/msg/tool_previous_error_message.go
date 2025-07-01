package msg

import (
	_ "embed"
	"strings"
	"text/template"
)

//go:embed templates/tool_previous_error_message.tmpl
var toolPreviousErrorMessage string

var toolPreviousErrorMessageTmpl = template.Must(template.New("tool_previous_error_message").Parse(toolPreviousErrorMessage))

type ToolPreviousErrorMessageParams struct{}

func RenderToolPreviousErrorMessage(params ToolPreviousErrorMessageParams) (string, error) {
	var sb strings.Builder
	if err := toolPreviousErrorMessageTmpl.Execute(&sb, params); err != nil {
		return "", err
	}

	return strings.TrimSpace(sb.String()), nil
}
