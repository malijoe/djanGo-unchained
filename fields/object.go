package fields

import (
	"fmt"
	"reflect"
)

type objectField struct {
	InternalValue interface{}
	ObjType       reflect.Type
	Meta
}

func (f *objectField) MetaData() *Meta {
	return &f.Meta
}

func (f *objectField) ToInternalValue(value interface{}) error {
	f.InternalValue = reflect.Indirect(reflect.New(f.ObjType)).Interface()
	if value != nil {
		v := reflect.ValueOf(value)
		if v.Type() != f.ObjType {
			return NewFieldError(f.Source, fmt.Errorf("%w %v %T", ErrorInvalidValue, value, value))
		}
		f.InternalValue = v.Interface()
	}

	return nil
}

func (f objectField) Internal() interface{} {
	return f.InternalValue
}

func (f objectField) ToRepresentation() interface{} {
	return f.InternalValue
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

	v := reflect.New(f.ObjType).Interface()
	if err := unmarshal(v); err != nil {
		return err
	}

	v = reflect.Indirect(reflect.ValueOf(v))
	return f.ToInternalValue(v)
}

type ObjectField struct {
	*objectField
	FullSerializer
}

func NewObjectField(objType interface{}, meta Meta) ObjectField {
	field := objectField{
		Meta:    meta,
		ObjType: reflect.TypeOf(objType),
	}

	return ObjectField{
		objectField:    &field,
		FullSerializer: DefaultFullSerializer(&field),
	}
}
