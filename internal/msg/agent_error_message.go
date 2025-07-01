package msg

import (
	_ "embed"
	"strings"
	"text/template"
)

//go:embed templates/agent_error_message.tmpl
var agentErrorMessage string

var agentErrorMessageTmpl = template.Must(template.New("agent_error_message").Parse(agentErrorMessage))

type AgentErrorMessageParams struct{}

func RenderAgentErrorMessage(params AgentErrorMessageParams) (string, error) {
	var sb strings.Builder
	if err := agentErrorMessageTmpl.Execute(&sb, params); err != nil {
		return "", err
	}

	return strings.TrimSpace(sb.String()), nil
}
