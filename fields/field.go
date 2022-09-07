package fields

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"reflect"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/lib/pq"
)

type Internalizable interface {
	ToInternalValue(interface{}) error
	Internal() interface{}
}

type Representable interface {
	ToRepresentation() interface{}
}

type field interface {
	Internalizable
	Representable
	MetaData() *Meta
}

type Descriptor[T any] Field[T]

type Field[T any] struct {
	InternalValue    T                      `json:"-"`
	Default          T                      `json:"default"`
	ReadOnly         bool                   `json:"read_only"`
	WriteOnly        bool                   `json:"write_only"`
	Required         bool                   `json:"required"`
	AllowNull        bool                   `json:"allow_null"`
	AllowBlank       bool                   `json:"allow_blank"`
	Source           string                 `json:"source"`
	Label            string                 `json:"label"`
	HelpText         string                 `json:"help_text"`
	Type             FieldType              `json:"type"`
	PlaceHolder      string                 `json:"place_holder"`
	DependantFields  []string               `json:"dependant_fields"`
	Validators       []func(Field[T]) error `json:"-"`
	ToRepresentation func(T) interface{}    `json:"-"`
}

func (f *Field[T]) Unmarshal(unmarshal func(interface{}) error) error {
	if f.WriteOnly {
		return nil
	}

	if err := unmarshal(&f.InternalValue); err != nil {
		return err
	}

	isZero := reflect.ValueOf(f.InternalValue).IsZero()
	if isZero {
		var zero T
		if reflect.DeepEqual(zero, f.Default) {
			if !(f.AllowNull) || f.AllowBlank {
				return NewFieldError(f.Source, ErrorNullNotAllowed)
			}
		} else {
			f.InternalValue = f.Default
		}
	}

	var errs error
	for _, validator := range f.Validators {
		if err := validator(*f); err != nil {
			errs = multierror.Append(errs, err)
		}
	}

	return errs
}

func (f *Field[T]) Marshal() interface{} {
	if f.WriteOnly {
		return nil
	}
	if f.ToRepresentation != nil {
		return f.ToRepresentation(f.InternalValue)
	}

	return f.InternalValue
}

func (f *Field[T]) UnmarshalJSON(data []byte) error {
	unmarshal := func(i interface{}) error {
		return json.Unmarshal(data, i)
	}
	return f.Unmarshal(unmarshal)
}

func (f Field[T]) MarshalJSON() ([]byte, error) {
	return json.Marshal(f.Marshal())
}

func (f *Field[T]) UnmarshalYAML(unmarshal func(interface{}) error) error {
	return f.Unmarshal(unmarshal)
}

func (f Field[T]) MarshalYAML() (interface{}, error) {
	return f.Marshal(), nil
}

func (f Field[T]) Scan(value interface{}) error {
	if scanner, ok := any(f.InternalValue).(sql.Scanner); ok {
		return scanner.Scan(value)
	}

	v, ok := value.(T)
	if !ok {
		return NewFieldError(f.Source, fmt.Errorf("%w %v %T", ErrorInvalidValue, value, value))
	}

	f.InternalValue = v
	return nil
}

func (f Field[T]) Value() (driver.Value, error) {
	if valuer, ok := any(f.InternalValue).(driver.Valuer); ok {
		return valuer.Value()
	}

	return f.InternalValue, nil
}

type ChoiceField[T any] struct {
	Options []T
	Field[T]
}

type TextField = Field[string]
type IntegerField = Field[int]
type BooleanField = Field[sql.NullBool]
type TimeField = Field[time.Time]
type TextChoiceField = ChoiceField[string]
type IntegerChoiceField = ChoiceField[int]
type TextArrayField = Field[pq.StringArray]
