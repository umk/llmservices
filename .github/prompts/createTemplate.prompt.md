---
mode: 'agent'
description: 'Scaffold infrastructure for a strongly-typed Go text template'
---
You are an expert Go developer. Your task is to create infrastructure for a strongly typed Go template. For example, if the user wants to create a template for a user message, you would create:

  1. A template file `internal/msg/templates/user_message.tmpl`. You must follow user instructions to implement the template file based on its purpose and parameters. If the user didn't provide instructions for the template, leave the template empty.
  2. A Go file `internal/msg/user_message.go` that contains program logic for rendering the template.

You must follow this example when generating the Go source file:

```go
package msg

import (
	_ "embed"
	"strings"
	"text/template"
)

//go:embed templates/user_message.tmpl
var userMessage string

var userMessageTmpl = template.Must(template.New("user_message").Parse(userMessage))

type UserMessageParams struct {
	Content string
}

func RenderUserMessage(params UserMessageParams) (string, error) {
	var sb strings.Builder
	if err := userMessageTmpl.Execute(&sb, p); err != nil {
		return "", err
	}

	return strings.TrimSpace(sb.String()), nil
}
```

The parameters structure must contain fields that are used in the template. The field names must be in PascalCase format and the comments must describe the purpose of the field. You may implement custom logic for rendering the template, but only if requested by the user. Otherwise you must strictly follow the example above.

If the user doesn't explain the purpose of the template (which is necessary to derive the names of files and types, such as "user message" in the example above), you must ask the user for clarification.
