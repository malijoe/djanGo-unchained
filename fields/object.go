package fields

import (
	"fmt"
	"reflect"
)

type objectField struct {
	value   interface{}
	objType reflect.Type
	Meta
}

func (f *objectField) MetaData() *Meta {
	return &f.Meta
}

func (f *objectField) ToInternalValue(value interface{}) error {
	f.value = reflect.Indirect(reflect.New(f.objType)).Interface()
	if value != nil {
		v := reflect.ValueOf(value)
		if v.Type() != f.objType {
			return NewFieldError(f.Source, fmt.Errorf("%w %v %T", ErrorInvalidValue, value, value))
		}
		f.value = v.Interface()
	}

	return nil
}

func (f objectField) Internal() interface{} {
	return f.value
}

func (f objectField) ToRepresentation() interface{} {
	return f.value
}

func (f objectField) Marshal() interface{} {
	if f.WriteOnly {
		return nil
	}
	return f.ToRepresentation()
}

func (f *objectField) Unmarshal(unmarshal func(interface{}) error) error {
	if f.ReadOnly {
		return nil
	}

	v := reflect.New(f.objType).Interface()
	if err := unmarshal(v); err != nil {
		return err
	}

	v = reflect.Indirect(reflect.ValueOf(v))
	return f.ToInternalValue(v)
}

type ObjectField struct {
	objectField
	FullSerializer
}

func NewObjectField(objType interface{}, meta Meta) *ObjectField {
	field := objectField{
		Meta:    meta,
		objType: reflect.TypeOf(objType),
	}

	return &ObjectField{
		objectField:    field,
		FullSerializer: DefaultFullSerializer(&field),
	}
}
