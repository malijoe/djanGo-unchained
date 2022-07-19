package fields

import (
	"fmt"
	"time"
)

type timeField struct {
	InternalValue time.Time
	Meta
}

func (f *timeField) MetaData() *Meta {
	return &f.Meta
}

func (f *timeField) ToInternalValue(value interface{}) error {
	f.InternalValue = time.Time{}
	if value != nil {
		v, ok := value.(time.Time)
		if !ok {
			return NewFieldError(f.Source, fmt.Errorf("%w %v %T", ErrorInvalidValue, value, value))
		}
		f.InternalValue = v
	}

	return nil
}

func (f timeField) Internal() interface{} {
	return f.InternalValue
}

func (f timeField) ToRepresentation() interface{} {
	return f.InternalValue
}

func (f *timeField) Unmarshal(unmarshal func(interface{}) error) error {
	if err := unmarshal(&f.InternalValue); err != nil {
		return err
	}

	return nil
}

func (f timeField) Marshal() interface{} {
	return f.ToRepresentation()
}

type TimeField struct {
	*timeField
	FullSerializer
}

func NewTimeField(meta Meta) TimeField {
	field := timeField{
		Meta: meta,
	}

	return TimeField{
		timeField:      &field,
		FullSerializer: DefaultFullSerializer(&field),
	}
}
