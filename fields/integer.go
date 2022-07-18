package fields

import (
	"fmt"
	"strconv"

	"github.com/malijoe/djanGo-unchained/utils"
)

type integerField struct {
	value int
	Meta
}

func (f *integerField) MetaData() *Meta {
	return &f.Meta
}

func (f *integerField) ToInternalValue(value interface{}) error {
	// TODO: add function to reset field value to it's default or zero value
	// reset the internal value to avoid conflicts
	f.value = 0
	if value != nil {
		v, ok := utils.ParseInt(value)
		if !ok {
			return NewFieldError(f.Source, fmt.Errorf("%w %v %T", ErrorInvalidValue, value, value))
		}
		f.value = v
	}

	return nil
}

func (f integerField) ToRepresentation() interface{} {
	return f.value
}

func (f integerField) Internal() interface{} {
	return f.value
}

func (f integerField) Marshal() interface{} {
	return f.ToRepresentation()
}

func (f *integerField) Unmarshal(unmarshal func(interface{}) error) error {
	return unmarshal(&f.value)
}

func (f *integerField) UnmarshalParam(param string) error {
	value, err := strconv.Atoi(param)
	if err != nil {
		return err
	}
	f.value = value
	return nil
}

type IntegerField struct {
	*integerField
	FullSerializer
}

func NewIntegerField(meta Meta) IntegerField {
	field := &integerField{
		Meta: meta,
	}
	return IntegerField{
		integerField:   field,
		FullSerializer: DefaultFullSerializer(field),
	}
}
