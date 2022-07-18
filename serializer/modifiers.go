package serializer

import (
	"github.com/Malijoe/djanGo-unchained/fields"
	"github.com/Malijoe/djanGo-unchained/models"
	"github.com/hashicorp/go-multierror"
)

type FieldModifier func(field fields.Field) error

func ManageModelFields(model models.Model, fields []string, mod FieldModifier) error {
	var errs error
	for _, field := range fields {
		if err := mod(models.GetField(model, field)); err != nil {
			errs = multierror.Append(errs, err)
		}
	}

	return errs
}

func SetFieldReadOnly(field fields.Field) error {
	if field != nil {
		meta := field.MetaData()
		meta.ReadOnly = true
	}
	return nil
}

func SetFieldWriteOnly(field fields.Field) error {
	if field != nil {
		meta := field.MetaData()
		meta.WriteOnly = true
	}
	return nil
}
