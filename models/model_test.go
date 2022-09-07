package models

import (
	"database/sql"
	"fmt"
	"testing"

	"github.com/malijoe/djanGo-unchained/fields"
)

// testModel models.Model implementation for testing
type testModel struct {
	Field1 fields.IntegerField    `json:"field_1"`
	Field2 fields.TextField       `json:"field_2"`
	Field3 fields.TextChoiceField `json:"field_3"`
	Field4 fields.BooleanField    `json:"field_4"`
}

func (t *testModel) Init() {
	*t = testModel{
		Field1: fields.NewIntegerField(
			fields.Meta{
				Source: "field_1",
			},
		),
		Field2: fields.NewTextField(
			fields.Meta{
				Source: "field_2",
			},
		),
		Field3: fields.NewTextChoiceField(
			[]string{"choice_1", "choice_2"},
			fields.Meta{
				Source: "field_3",
			},
		),
		Field4: fields.NewBooleanField(
			fields.Meta{
				Source: "field_4",
			},
		),
	}
}

func (t *testModel) GetFields() []fields.Field {
	return []fields.Field{t.Field1, t.Field2, t.Field3, t.Field4}
}

func TestModel(t *testing.T) {
	model := &testModel{}
	model.Init()
	_ = model.Field1.ToInternalValue(3)
	_ = model.Field2.ToInternalValue("hello")
	_ = model.Field3.ToInternalValue("choice_1")
	_ = model.Field4.ToInternalValue(true)
	type dataDef struct {
		ShouldFind    bool
		ExpectedValue interface{}
	}
	data := map[string]dataDef{
		"field_1": {
			ShouldFind:    true,
			ExpectedValue: 3,
		},
		"field_7": {
			ShouldFind: false,
		},
		"field_3": {
			ShouldFind:    true,
			ExpectedValue: "choice_1",
		},
		"field_8": {
			ShouldFind: false,
		},
		"field_4": {
			ShouldFind: true,
			ExpectedValue: sql.NullBool{
				Valid: true,
				Bool:  true,
			},
		},
	}

	for fieldName, def := range data {
		t.Run(fmt.Sprintf("TestModel_HasField_%s", fieldName), testModel_HasField(model, fieldName, def.ShouldFind))
		if def.ShouldFind {
			t.Run(fmt.Sprintf("TestModel_GetField_%s", fieldName), testModel_GetField(model, fieldName, def.ExpectedValue))
		}
	}
}

func testModel_GetField(model Model, fieldName string, expected interface{}) func(*testing.T) {
	return func(t *testing.T) {
		field := GetField(model, fieldName)
		if field == nil {
			t.Errorf("did not find expected field %s", fieldName)
			return
		}

		if value := field.Internal(); value != expected {
			t.Errorf("field %s did not have expected value %v when retrieved", fieldName, expected)
		}
	}
}

func testModel_HasField(model Model, fieldName string, shouldFind bool) func(*testing.T) {
	return func(t *testing.T) {
		if HasField(model, fieldName) != shouldFind {
			t.Errorf("did not find expected field %s", fieldName)
		}
	}
}
