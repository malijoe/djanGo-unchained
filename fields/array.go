package fields

import (
	"fmt"
	"reflect"

	"github.com/lib/pq"
)

type stringArrayField struct {
	values pq.StringArray
	Meta
}

func (f *stringArrayField) MetaData() *Meta {
	return &f.Meta
}

func (f *stringArrayField) ToInternalValue(value interface{}) error {
	f.values = pq.StringArray{}
	if value != nil {
		switch v := value.(type) {
		case pq.StringArray:
			f.values = v
		case []string:
			f.values = v
		case []byte, string:
			if err := f.values.Scan(value); err != nil {
				return err
			}
		default:
			return NewFieldError(f.Source, fmt.Errorf("%w %v %T", ErrorInvalidValue, value, value))
		}
	}
	return nil
}

func (f stringArrayField) Internal() interface{} {
	return f.values
}

func (f stringArrayField) ToRepresentation() interface{} {
	return []string(f.values)
}

func (f stringArrayField) Marshal() interface{} {
	return f.ToRepresentation()
}

func (f *stringArrayField) Unmarshal(unmarshal func(interface{}) error) error {
	return unmarshal(&f.values)
}

type StringArrayField struct {
	*stringArrayField
	FullSerializer
}

func NewStringArrayField(meta Meta) StringArrayField {
	field := stringArrayField{
		Meta: meta,
	}

	return StringArrayField{
		stringArrayField: &field,
		FullSerializer:   DefaultFullSerializer(&field),
	}
}

type objectArrayField struct {
	values   []interface{}
	elemType reflect.Type
	Meta
}

func (f *objectArrayField) MetaData() *Meta {
	return &f.Meta
}

func (f *objectArrayField) ToInternalValue(value interface{}) error {
	f.values = []interface{}{}
	if value != nil {
		v := reflect.ValueOf(value)
		if !(v.Kind() == reflect.Array || v.Kind() == reflect.Slice) || v.Type().Elem() != f.elemType {
			return NewFieldError(f.Source, fmt.Errorf("%w %v %T", ErrorInvalidValue, value, value))
		}

		values := make([]interface{}, v.Len())
		for i := 0; i < v.Len(); i++ {
			obj := v.Index(i)
			values[i] = obj.Interface()
		}
		f.values = values
	}
	return nil
}

func (f objectArrayField) Internal() interface{} {
	sl := reflect.New(reflect.SliceOf(f.elemType)).Elem()
	for _, v := range f.values {
		sl = reflect.Append(sl, reflect.ValueOf(v))
	}

	return sl.Interface()
}

func (f objectArrayField) ToRepresentation() interface{} {
	sl := reflect.New(reflect.SliceOf(f.elemType)).Elem()
	for _, v := range f.values {
		sl = reflect.Append(sl, reflect.ValueOf(v))
	}
	return sl.Interface()
}

func (f objectArrayField) Marshal() interface{} {
	return f.ToRepresentation()
}

func (f *objectArrayField) Unmarshal(unmarshal func(interface{}) error) error {
	sl := reflect.New(reflect.SliceOf(f.elemType)).Interface()
	if err := unmarshal(sl); err != nil {
		return err
	}
	return f.ToInternalValue(sl)
}

type ObjectArrayField struct {
	*objectArrayField
	FullSerializer
}

func NewObjectArrayField(elem interface{}, meta Meta) ObjectArrayField {
	field := objectArrayField{
		Meta:     meta,
		elemType: reflect.TypeOf(elem),
	}
	return ObjectArrayField{
		objectArrayField: &field,
		FullSerializer:   DefaultFullSerializer(&field),
	}
}
