package fields

import (
	"fmt"
	"testing"

	"github.com/Malijoe/djanGo-unchained/utils"
)

func TestIntegerField(t *testing.T) {
	field := NewIntegerField(
		Meta{
			Source: "field",
		},
	)
	ivData := map[bool][]interface{}{
		false: {
			3, 2.00, "5", int64(3), int32(5),
		},
		true: {
			[]byte("hello"),
			"string",
			struct{ Random interface{} }{Random: 5},
		},
	}

	for shouldFail, values := range ivData {
		for i, value := range values {
			var pof string
			if shouldFail {
				pof = "Fail"
			} else {
				pof = "Pass"
			}
			t.Run(fmt.Sprintf("IntegerField_ToInternalValue_%s_%02d", pof, i), testField_ToInternalValue(field, value, shouldFail))
			if !shouldFail {
				expected, _ := utils.ParseInt(value)
				t.Run(fmt.Sprintf("IntegerField_ToRepresentation_%s_%02d", pof, i), testField_ToRepresentation(field, expected))
				t.Run(fmt.Sprintf("IntegerField_Internal_%s_%02d", pof, i), testField_Internal(field, expected))
			} else {
				// if an invalid value was attempted the representational value should be 0
				t.Run(fmt.Sprintf("IntegerField_ToRepresentation_%s_%02d", pof, i), testField_ToRepresentation(field, 0))
				t.Run(fmt.Sprintf("IntegerField_Internal_%s_%02d", pof, i), testField_Internal(field, 0))
			}
		}
	}

}
