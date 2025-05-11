package validator

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

var V = validator.New(validator.WithRequiredStructEnabled())

func ValidateUnion(sl validator.StructLevel) {
	// Get current value and its type information
	current := sl.Current()
	t := current.Type()
	structName := t.Name()

	n := 0
	// Iterate through struct fields
	for i := range current.NumField() {
		field := t.Field(i)

		// Skip fields not starting with "Of"
		if !strings.HasPrefix(field.Name, "Of") {
			panic(fmt.Sprintf("Field %s in %s must start with 'Of'", field.Name, structName))
		}

		// Ensure field is a pointer type
		if field.Type.Kind() != reflect.Ptr {
			panic(fmt.Sprintf("Field %s in %s must be a pointer type", field.Name, structName))
		}

		// Count non-nil fields
		if !current.Field(i).IsNil() {
			if n++; n > 1 {
				break
			}
		}
	}

	// Report error if not exactly one field is set
	if n != 1 {
		sl.ReportError(current.Interface(), structName, structName, "oneOfRequired", "")
	}
}
