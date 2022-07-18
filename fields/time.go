package fields

import (
	"fmt"
	"time"
)

type timeField struct {
	value time.Time
	Meta
}

func (f *timeField) MetaData() *Meta {
	return &f.Meta
}

func (f *timeField) ToInternalValue(value interface{}) error {
	f.value = time.Time{}
	if value != nil {
		v, ok := value.(time.Time)
		if ok {
			return NewFieldError(f.Source, fmt.Errorf("%w %v %T", ErrorInvalidValue, value, value))
		}
		f.value = v
	}

	return nil
}

func (f timeField) Internal() interface{} {
	return f.value
}

func (f timeField) ToRepresentation() interface{} {
	return f.value
}

func (f *timeField) Unmarshal(unmarshal func(interface{}) error) error {
	if err := unmarshal(&f.value); err != nil {
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
