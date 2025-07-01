package msg

import (
	_ "embed"
	"strings"
	"text/template"

	"github.com/umk/llmservices/pkg/adapter"
)

//go:embed templates/agent_system_message.tmpl
var agentSystemMessage string

var agentSystemMessageTmpl = template.Must(template.New("agent_system_message").Parse(agentSystemMessage))

type AgentSystemMessageParams struct {
	Description string
	Tools       []adapter.Tool
}

func RenderAgentSystemMessage(params AgentSystemMessageParams) (string, error) {
	var sb strings.Builder
	if err := agentSystemMessageTmpl.Execute(&sb, params); err != nil {
		return "", err
	}

	return strings.TrimSpace(sb.String()), nil
}
