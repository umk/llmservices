package thread

import (
	"context"
	"fmt"
	"slices"

	"github.com/umk/llmservices/internal/msg"
	"github.com/umk/llmservices/pkg/adapter"
)

type ResponseParams struct {
	adapter.CompletionParams
	Iterations int             `json:"iterations" validate:"required,min=1"`
	Handler    ResponseHandler `json:"-"`
}

type Response struct {
	Thread Thread `json:"thread" validate:"required"`
	Done   bool   `json:"done"`
}

type ResponseHandler interface {
	Call(ctx context.Context, fn adapter.ToolCallFunction) (string, error)
}

func (c *Client) Response(ctx context.Context, thread Thread, params ResponseParams) (Response, error) {
	for range params.Iterations {
		resp, err := c.Completion(ctx, thread, params.CompletionParams)
		if err != nil {
			return Response{}, err
		}

		r, err := resp.Thread.Response()
		if err != nil {
			return Response{}, err
		}

		if len(r.ToolCalls) == 0 {
			return Response{
				Thread: resp.Thread,
				Done:   true,
			}, nil
		}

		for _, c := range r.ToolCalls {
			if !slices.ContainsFunc(params.Tools, func(t adapter.Tool) bool {
				return t.Function.Name == c.Function.Name
			}) {
				return Response{}, fmt.Errorf("calling not existing function: %s", c.Function.Name)
			}
		}

		if params.Handler == nil {
			return Response{}, fmt.Errorf("function caller is not specified")
		}

		f := &resp.Thread.Frames[len(resp.Thread.Frames)-1]

		for i, c := range r.ToolCalls {
			resp, err := params.Handler.Call(ctx, c.Function)
			if err != nil {
				m, renderErr := msg.RenderToolErrorMessage(msg.ToolErrorMessageParams{
					Error: err.Error(),
				})
				if renderErr != nil {
					m = err.Error()
				}
				f.Messages = append(f.Messages, adapter.CreateToolMessage(c.ID, m))
				for i++; i < len(r.ToolCalls); i++ {
					c := r.ToolCalls[i]
					m, err := msg.RenderToolPreviousErrorMessage(msg.ToolPreviousErrorMessageParams{})
					if err != nil {
						m = "Ignored because one of the previous calls ended with an error."
					}
					f.Messages = append(f.Messages, adapter.CreateToolMessage(c.ID, m))
				}
				return Response{}, err
			}
			f.Messages = append(f.Messages, adapter.CreateToolMessage(c.ID, resp))
		}

		thread = resp.Thread
	}

	return Response{
		Thread: thread,
		Done:   false,
	}, nil
}
