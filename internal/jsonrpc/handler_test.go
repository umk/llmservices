package jsonrpc

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
)

// MockRPCContext implements the RPCContext interface for testing
type MockRPCContext struct {
	req         *rpcRequest
	requestBody any
	response    any
	err         error
}

func (m *MockRPCContext) GetRequestBody(v any) error {
	if m.err != nil {
		return m.err
	}

	// If requestBody is set, use it to populate v
	if m.requestBody != nil {
		data, err := json.Marshal(m.requestBody)
		if err != nil {
			return err
		}
		return json.Unmarshal(data, v)
	}

	// Otherwise, use the request params if they exist
	if m.req != nil && m.req.Params != nil {
		return json.Unmarshal(*m.req.Params, v)
	}

	return nil
}

func (m *MockRPCContext) GetResponse(v any) (any, error) {
	if m.err != nil {
		return nil, m.err
	}
	m.response = v
	return v, nil
}

func TestNewHandler(t *testing.T) {
	funcs := map[string]HandlerFunc{
		"testMethod": func(ctx context.Context, c RPCContext) (any, error) {
			return "result", nil
		},
	}

	handler := NewHandler(funcs)

	if handler == nil {
		t.Fatal("Expected non-nil handler")
	}

	if len(handler.funcs) != len(funcs) {
		t.Errorf("Expected %d functions, got %d", len(funcs), len(handler.funcs))
	}

	if _, ok := handler.funcs["testMethod"]; !ok {
		t.Errorf("Expected 'testMethod' to be in handler functions")
	}
}

func TestHandler_Handle_ParseError(t *testing.T) {
	handler := NewHandler(nil)

	// Invalid JSON
	data := []byte(`{"jsonrpc": "2.0", "method": "test", "params": {,}}`)

	resp, err := handler.Handle(context.Background(), data)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	var result rpcResponse
	if err := json.Unmarshal(resp, &result); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if result.Error == nil {
		t.Fatal("Expected error response")
	}

	if result.Error.Code != -32700 {
		t.Errorf("Expected parse error code -32700, got %d", result.Error.Code)
	}

	if result.Error.Message != "Parse error" {
		t.Errorf("Expected 'Parse error' message, got %s", result.Error.Message)
	}
}

func TestHandler_Handle_InvalidRequest(t *testing.T) {
	handler := NewHandler(nil)

	testCases := []struct {
		name        string
		requestJSON string
		expectCode  int
		expectMsg   string
	}{
		{
			name:        "Missing JSONRPC Version",
			requestJSON: `{"method": "test", "params": {}, "id": 1}`,
			expectCode:  -32600,
			expectMsg:   "Invalid request",
		},
		{
			name:        "Wrong JSONRPC Version",
			requestJSON: `{"jsonrpc": "1.0", "method": "test", "params": {}, "id": 1}`,
			expectCode:  -32600,
			expectMsg:   "Invalid request",
		},
		{
			name:        "Missing Method",
			requestJSON: `{"jsonrpc": "2.0", "params": {}, "id": 1}`,
			expectCode:  -32600,
			expectMsg:   "Invalid request",
		},
		{
			name:        "Empty Method",
			requestJSON: `{"jsonrpc": "2.0", "method": "", "params": {}, "id": 1}`,
			expectCode:  -32600,
			expectMsg:   "Invalid request",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resp, err := handler.Handle(context.Background(), []byte(tc.requestJSON))
			if err != nil {
				t.Fatalf("Expected no error, got %v", err)
			}

			var result rpcResponse
			if err := json.Unmarshal(resp, &result); err != nil {
				t.Fatalf("Failed to unmarshal response: %v", err)
			}

			if result.Error == nil {
				t.Fatal("Expected error response")
			}

			if result.Error.Code != tc.expectCode {
				t.Errorf("Expected error code %d, got %d", tc.expectCode, result.Error.Code)
			}

			if result.Error.Message != tc.expectMsg {
				t.Errorf("Expected message '%s', got '%s'", tc.expectMsg, result.Error.Message)
			}
		})
	}
}

func TestHandler_Handle_MethodNotFound(t *testing.T) {
	handler := NewHandler(map[string]HandlerFunc{
		"existingMethod": func(ctx context.Context, c RPCContext) (any, error) {
			return "result", nil
		},
	})

	data := []byte(`{"jsonrpc": "2.0", "method": "nonExistingMethod", "params": {}, "id": 1}`)

	resp, err := handler.Handle(context.Background(), data)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	var result rpcResponse
	if err := json.Unmarshal(resp, &result); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if result.Error == nil {
		t.Fatal("Expected error response")
	}

	if result.Error.Code != -32601 {
		t.Errorf("Expected method not found code -32601, got %d", result.Error.Code)
	}

	if result.Error.Message != "Method not found" {
		t.Errorf("Expected 'Method not found' message, got %s", result.Error.Message)
	}
}

