package fields

import (
	"database/sql"
	"fmt"
	"strconv"
	"testing"
)

func TestBooleanField(t *testing.T) {
	field := NewBooleanField(
		Meta{
			Source: "field",
		},
	)

	data := map[bool][]interface{}{
		false: {
			true,
			sql.NullBool{
				Valid: false,
			},
			Bool{
				NullBool: sql.NullBool{
					Valid: true,
					Bool:  false,
				},
			},
			"true",
		},
		true: {
			"truth",
			"hi",
			struct {
				Internal bool
			}{Internal: true},
		},
	}

	//defaultValue := Bool{}
	toExpectedValue := func(value interface{}) interface{} {
		expected := Bool{}
		switch v := value.(type) {
		case bool:
			expected.Valid = true
			expected.Bool = v
		case sql.NullBool:
			expected.NullBool = v
		case Bool:
			expected = v
		case string:
			b, err := strconv.ParseBool(v)
			if err == nil {
				expected.Valid = true
				expected.Bool = b
			}
		}

		return expected
	}
	for shouldFail, values := range data {
		for i, value := range values {
			expectedValue := toExpectedValue(value)
			t.Run(fmt.Sprintf("BooleanField_ToInternalValue_%s_%02d", getPoF(shouldFail), i), testField_ToInternalValue(field, value, shouldFail))
			t.Run(fmt.Sprintf("BooleanField_ToRepresentation_%s_%02d", getPoF(shouldFail), i), testField_ToRepresentation(field, expectedValue))
			t.Run(fmt.Sprintf("BooleanField_Internal_%s_%02d", getPoF(shouldFail), i), testField_Internal(field, expectedValue))
		}
	}
}
