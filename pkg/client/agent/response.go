package agent

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"unicode"

	"github.com/umk/llmservices/internal/msg"
	"github.com/umk/llmservices/pkg/adapter"
	thread_ "github.com/umk/llmservices/pkg/client/thread"
)

var responseRx = regexp.MustCompile(`<(thought|action|action_input|observation|answer)>`)

type ResponseParams struct {
	adapter.CompletionParams
	Description string          `json:"description"`
	Iterations  int             `json:"iterations" validate:"required,min=1"`
	Retries     *int            `json:"retries,omitempty" validate:"min=0"`
	Handler     ResponseHandler `json:"-"`
}

type Response struct {
	Thread thread_.Thread `json:"thread" validate:"required"`
	Answer string         `json:"answer,omitempty"`
	Error  string         `json:"error,omitempty"`
	Done   bool           `json:"done"`
}

type ResponseHandler interface {
	Call(ctx context.Context, fn adapter.ToolCallFunction) (string, error)
	Thought(ctx context.Context, content string) error
}

func (c *Client) Response(ctx context.Context, thread thread_.Thread, params ResponseParams) (Response, error) {
	if params.ResponseFormat.OfResponseFormatJSONSchema != nil {
		return Response{}, errors.New("cannot use structured output for agent")
	}

	t, err := setSystemMessage(thread, params)
	if err != nil {
		return Response{}, err
	}

	// Cannot use the built-in functionality for tools calling, so just clear
	// the tools in the request parameters.
	params.Tools = nil

	// By default share the budget of retries with iterations.
	iteration := params.Iterations
	retries := &iteration

	if params.Retries != nil {
		retries = new(int)
		*retries = *params.Retries
	}

	// Iterate by getting completions and calling tools
	for {
		iteration--

		r, err := c.responseIterate(ctx, t, params, retries)
		if err != nil {
			return Response{}, err
		}

		if r.Done || iteration == 0 {
			return r, err
		}
	}
}

func (c *Client) responseIterate(
	ctx context.Context,
	thread thread_.Thread,
	params ResponseParams,
	retries *int,
) (Response, error) {
	var output structuredCompl

	for t := thread; ; *retries-- {
		o, err := c.getStructuredCompl(ctx, &t, params.CompletionParams)
		if err != nil {
			return Response{}, err
		}

		if len(o.thoughts) > 0 || o.action != "" || o.done {
			thread = t
			output = o

			break
		}

		if *retries == 0 {
			m, err := msg.RenderAgentFatalErrorMessage(msg.AgentFatalErrorMessageParams{})
			if err != nil {
				return Response{}, err
			}

			return Response{
				Thread: thread, // don't include futile retries into response
				Error:  m,
				Done:   true,
			}, nil
		}

		m, err := msg.RenderAgentErrorMessage(msg.AgentErrorMessageParams{})
		if err != nil {
			return Response{}, err
		}

		f := &t.Frames[len(t.Frames)-1]
		f.Messages = append(f.Messages, adapter.CreateUserMessage(m))
	}

	for _, t := range output.thoughts {
		if err := params.Handler.Thought(ctx, t); err != nil {
			return Response{}, err
		}
	}

	if output.done {
		return Response{
			Thread: thread,
			Answer: output.answer,
			Done:   true,
		}, nil
	}

	if output.action != "" {
		resp, err := params.Handler.Call(ctx, adapter.ToolCallFunction{
			Name:      output.action,
			Arguments: output.parameter,
		})
		if err != nil {
			m, renderErr := msg.RenderToolErrorMessage(msg.ToolErrorMessageParams{
				Error: err.Error(),
			})
			if renderErr != nil {
				m = err.Error()
			}
			resp = m
		}

		v := fmt.Sprintf("<observation>%s</observation>", resp)

		f := &thread.Frames[len(thread.Frames)-1]
		f.Messages = append(f.Messages, adapter.CreateUserMessage(v))
	}

	return Response{
		Thread: thread,
		Done:   false,
	}, nil
}

func (c *Client) getStructuredCompl(ctx context.Context, thread *thread_.Thread, params adapter.CompletionParams) (structuredCompl, error) {
	resp, err := (*thread_.Client)(c).Completion(ctx, *thread, params)
	if err != nil {
		return structuredCompl{}, err
	}

	r, err := resp.Thread.Response()
	if err != nil {
		return structuredCompl{}, err
	}

	if r.Refusal != nil || r.Content == nil {
		return structuredCompl{
			answer: *r.Refusal,
			done:   true,
		}, nil
	}

	*thread = resp.Thread

	return parseResponse(*r.Content), nil
}

func setSystemMessage(thread thread_.Thread, params ResponseParams) (thread_.Thread, error) {
	// Create system message that contains agent description and the tools
	// available for calling.
	c, err := msg.RenderAgentSystemMessage(msg.AgentSystemMessageParams{
		Description: params.Description,
		Tools:       params.Tools,
	})
	if err != nil {
		return thread_.Thread{}, err
	}

	s := adapter.CreateSystemMessage(c)

	// Assign the system message to the thread, or replace the system message.
	for i := 0; ; {
		if i < len(thread.Frames) {
			f := &thread.Frames[i]
			if len(f.Messages) == 0 {
				continue
			}

			thread.Frames = thread.Frames[i:]
			if m := &f.Messages[0]; m.OfSystemMessage != nil {
				*m = s
				f.Tokens = 0

				return thread, nil
			}
		}

		f := thread_.MessagesFrame{Messages: []adapter.Message{s}}
		thread.Frames = append([]thread_.MessagesFrame{f}, thread.Frames...)

		return thread, nil
	}
}

type structuredCompl struct {
	thoughts    []string
	action      string
	parameter   string
	observation string
	answer      string
	done        bool
}

func parseResponse(response string) (output structuredCompl) {
	indexes := responseRx.FindAllStringSubmatchIndex(response, -1)
	indexes = append(indexes, []int{len(response)})
	for i := range len(indexes) - 1 {
		tag := response[indexes[i][2]:indexes[i][3]]
		start, end := indexes[i][1], indexes[i+1][0]
		content := strings.TrimSpace(response[start:end])
		closing := fmt.Sprintf("</%s>", tag)
		if strings.HasSuffix(content, closing) {
			content = strings.TrimSuffix(strings.TrimRightFunc(content, unicode.IsSpace), closing)
		}
		content = strings.TrimSpace(content)
		switch tag {
		case "thought":
			if content != "" {
				output.thoughts = append(output.thoughts, content)
			}
		case "action":
			output.action = content
		case "action_input":
			output.parameter = content
		case "observation":
			output.observation = content
		case "answer":
			output.answer = content
			output.done = true
		}
	}

	return
}