func TestHandler_Handle_Success(t *testing.T) {
	handler := NewHandler(map[string]HandlerFunc{
		"testMethod": func(ctx context.Context, c RPCContext) (any, error) {
			var params struct {
				Name string `json:"name"`
				Age  int    `json:"age"`
			}
			if err := c.GetRequestBody(&params); err != nil {
				return nil, err
			}

			return map[string]any{
				"greeting": "Hello, " + params.Name,
				"age":      params.Age,
			}, nil
		},
	})

	data := []byte(`{
		"jsonrpc": "2.0", 
		"method": "testMethod", 
		"params": {"name": "John", "age": 30}, 
		"id": 123
	}`)

	resp, err := handler.Handle(context.Background(), data)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	var result rpcResponse
	if err := json.Unmarshal(resp, &result); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if result.Error != nil {
		t.Fatalf("Expected no error in response, got %v", result.Error)
	}

	resultMap, ok := result.Result.(map[string]any)
	if !ok {
		t.Fatalf("Expected map result, got %T", result.Result)
	}

	greeting, ok := resultMap["greeting"].(string)
	if !ok || greeting != "Hello, John" {
		t.Errorf("Expected greeting 'Hello, John', got %v", resultMap["greeting"])
	}

	age, ok := resultMap["age"].(float64) // JSON numbers are floats
	if !ok || int(age) != 30 {
		t.Errorf("Expected age 30, got %v", resultMap["age"])
	}

	// Check ID is passed through
	var id float64
	idBytes, _ := json.Marshal(result.Id)
	json.Unmarshal(idBytes, &id)
	if id != 123 {
		t.Errorf("Expected id 123, got %v", id)
	}
}

func TestHandler_Handle_HandlerError(t *testing.T) {
	customErr := Error{
		Code:    -32000,
		Message: "Custom application error",
		Data:    "Additional error data",
	}

	handler := NewHandler(map[string]HandlerFunc{
		"errorMethod": func(ctx context.Context, c RPCContext) (any, error) {
			return nil, customErr
		},
		"regularError": func(ctx context.Context, c RPCContext) (any, error) {
			return nil, errors.New("regular error")
		},
	})

	t.Run("RPC Error", func(t *testing.T) {
		data := []byte(`{"jsonrpc": "2.0", "method": "errorMethod", "id": 1}`)

		resp, err := handler.Handle(context.Background(), data)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		var result rpcResponse
		if err := json.Unmarshal(resp, &result); err != nil {
			t.Fatalf("Failed to unmarshal response: %v", err)
		}

		if result.Error == nil {
			t.Fatal("Expected error response")
		}

		if result.Error.Code != customErr.Code {
			t.Errorf("Expected code %d, got %d", customErr.Code, result.Error.Code)
		}

		if result.Error.Message != customErr.Message {
			t.Errorf("Expected message '%s', got '%s'", customErr.Message, result.Error.Message)
		}

		if result.Error.Data != customErr.Data {
			t.Errorf("Expected data %v, got %v", customErr.Data, result.Error.Data)
		}
	})

	t.Run("Regular Error", func(t *testing.T) {
		data := []byte(`{"jsonrpc": "2.0", "method": "regularError", "id": 1}`)

		resp, err := handler.Handle(context.Background(), data)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		var result rpcResponse
		if err := json.Unmarshal(resp, &result); err != nil {
			t.Fatalf("Failed to unmarshal response: %v", err)
		}

		if result.Error == nil {
			t.Fatal("Expected error response")
		}

		if result.Error.Code != -32603 {
			t.Errorf("Expected internal error code -32603, got %d", result.Error.Code)
		}

		if result.Error.Message != "Internal error" {
			t.Errorf("Expected message 'Internal error', got '%s'", result.Error.Message)
		}
	})
}

func TestHandler_Handle_Notification(t *testing.T) {
	methodCalled := false

	handler := NewHandler(map[string]HandlerFunc{
		"notificationMethod": func(ctx context.Context, c RPCContext) (any, error) {
			methodCalled = true
			return "result", nil
		},
	})

	// Notification request (no ID)
	data := []byte(`{"jsonrpc": "2.0", "method": "notificationMethod", "params": {}}`)

	resp, err := handler.Handle(context.Background(), data)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if resp != nil {
		t.Errorf("Expected nil response for notification, got %s", string(resp))
	}

	if !methodCalled {
		t.Error("Expected notification method to be called")
	}
}
