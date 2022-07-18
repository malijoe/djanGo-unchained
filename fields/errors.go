package fields

import (
	"errors"
	"fmt"
)

type FieldError struct {
	Err       error
	fieldName string
}

func NewFieldError(field string, err ...error) *FieldError {
	e := FieldError{
		fieldName: field,
	}
	if len(err) > 0 {
		e.Err = err[0]
	}
	return &e
}

func (e *FieldError) Error() string {
	message := fmt.Sprintf("error processing field %s", e.fieldName)
	if e.Err != nil {
		message = fmt.Sprintf("%s: %v", message, e.Err)
	}
	return message
}

func (e *FieldError) Is(target error) bool {
	t, ok := target.(*FieldError)
	if !ok {
		return false
	}
	return e.fieldName == t.fieldName
}

func (e *FieldError) Unwrap() error {
	return e.Err
}

var ErrorInvalidValue = errors.New("invalid field value")

var ErrorMissingRequiredField = errors.New("missing required field")

var ErrorNullNotAllowed = errors.New("null is not allowed")
