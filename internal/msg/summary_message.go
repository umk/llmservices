package msg

import (
	_ "embed"
	"strings"
	"text/template"
)

//go:embed templates/summary_message.tmpl
var summaryMessage string

var summaryMessageTmpl = template.Must(template.New("summary_message").Parse(summaryMessage))

type SummaryMessageParams struct {
	// The conversation summary to be presented
	Summary string
}

func RenderSummaryMessage(p SummaryMessageParams) (string, error) {
	var sb strings.Builder
	if err := summaryMessageTmpl.Execute(&sb, p); err != nil {
		return "", err
	}

	return strings.TrimSpace(sb.String()), nil
}
