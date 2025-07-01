package msg

import (
	_ "embed"
	"strings"
	"text/template"
)

//go:embed templates/agent_fatal_error_message.tmpl
var agentFatalErrorMessage string

var agentFatalErrorMessageTmpl = template.Must(template.New("agent_fatal_error_message").Parse(agentFatalErrorMessage))

type AgentFatalErrorMessageParams struct{}

func RenderAgentFatalErrorMessage(params AgentFatalErrorMessageParams) (string, error) {
	var sb strings.Builder
	if err := agentFatalErrorMessageTmpl.Execute(&sb, params); err != nil {
		return "", err
	}

	return strings.TrimSpace(sb.String()), nil
}
