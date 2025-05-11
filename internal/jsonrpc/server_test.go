package jsonrpc

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"strings"
	"testing"
	"time"
)

// mockReader simulates various reader behaviors for testing
type mockReader struct {
	data    []byte
	pos     int
	err     error
	delay   time.Duration
	readErr bool
}

func (r *mockReader) Read(p []byte) (int, error) {
	if r.readErr {
		return 0, r.err
	}

	if r.delay > 0 {
		time.Sleep(r.delay)
	}

	if r.pos >= len(r.data) {
		return 0, io.EOF
	}

	n := copy(p, r.data[r.pos:])
	r.pos += n
	if r.pos >= len(r.data) && r.err != nil {
		return n, r.err
	}
	return n, nil
}

func TestNewServer(t *testing.T) {
	t.Run("WithDefaultOptions", func(t *testing.T) {
		handler := NewHandler(nil)
		server := NewServer(handler)

		if server == nil {
			t.Fatal("Expected non-nil server")
		}

		if server.handler != handler {
			t.Error("Handler not correctly set")
		}

		// Default request size is 4K
		if server.bufferPool == nil {
			t.Error("Buffer pool not initialized")
		}
	})

	t.Run("WithCustomOptions", func(t *testing.T) {
		handler := NewHandler(nil)
		customSize := 8192
		server := NewServer(handler, WithRequestSize(customSize))

		if server == nil {
			t.Fatal("Expected non-nil server")
		}

		if server.handler != handler {
			t.Error("Handler not correctly set")
		}
	})
}

func TestServer_Run(t *testing.T) {
	t.Run("ProcessSingleRequest", func(t *testing.T) {
		requestJSON := `{"jsonrpc":"2.0","method":"test","id":1}` + "\n"

		// Create a handler that returns a predefined response
		handler := NewHandler(map[string]HandlerFunc{
			"test": func(ctx context.Context, c RPCContext) (any, error) {
				return "success", nil
			},
		})

		server := NewServer(handler)

		input := strings.NewReader(requestJSON)
		var output bytes.Buffer

		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
		defer cancel()

		err := server.Run(ctx, input, &output)
		// After timeout, we should get a context canceled error or nil
		if err != nil && err != context.DeadlineExceeded {
			t.Errorf("Server.Run() unexpected error = %v", err)
		}

		// Check that we got a valid response
		outputStr := output.String()
		var resp rpcResponse
		if err := json.Unmarshal([]byte(strings.TrimSuffix(outputStr, "\n")), &resp); err != nil {
			t.Errorf("Failed to parse response: %v", err)
		}

		if resp.Error != nil {
			t.Errorf("Got error response: %v", resp.Error)
		}

		if resp.Result != "success" {
			t.Errorf("Expected result 'success', got %v", resp.Result)
		}
	})

	t.Run("ProcessMultipleRequests", func(t *testing.T) {
		request1 := `{"jsonrpc":"2.0","method":"method1","id":1}` + "\n"
		request2 := `{"jsonrpc":"2.0","method":"method2","id":2}` + "\n"

		// Create a handler with multiple methods
		handler := NewHandler(map[string]HandlerFunc{
			"method1": func(ctx context.Context, c RPCContext) (any, error) {
				return "result1", nil
			},
			"method2": func(ctx context.Context, c RPCContext) (any, error) {
				return "result2", nil
			},
		})

		server := NewServer(handler)

		input := strings.NewReader(request1 + request2)
		var output bytes.Buffer

		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
		defer cancel()

		err := server.Run(ctx, input, &output)
		// After timeout, we should get a context canceled error or nil
		if err != nil && err != context.DeadlineExceeded {
			t.Errorf("Server.Run() unexpected error = %v", err)
		}

		// Check we got both responses (order may vary due to concurrency)
		outputStr := output.String()
		if !strings.Contains(outputStr, `"result":"result1"`) || !strings.Contains(outputStr, `"result":"result2"`) {
			t.Errorf("Expected both results. Got: %s", outputStr)
		}

		if !strings.Contains(outputStr, `"id":1`) || !strings.Contains(outputStr, `"id":2`) {
			t.Errorf("Expected both IDs. Got: %s", outputStr)
		}
	})

	t.Run("HandleReaderError", func(t *testing.T) {
		mockReader := &mockReader{
			err:     io.ErrUnexpectedEOF,
			readErr: true,
		}

		handler := NewHandler(nil)
		server := NewServer(handler)
		var output bytes.Buffer

		err := server.Run(context.Background(), mockReader, &output)
		if err != io.ErrUnexpectedEOF {
			t.Errorf("Expected io.ErrUnexpectedEOF, got %v", err)
		}
	})

	t.Run("HandleEOF", func(t *testing.T) {
		mockReader := &mockReader{
			data: []byte{}, // Empty data will trigger EOF
		}

		handler := NewHandler(nil)
		server := NewServer(handler)
		var output bytes.Buffer

		err := server.Run(context.Background(), mockReader, &output)
		if err != nil {
			t.Errorf("Expected no error on EOF, got %v", err)
		}
	})

	t.Run("HandleNullResponse", func(t *testing.T) {
		requestJSON := `{"jsonrpc":"2.0","method":"notification"}` + "\n" // No ID = notification

		// Create a handler that processes notifications
		handler := NewHandler(map[string]HandlerFunc{
			"notification": func(ctx context.Context, c RPCContext) (any, error) {
				// For a notification, return doesn't matter
				return nil, nil
			},
		})

		server := NewServer(handler)

		input := strings.NewReader(requestJSON)
		var output bytes.Buffer

		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
		defer cancel()

		err := server.Run(ctx, input, &output)
		// After timeout, we should get a context canceled error or nil
		if err != nil && err != context.DeadlineExceeded {
			t.Errorf("Server.Run() unexpected error = %v", err)
		}

		// Should have no output for a notification request
		if output.Len() > 0 {
			t.Errorf("Expected empty output for notification, got: %s", output.String())
		}
	})

	t.Run("HandleHandlerError", func(t *testing.T) {
		requestJSON := `{"jsonrpc":"2.0","method":"errorMethod","id":1}` + "\n"

		// Create a handler that returns an error
		handler := NewHandler(map[string]HandlerFunc{
			"errorMethod": func(ctx context.Context, c RPCContext) (any, error) {
				return nil, errors.New("handler error")
			},
		})

		server := NewServer(handler)

		input := strings.NewReader(requestJSON)
		var output bytes.Buffer

		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
		defer cancel()

		err := server.Run(ctx, input, &output)
		// After timeout, we should get a context canceled error or nil
		if err != nil && err != context.DeadlineExceeded {
			t.Errorf("Server.Run() unexpected error = %v", err)
		}

		// Should get an error response
		outputStr := output.String()
		if !strings.Contains(outputStr, `"error"`) {
			t.Errorf("Expected error in response, got: %s", outputStr)
		}
	})
}

