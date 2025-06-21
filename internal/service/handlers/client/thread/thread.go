package thread

import (
	"context"

	"github.com/umk/jsonrpc2"
	serviceclient "github.com/umk/llmservices/internal/service/client"
	"github.com/umk/llmservices/pkg/client/thread"
)

func GetResponseRPC(ctx context.Context, c jsonrpc2.RPCContext) (any, error) {
	var req getThreadResponseRequest
	if err := c.GetRequestBody(&req); err != nil {
		return nil, err
	}

	cl, err := GetClient(req.ClientID)
	if err != nil {
		return nil, err
	}

	resp, err := cl.Response(ctx, req.Thread, thread.ResponseParams{
		CompletionParams: req.Params,
		Caller: &serviceclient.FunctionCaller{
			Client: serviceclient.Client,
		},
	})
	if err != nil {
		return nil, newResponseError(err)
	}

	return c.GetResponse(getThreadResponseResponse{
		Thread: resp,
	})
}

func GetCompletionRPC(ctx context.Context, c jsonrpc2.RPCContext) (any, error) {
	var req getThreadCompletionRequest
	if err := c.GetRequestBody(&req); err != nil {
		return nil, err
	}

	cl, err := GetClient(req.ClientID)
	if err != nil {
		return nil, err
	}

	resp, err := cl.Completion(ctx, req.Thread, req.Params)
	if err != nil {
		return nil, newCompletionError(err)
	}

	return c.GetResponse(getThreadCompletionResponse{
		Completion: resp,
	})
}

func GetSummaryRPC(ctx context.Context, c jsonrpc2.RPCContext) (any, error) {
	var req getThreadSummaryRequest
	if err := c.GetRequestBody(&req); err != nil {
		return nil, err
	}

	// If the generator is not specified, use the same client as the summarizer.
	if req.GenClientID == nil {
		req.GenClientID = &req.ClientID
	}

	// Get both the summarizer and the generator clients.
	cl, err := GetClient(req.ClientID)
	if err != nil {
		return nil, err
	}

	gen, err := GetClient(*req.GenClientID)
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

	return c.GetResponse(getThreadSummaryResponse{
		Thread: t,
	})
}
