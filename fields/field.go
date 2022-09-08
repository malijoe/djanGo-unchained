package fields

import (
	"errors"
	"reflect"
)

var (
	SkipField     = errors.New("skip field")
	Required      = errors.New("required")
	InvalidParent = errors.New("invalid parent")
)

// Empty used to represent no data being provided for a given value
// This is required because nil may be a valid value
type Empty any

type Field[T any] struct {
	ReadOnly, WriteOnly, Required bool
	Default, Initial              T
	Label, HelpText               string
	AllowNull                     bool
	_Validators                   []func(T) error
	FieldName                     string
	Source                        string
	value                         *T
	Parent                        any
}

func (f *Field[T]) Bind(field_name string, parent any) error {
	if reflect.Indirect(reflect.ValueOf(parent)).Kind() != reflect.Struct {
		return InvalidParent
	}
	f.Parent = parent
	f.FieldName = field_name
	if f.Source == "" {
		f.Source = f.FieldName
	}
	return nil
}

func (f *Field[T]) Validators() []func(T) error {
	return f._Validators
}

// GetValue given the incoming primitive data, return the value for this field that should be validated and transformed
// to a native value
func (f *Field[T]) GetValue(dictionary map[string]any) any {
	return dictionary[f.FieldName]
}

func (f *Field[T]) GetDefault() (T, error) {
	var zero T
	if reflect.ValueOf(f.Default).IsZero() {
		return zero, SkipField
	}
	return f.Default, nil
}

func (f *Field[T]) Unmarshal(unmarshal func(any) error) error {
	if f.Parent == nil {
		panic("unmarshalling unbound field")
	}
	if f.ReadOnly {
		return nil
	}
	if err := unmarshal(f.Parent); err != nil {
		return err
	}
	return nil
}

func (f *Field[T]) Marshal() any {
	if f.WriteOnly {
		return nil
	}
	return *f.value
}
