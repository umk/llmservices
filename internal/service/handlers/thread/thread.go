package thread

import (
	"context"

	"github.com/umk/jsonrpc2"
	"github.com/umk/llmservices/internal/service/callbacks"
	"github.com/umk/llmservices/pkg/client/thread"
)

func GetResponseRPC(ctx context.Context, c jsonrpc2.RPCContext) (any, error) {
	var req GetResponseRequest
	if err := c.Request(&req); err != nil {
		return nil, err
	}

	cl, err := GetClient(ctx, req.ClientID)
	if err != nil {
		return nil, err
	}

	req.Params.Handler = callbacks.Callback{}

	resp, err := cl.Response(ctx, req.Thread, req.Params)
	if err != nil {
		return nil, newResponseError(err)
	}

	return c.Response(GetResponseResponse{
		Response: resp,
	})
}

func GetCompletionRPC(ctx context.Context, c jsonrpc2.RPCContext) (any, error) {
	var req GetCompletionRequest
	if err := c.Request(&req); err != nil {
		return nil, err
	}

	cl, err := GetClient(ctx, req.ClientID)
	if err != nil {
		return nil, err
	}

	resp, err := cl.Completion(ctx, req.Thread, req.Params)
	if err != nil {
		return nil, newCompletionError(err)
	}

	return c.Response(GetCompletionResponse{
		Completion: resp,
	})
}

func GetSummaryRPC(ctx context.Context, c jsonrpc2.RPCContext) (any, error) {
	var req GetSummaryRequest
	if err := c.Request(&req); err != nil {
		return nil, err
	}

	// If the generator is not specified, use the same client as the summarizer.
	if req.GenClientID == nil {
		req.GenClientID = &req.ClientID
	}

	// Get both the summarizer and the generator clients.
	cl, err := GetClient(ctx, req.ClientID)
	if err != nil {
		return nil, err
	}

	gen, err := GetClient(ctx, *req.GenClientID)
	if err != nil {
		return nil, err
	}

	// Calculate the number tokens for each frame in the thread
	thread.SetFrameTokens(&req.Thread, gen.Samples)

	// Summarize the thread.
	var opts []thread.SummarizerOption
	if req.MaxMessages != nil {
		opts = append(opts, thread.WithMaxMessages(*req.MaxMessages))
	}
	if req.MaxTokens != nil {
		opts = append(opts, thread.WithMaxTokens(*req.MaxTokens))
	}

	if len(opts) != 0 {
		return nil, errSummarizerParams
	}

	s := thread.NewSummarizer(cl, req.Fraction, opts...)

	t, err := s.Summarize(ctx, req.Thread)
	if err != nil {
		return nil, newSummarizerError(err)
	}

	// Once again calculate the number tokens for each frame in the thread
	thread.SetFrameTokens(&t, gen.Samples)

	return c.Response(GetSummaryResponse{
		Thread: t,
	})
}
