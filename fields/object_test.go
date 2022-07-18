package fields

import (
	"fmt"
	"testing"
)

func TestObjectField(t *testing.T) {
	// type for testing object field methods
	type testType struct {
		Internal int
	}
	field := NewObjectField(
		testType{},
		Meta{
			Source: "field",
		},
	)

	testData := map[bool][]interface{}{
		false: {
			testType{
				Internal: 5,
			},
			testType{
				Internal: 10,
			},
			testType{
				Internal: 15,
			},
		},
		true: {
			false,
			5,
			struct {
				Internal int
			}{Internal: 5},
		},
	}
	defaultValue := testType{}
	for shouldFail, values := range testData {
		for i, value := range values {
			expectedValue := value
			if shouldFail {
				expectedValue = defaultValue
			}
			t.Run(fmt.Sprintf("ObjectField_ToInternalValue_%s_%02d", getPoF(shouldFail), i), testField_ToInternalValue(field, value, shouldFail))
			t.Run(fmt.Sprintf("ObjectField_ToRepresentation_%s_%02d", getPoF(shouldFail), i), testField_ToRepresentation(field, expectedValue))
			t.Run(fmt.Sprintf("ObjectField_Internal_%s_%02d", getPoF(shouldFail), i), testField_Internal(field, expectedValue))
		}
	}
}
