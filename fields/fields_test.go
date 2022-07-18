package fields

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"

	"gopkg.in/yaml.v2"
)

func getPoF(shouldFail bool) string {
	if shouldFail {
		return "Fail"
	}
	return "Pass"
}

func testField_ToInternalValue(field Field, value interface{}, shouldFail bool) func(*testing.T) {
	return func(t *testing.T) {
		err := field.ToInternalValue(value)
		if err != nil != shouldFail {
			if shouldFail {
				t.Errorf("expected to fail when setting field %T to value %v %T", field, value, value)
			} else {
				t.Errorf("failed when setting field %T to value %v %T: %s", field, value, value, err.Error())
			}
		}
	}
}

func testField_ToRepresentation(field Field, expected interface{}) func(*testing.T) {
	return func(t *testing.T) {
		value := field.ToRepresentation()
		if !compareValues(value, expected) {
			t.Errorf("expected field value representation %v but got %v instead", expected, value)
		}
	}
}

func testField_Internal(field Field, expected interface{}) func(*testing.T) {
	return func(t *testing.T) {
		value := field.Internal()
		if !compareValues(value, expected) {
			t.Errorf("expected field internal value %v but got %v instead", expected, value)
		}
	}
}

func compareValues(value, expected interface{}) bool {
	return reflect.DeepEqual(value, expected)
}

func testField_UnmarshalJSON(value interface{}, handler func(func(interface{}) error) error) error {
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("error marshalling json data: %w", err)
	}

	unmarshal := func(i interface{}) error {
		return json.Unmarshal(data, i)
	}

	return handler(unmarshal)
}

func testField_UnmarshalYAML(value interface{}, handler func(func(interface{}) error) error) error {
	data, err := yaml.Marshal(value)
	if err != nil {
		return fmt.Errorf("error marshalling yaml data: %w", err)
	}

	unmarshal := func(i interface{}) error {
		return yaml.Unmarshal(data, i)
	}

	return handler(unmarshal)
}
