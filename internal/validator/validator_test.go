package validator

import (
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"github.com/umk/llmservices/internal/pointer"
)

// TestValidateUnion tests the ValidateUnion function
func TestValidateUnion(t *testing.T) {
	// Create a new validator instance for testing
	v := validator.New()

	// Register our custom validator
	v.RegisterStructValidation(ValidateUnion, TestUnion{})

	tests := []struct {
		name        string
		input       TestUnion
		expectError bool
	}{
		{
			name: "valid - one field set",
			input: TestUnion{
				OfString: pointer.Ptr("test"),
				OfInt:    nil,
				OfBool:   nil,
			},
			expectError: false,
		},
		{
			name: "invalid - no fields set",
			input: TestUnion{
				OfString: nil,
				OfInt:    nil,
				OfBool:   nil,
			},
			expectError: true,
		},
		{
			name: "invalid - multiple fields set",
			input: TestUnion{
				OfString: pointer.Ptr("test"),
				OfInt:    pointer.Ptr(42),
				OfBool:   nil,
			},
			expectError: true,
		},
		{
			name: "valid - different field set",
			input: TestUnion{
				OfString: nil,
				OfInt:    nil,
				OfBool:   pointer.Ptr(true),
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := v.Struct(tt.input)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}

	// Test the panic cases separately
	t.Run("panic - field without Of prefix", func(t *testing.T) {
		v := validator.New()
		v.RegisterStructValidation(ValidateUnion, InvalidPrefixUnion{})

		defer func() {
			r := recover()
			assert.NotNil(t, r)
			assert.Contains(t, r.(string), "must start with 'Of'")
		}()

		_ = v.Struct(InvalidPrefixUnion{})
	})

	t.Run("panic - field not pointer type", func(t *testing.T) {
		v := validator.New()
		v.RegisterStructValidation(ValidateUnion, NonPointerUnion{})

		defer func() {
			r := recover()
			assert.NotNil(t, r)
			assert.Contains(t, r.(string), "must be a pointer type")
		}()

		_ = v.Struct(NonPointerUnion{})
	})
}

// Test structures
type TestUnion struct {
	OfString *string
	OfInt    *int
	OfBool   *bool
}

type InvalidPrefixUnion struct {
	OfString *string
	BadName  *string // Name doesn't start with "Of"
}

type NonPointerUnion struct {
	OfString *string
	OfInt    int // Not a pointer type
}
