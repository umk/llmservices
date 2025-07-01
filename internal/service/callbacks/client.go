package callbacks

import (
	"context"
)

func GetFunctionCallRPC(ctx context.Context, req GetFunctionCallRequest, resp *GetFunctionCallResponse) error {
	return (*Client(ctx)).Call(ctx, "getFunctionCall", req, resp)
}

func PushThoughtRPC(ctx context.Context, req PushThoughtRequest) error {
	return (*Client(ctx)).Notify(ctx, "pushThought", req)
}
