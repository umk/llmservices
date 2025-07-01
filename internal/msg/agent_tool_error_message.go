package msg

import (
	_ "embed"
	"strings"
	"text/template"
)

//go:embed templates/agent_tool_error_message.tmpl
var agentToolErrorMessage string

var agentToolErrorMessageTmpl = template.Must(template.New("agent_tool_error_message").Parse(agentToolErrorMessage))

type AgentToolErrorMessageParams struct {
	Name string
}

func RenderAgentToolErrorMessage(params AgentToolErrorMessageParams) (string, error) {
	var sb strings.Builder
	if err := agentToolErrorMessageTmpl.Execute(&sb, params); err != nil {
		return "", err
	}

	return strings.TrimSpace(sb.String()), nil
}
