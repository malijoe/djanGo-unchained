package models

import "github.com/malijoe/djanGo-unchained/fields"

type Model interface {
	Init()
	GetFields() []fields.Field
}

func GetField(model Model, fieldName string) fields.Field {
	for _, field := range model.GetFields() {
		field_meta := field.MetaData()
		if field_meta.Source == fieldName {
			return field
		}
	}

	return nil
}

func HasField(model Model, fieldName string) bool {
	return GetField(model, fieldName) != nil
}
