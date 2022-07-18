package fields

import (
	"fmt"
	"testing"

	"github.com/lib/pq"
)

func TestStringArrayField(t *testing.T) {
	field := NewStringArrayField(
		Meta{
			Source: "field",
		},
	)

	data := map[bool][]interface{}{
		false: {
			[]string{
				"hello",
				"there",
			},
			pq.StringArray{
				"hello",
				"there",
				"general",
			},
			`{ "hello", "there", "general" , "kenobi" }`,
		},
		true: {
			"hello, there, general, kenobi",
			struct {
				Internal []string
			}{Internal: []string{"hello", "there"}},
			[]byte("general"),
		},
	}

	toExpectedValue := func(value interface{}, represent bool) interface{} {
		expected := pq.StringArray{}
		switch v := value.(type) {
		case []string:
			expected = v
		case pq.StringArray:
			expected = v
		case []byte, string:
			_ = expected.Scan(v)
		}
		if represent {
			return []string(expected)
		}
		return expected
	}

	for shouldFail, values := range data {
		for i, value := range values {
			t.Run(fmt.Sprintf("StringArrayField_ToInternalValue_%s_%02d", getPoF(shouldFail), i), testField_ToInternalValue(field, value, shouldFail))
			t.Run(fmt.Sprintf("StringArrayField_ToRepresentation_%s_%02d", getPoF(shouldFail), i), testField_ToRepresentation(field, toExpectedValue(value, true)))
			t.Run(fmt.Sprintf("StringArrayField_Internal_%s_%02d", getPoF(shouldFail), i), testField_Internal(field, toExpectedValue(value, false)))
		}
	}
}

func TestObjectArrayField(t *testing.T) {
	type testType struct {
		Internal string
	}
	field := NewObjectArrayField(
		testType{},
		Meta{
			Source: "field",
		},
	)

	data := map[bool][]interface{}{
		false: {
			[]testType{
				{Internal: "hello"},
				{Internal: "there"},
				{Internal: "general"},
				{Internal: "kenobi"},
			},
		},
		true: {
			[]string{"hello", "there", "general", "kenobi"},
		},
	}

	toExpectedValue := func(value interface{}) interface{} {
		var expected []testType
		switch v := value.(type) {
		case []testType:
			expected = v
		}
		return expected
	}

	for shouldFail, values := range data {
		for i, value := range values {
			t.Run(fmt.Sprintf("StringArrayField_ToInternalValue_%s_%02d", getPoF(shouldFail), i), testField_ToInternalValue(field, value, shouldFail))
			t.Run(fmt.Sprintf("StringArrayField_ToRepresentation_%s_%02d", getPoF(shouldFail), i), testField_ToRepresentation(field, toExpectedValue(value)))
			t.Run(fmt.Sprintf("StringArrayField_Internal_%s_%02d", getPoF(shouldFail), i), testField_Internal(field, toExpectedValue(value)))
		}
	}

}
