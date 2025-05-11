package jsonrpc

import (
	"encoding/json"
	"errors"
	"strings"
	"testing"
)

type person struct {
	Name string `json:"name" validate:"required"`
	Age  int    `json:"age" validate:"gte=0"`
}

func TestRPCParseError(t *testing.T) {
	baseErr := errors.New("test error")
	parseErr := rpcParseError{err: baseErr}

	// Test Error method
	if parseErr.Error() != "failed to parse RPC request: test error" {
		t.Errorf("Expected specific error message, got: %s", parseErr.Error())
	}

	// Test Unwrap method
	if !errors.Is(parseErr, baseErr) {
		t.Errorf("Expected parseErr to wrap baseErr, but it didn't")
	}

	// Test RPCError method
	rpcErr := parseErr.RPCError()
	if rpcErr.Code != -32602 {
		t.Errorf("Expected code -32602, got: %d", rpcErr.Code)
	}
	if rpcErr.Message != baseErr.Error() {
		t.Errorf("Expected message to be base error message, got: %s", rpcErr.Message)
	}
}

func TestGetRequestBody_WithParams(t *testing.T) {
	// Create sample JSON params
	params := json.RawMessage(`{"name":"John","age":30}`)
	req := &rpcRequest{
		Params: &params,
	}

	ctx := &rpcContext{req: req}
	var data person

	err := ctx.GetRequestBody(&data)
	if err != nil {
		t.Fatalf("GetRequestBody failed: %v", err)
	}

	if data.Name != "John" {
		t.Errorf("Expected name 'John', got '%s'", data.Name)
	}

	if data.Age != 30 {
		t.Errorf("Expected age 30, got %d", data.Age)
	}
}

func TestGetRequestBody_InvalidJSON(t *testing.T) {
	// Create invalid JSON params
	params := json.RawMessage(`{"name":"John","age":invalid}`)
	req := &rpcRequest{
		Params: &params,
	}

	ctx := &rpcContext{req: req}
	var data person

	err := ctx.GetRequestBody(&data)
	if err == nil {
		t.Fatal("Expected error for invalid JSON, got nil")
	}

	// Verify it's an rpcParseError
	rpcParseErr, ok := err.(rpcParseError)
	if !ok {
		t.Errorf("Expected rpcParseError, got %T", err)
	}

	// Check if RPCError method works as expected
	rpcErr := rpcParseErr.RPCError()
	if rpcErr.Code != -32602 {
		t.Errorf("Expected code -32602, got: %d", rpcErr.Code)
	}
}

func TestGetRequestBody_NoParams(t *testing.T) {
	// Request with nil params
	req := &rpcRequest{
		Params: nil,
	}

	ctx := &rpcContext{req: req}
	var data person

	err := ctx.GetRequestBody(&data)
	if err == nil {
		t.Fatal("Expected validation error for missing required field, got nil")
	}

	// Verify it's an rpcParseError (due to validation)
	_, ok := err.(rpcParseError)
	if !ok {
		t.Errorf("Expected rpcParseError, got %T", err)
	}
}

func TestGetRequestBody_ValidationFailed(t *testing.T) {
	// Create params with invalid value (required field missing)
	params := json.RawMessage(`{"age":30}`)
	req := &rpcRequest{
		Params: &params,
	}

	ctx := &rpcContext{req: req}
	var data person

	err := ctx.GetRequestBody(&data)
	if err == nil {
		t.Fatal("Expected validation error, got nil")
	}

	// Verify it's an rpcParseError (from validation)
	rpcParseErr, ok := err.(rpcParseError)
	if !ok {
		t.Errorf("Expected rpcParseError, got %T", err)
	}

	// Check error message contains validation info
	if rpcErr := rpcParseErr.RPCError(); !strings.Contains(rpcErr.Message, "required") {
		t.Errorf("Expected validation error to mention 'required' field, got: %s", rpcErr.Message)
	}
}

func TestGetRequestBody_NegativeAgeValidationFailed(t *testing.T) {
	// Create params with invalid value (age negative)
	params := json.RawMessage(`{"name":"John","age":-5}`)
	req := &rpcRequest{
		Params: &params,
	}

	ctx := &rpcContext{req: req}
	var data person

	err := ctx.GetRequestBody(&data)
	if err == nil {
		t.Fatal("Expected validation error, got nil")
	}

	// Verify it's an rpcParseError (from validation)
	rpcParseErr, ok := err.(rpcParseError)
	if !ok {
		t.Errorf("Expected rpcParseError, got %T", err)
	}

	// Check error message contains validation info
	if rpcErr := rpcParseErr.RPCError(); !strings.Contains(rpcErr.Message, "gte") {
		t.Errorf("Expected validation error to mention 'gte', got: %s", rpcErr.Message)
	}
}

func TestGetResponse_Valid(t *testing.T) {
	ctx := &rpcContext{}
	response := person{
		Name: "John",
		Age:  30,
	}

	result, err := ctx.GetResponse(response)
	if err != nil {
		t.Fatalf("GetResponse failed: %v", err)
	}

	// Check response is returned as-is
	typedResult, ok := result.(person)
	if !ok {
		t.Fatalf("Expected result of type testStruct, got %T", result)
	}

	if typedResult.Name != "John" || typedResult.Age != 30 {
		t.Errorf("Response data doesn't match input")
	}
}

func TestGetResponse_Invalid(t *testing.T) {
	ctx := &rpcContext{}
	response := person{
		// Missing required Name field
		Age: 30,
	}

	result, err := ctx.GetResponse(response)
	if err == nil {
		t.Fatal("Expected validation error, got nil")
	}

	if result != nil {
		t.Errorf("Expected nil result, got %v", result)
	}

	if err.Error() != "invalid response from server" {
		t.Errorf("Unexpected error message: %s", err.Error())
	}
}