func TestReadInput(t *testing.T) {
	t.Run("SingleLine", func(t *testing.T) {
		input := "test input\n"
		reader := bufio.NewReader(strings.NewReader(input))

		data := make([]byte, 0, 64)
		err := readInput(reader, &data)

		if err != nil {
			t.Errorf("readInput() error = %v", err)
		}

		if string(data) != "test input" {
			t.Errorf("readInput() = %q, want %q", string(data), "test input")
		}
	})

	t.Run("MultiLine", func(t *testing.T) {
		input := "line 1\nline 2\n"
		reader := bufio.NewReader(strings.NewReader(input))

		// First call should read "line 1"
		data := make([]byte, 0, 64)
		err := readInput(reader, &data)

		if err != nil {
			t.Errorf("readInput() error = %v", err)
		}

		if string(data) != "line 1" {
			t.Errorf("readInput() = %q, want %q", string(data), "line 1")
		}

		// Second call should read "line 2"
		data = data[:0]
		err = readInput(reader, &data)

		if err != nil {
			t.Errorf("readInput() error = %v", err)
		}

		if string(data) != "line 2" {
			t.Errorf("readInput() = %q, want %q", string(data), "line 2")
		}
	})

	t.Run("LongLine", func(t *testing.T) {
		// Generate a long line that will exceed the default buffer
		longLine := strings.Repeat("a", 8192) + "\n"
		reader := bufio.NewReader(strings.NewReader(longLine))

		data := make([]byte, 0, 64)
		err := readInput(reader, &data)

		if err != nil {
			t.Errorf("readInput() error = %v", err)
		}

		if len(data) != 8192 {
			t.Errorf("readInput() len = %d, want %d", len(data), 8192)
		}
	})

	t.Run("Error", func(t *testing.T) {
		mockReader := &mockReader{
			err:     io.ErrUnexpectedEOF,
			readErr: true,
		}
		reader := bufio.NewReader(mockReader)

		data := make([]byte, 0, 64)
		err := readInput(reader, &data)

		if err != io.ErrUnexpectedEOF {
			t.Errorf("Expected io.ErrUnexpectedEOF, got %v", err)
		}
	})

	t.Run("NoNewline", func(t *testing.T) {
		// Input without a trailing newline
		input := "input without newline"
		reader := bufio.NewReader(strings.NewReader(input))

		data := make([]byte, 0, 64)
		err := readInput(reader, &data)

		if err != nil {
			t.Errorf("readInput() error = %v", err)
		}

		if string(data) != "input without newline" {
			t.Errorf("readInput() = %q, want %q", string(data), "input without newline")
		}
	})
}
