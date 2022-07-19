package fields

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strconv"
)

// Bool a wrapper for sql.NullBool to marshal/unmarshal YAML and JSON input
type Bool sql.NullBool

func (b *Bool) unmarshal(unmarshal func(interface{}) error) error {
	var ok bool
	s := &ok
	if err := unmarshal(&s); err != nil {
		return err
	}

	b.Valid = s != nil
	if b.Valid {
		b.Bool = *s
	}
	return nil
}

func (b *Bool) UnmarshalYAML(unmarshal func(interface{}) error) error {
	return b.unmarshal(unmarshal)
}

func (b *Bool) UnmarshalJSON(data []byte) error {
	unmarshal := func(i interface{}) error {
		return json.Unmarshal(data, i)
	}
	return b.unmarshal(unmarshal)
}

func (b Bool) marshal() interface{} {
	if b.Valid {
		return b.Bool
	}
	return nil
}

func (b Bool) MarshalJSON() ([]byte, error) {
	return json.Marshal(b.marshal())
}

func (b Bool) MarshalYAML() (interface{}, error) {
	return b.marshal(), nil
}

type booleanField struct {
	InternalValue sql.NullBool
	Meta
}

func (f *booleanField) MetaData() *Meta {
	return &f.Meta
}

func (f *booleanField) ToInternalValue(value interface{}) error {
	f.InternalValue = sql.NullBool{}
	if value != nil {
		switch v := value.(type) {
		case bool:
			f.InternalValue.Valid = true
			f.InternalValue.Bool = v
		case sql.NullBool:
			f.InternalValue = v
		case string:
			b, err := strconv.ParseBool(v)
			if err != nil {
				return err
			}
			f.InternalValue.Valid = true
			f.InternalValue.Bool = b
		default:
			return NewFieldError(f.Source, fmt.Errorf("%w %v %T", ErrorInvalidValue, value, value))
		}
	}
	return nil
}

func (f booleanField) Internal() interface{} {
	return f.InternalValue
}

func (f booleanField) ToRepresentation() interface{} {
	return Bool(f.InternalValue)
}

func (f booleanField) Marshal() interface{} {
	return f.ToRepresentation()
}

func (f *booleanField) Unmarshal(unmarshal func(interface{}) error) error {
	v := Bool(f.InternalValue)
	return unmarshal(&v)
}

type BooleanField struct {
	*booleanField
	FullSerializer
}

func NewBooleanField(meta Meta) BooleanField {
	field := &booleanField{
		Meta: meta,
	}

	return BooleanField{
		booleanField:   field,
		FullSerializer: DefaultFullSerializer(field),
	}
}
