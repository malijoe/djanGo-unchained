package fields

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"testing"
)

func TestField(t *testing.T) {
	intField := IntegerField{
		Default:  0,
		ReadOnly: true,
		Source:   "field_1",
		Validators: []func(IntegerField) error{
			func(f IntegerField) error {
				if f.InternalValue == 0 {
					return NewFieldError(f.Source, ErrorInvalidValue)
				}
				return nil
			},
			func(f IntegerField) error {
				fmt.Println(f.InternalValue)
				return nil
			},
		},
	}

	stringField := TextField{
		Default:  "hello",
		ReadOnly: false,
		Source:   "field_2",
		Validators: []func(TextField) error{
			func(f TextField) error {
				if f.InternalValue == "there" {
					return NewFieldError(f.Source, ErrorInvalidValue)
				}
				return nil
			},
		},
	}

	boolField := BooleanField{
		Default:   sql.NullBool{},
		WriteOnly: true,
		Source:    "field_3",
		AllowNull: false,
	}

	var test = struct {
		Field1 TextField    `json:"field_1"`
		Field2 IntegerField `json:"field_2"`
		Field3 BooleanField `json:"field_3"`
	}{
		Field1: stringField,
		Field2: intField,
		Field3: boolField,
	}

	data := `
    {
        "field_1": "",
        "field_2": 3,
        "field_3": true
    }
    `
	if err := json.Unmarshal([]byte(data), &test); err != nil {
		t.Error(err)
		return
	}

	if test.Field1.InternalValue != "hello" {
		t.Errorf("field_1 has unexpected value %s", test.Field1.InternalValue)
	}

	if test.Field2.InternalValue != 3 {
		t.Errorf("field_2 has unexpected value %v", test.Field2.InternalValue)
	}
}
