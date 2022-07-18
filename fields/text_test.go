package fields

import (
	"fmt"
	"testing"
)

func TestTextField(t *testing.T) {
	field := NewTextField(
		Meta{
			Source: "field",
		},
	)

	data := map[bool][]interface{}{
		false: {
			"hello",
			"there",
			[]byte("General Kenobi"),
		},
		true: {
			5,
			16,
			struct {
				Internal string
			}{Internal: "hello world"},
		},
	}

	defaultValue := ""
	for shouldFail, values := range data {
		for i, value := range values {
			var expectedValue interface{}
			switch v := value.(type) {
			case []byte:
				expectedValue = string(v)
			default:
				expectedValue = v
			}
			if shouldFail {
				expectedValue = defaultValue
			}

			t.Run(fmt.Sprintf("TextField_ToInternalValue_%s_%02d", getPoF(shouldFail), i), testField_ToInternalValue(field, value, shouldFail))
			t.Run(fmt.Sprintf("TextField_ToRepresentation_%s_%02d", getPoF(shouldFail), i), testField_ToRepresentation(field, expectedValue))
			t.Run(fmt.Sprintf("TextField_Internal_%s_%02d", getPoF(shouldFail), i), testField_Internal(field, expectedValue))
		}
	}
}
