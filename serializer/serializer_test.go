package serializer

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/Malijoe/djanGo-unchained/fields"
	"github.com/Malijoe/djanGo-unchained/models"
	"github.com/Malijoe/djanGo-unchained/utils"
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

func TestSerializer(t *testing.T) {
	model := testModel{}
	s1 := NewModelSerializer(
		&model,
		Meta{
			ReadOnlyFields:  []string{"field_1", "field_3"},
			WriteOnlyFields: []string{"field_2", "field_4"},
		},
	)

	data := `
        {
            "field_1": 5,
            "field_2": "hello",
            "field_3": "choice_2",
            "field_4": false
        }
    `

	handleMarshalTest := func(s *Serializer) func([]byte) error {
		return func(data []byte) error {
			// make sure the fields have no previous values
			s.Model.Init()
			if err := json.Unmarshal(data, s); err != nil {
				return err
			}

			for _, f := range s.WriteOnlyFields {
				field := models.GetField(s.Model, f)
				value := field.Internal()
				if value != nil && !utils.IsZero(value) {
					return fmt.Errorf("write only field %s has value after marshalling: %v", f, value)
				}
			}
			return nil
		}
	}

	for i, s := range []*Serializer{s1} {
		t.Run(fmt.Sprintf("TestSerializer_UnmarshalJSON_%02d", i), testSerializer_UnmarshalJSON([]byte(data), s))
		t.Run(fmt.Sprintf("TestSerializer_MarshalJSON_%02d", i), testSerializer_MarshalJSON(s, handleMarshalTest(s)))
	}
}

func testSerializer_UnmarshalJSON(data []byte, serializer *Serializer) func(t *testing.T) {
	return func(t *testing.T) {
		if err := json.Unmarshal(data, serializer); err != nil {
			t.Errorf("failed to unmarshal data: %s", err.Error())
			return
		}

		for _, ro := range serializer.ReadOnlyFields {
			field := models.GetField(serializer.Model, ro)
			value := field.Internal()
			if value != nil && !utils.IsZero(value) {
				t.Errorf("read only %s has value after unmarshalling: %v", ro, value)
			}
		}
	}
}

func testSerializer_MarshalJSON(data *Serializer, handle func([]byte) error) func(*testing.T) {
	return func(t *testing.T) {
		json_data, err := json.Marshal(data)
		if err != nil {
			t.Errorf("failed to marshal data: %s", err.Error())
			return
		}

		if err = handle(json_data); err != nil {
			t.Error(err)
		}
	}
}
