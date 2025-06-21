package msg

import (
	_ "embed"
	"strings"
	"text/template"
)

//go:embed templates/summarize_message.tmpl
var summarizeMessage string

var summarizeMessageTmpl = template.Must(template.New("summarize_message").Parse(summarizeMessage))

type SummarizeMessageParams struct{}

func RenderSummarizeMessage(params SummarizeMessageParams) (string, error) {
	var sb strings.Builder
	if err := summarizeMessageTmpl.Execute(&sb, params); err != nil {
		return "", err
	}

	return strings.TrimSpace(sb.String()), nil
}
