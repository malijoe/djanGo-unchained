package fields

type FieldValidator func(Field) error

type Zeroer interface {
	IsZero() bool
}

var ValidateAllowNullAndRequired FieldValidator = func(field Field) error {

	var (
		isZero = FieldIsNullOrZero(field)
		meta   = field.MetaData()
	)

	if isZero == meta.Required {
		return NewFieldError(meta.Source, ErrorMissingRequiredField)
	}

	if isZero != (meta.AllowNull || meta.AllowBlank) {
		return NewFieldError(meta.Source, ErrorNullNotAllowed)
	}
	return nil
}
