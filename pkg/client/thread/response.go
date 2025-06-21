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
	Caller FunctionCaller
}

type FunctionCaller interface {
	Call(ctx context.Context, fn adapter.ToolCallFunction) (string, error)
}

func (c *Client) Response(ctx context.Context, thread Thread, params ResponseParams) (Thread, error) {
	for {
		resp, err := c.Completion(ctx, thread, params.CompletionParams)
		if err != nil {
			return Thread{}, err
		}

		r, err := resp.Thread.Response()
		if err != nil {
			return Thread{}, err
		}

		if len(r.ToolCalls) == 0 {
			return resp.Thread, nil
		}

		for _, c := range r.ToolCalls {
			if !slices.ContainsFunc(params.Tools, func(t adapter.Tool) bool {
				return t.Function.Name == c.Function.Name
			}) {
				return Thread{}, fmt.Errorf("calling not existing function: %s", c.Function.Name)
			}
		}

		if params.Caller == nil {
			return Thread{}, fmt.Errorf("function caller is not specified")
		}

		f := &resp.Thread.Frames[len(resp.Thread.Frames)-1]

		for i, c := range r.ToolCalls {
			resp, err := params.Caller.Call(ctx, c.Function)
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
					f.Messages = append(f.Messages, adapter.CreateToolMessage(c.ID, msg.RenderToolPreviousErrorMessage()))
				}
				return Thread{}, err
			}
			f.Messages = append(f.Messages, adapter.CreateToolMessage(c.ID, resp))
		}

		thread = resp.Thread
	}
}
