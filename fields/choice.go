package fields

import (
	"fmt"
)

type TextChoiceField struct {
	Options []string
	TextField
}

func NewTextChoiceField(options []string, meta Meta) TextChoiceField {
	meta.Type = Choice
	meta.Validators = append([]FieldValidator{ValidateTextChoice(options)}, meta.Validators...)

	return TextChoiceField{
		Options:   options,
		TextField: NewTextField(meta),
	}
}

var ValidateTextChoice = func(options []string) FieldValidator {
	ops := make([]interface{}, len(options))
	for i := range options {
		ops[i] = options[i]
	}

	return ValidateChoice(ops)
}

type CompareFunc func(to, elem interface{}) bool

var LazyEqual CompareFunc = func(to, elem interface{}) bool {
	ok := to == elem
	return ok
}
var ValidateChoice = func(options []interface{}, compFunc ...CompareFunc) FieldValidator {
	comp := LazyEqual
	if len(compFunc) > 0 {
		comp = compFunc[0]
	}

	return func(field Field) error {
		value := field.Internal()

		for _, option := range options {
			// if a valid option is found return
			if comp(option, value) {
				return nil
			}
		}

		return NewFieldError(field.MetaData().Source, fmt.Errorf("%w valid options are %v", ErrorInvalidValue, options))
	}
}
