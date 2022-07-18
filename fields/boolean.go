package fields

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strconv"
)

// Bool a wrapper for sql.NullBool to marshal/unmarshal YAML and JSON input
type Bool struct {
	sql.NullBool
}

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
	value Bool
	Meta
}

func (f *booleanField) MetaData() *Meta {
	return &f.Meta
}

func (f *booleanField) ToInternalValue(value interface{}) error {
	f.value = Bool{}
	if value != nil {
		switch v := value.(type) {
		case bool:
			f.value.Valid = true
			f.value.Bool = v
		case sql.NullBool:
			f.value.NullBool = v
		case Bool:
			f.value = v
		case string:
			b, err := strconv.ParseBool(v)
			if err != nil {
				return err
			}
			f.value.Valid = true
			f.value.Bool = b
		default:
			return NewFieldError(f.Source, fmt.Errorf("%w %v %T", ErrorInvalidValue, value, value))
		}
	}
	return nil
}

func (f booleanField) Internal() interface{} {
	return f.value
}

func (f booleanField) ToRepresentation() interface{} {
	return f.value
}

func (f booleanField) Marshal() interface{} {
	return f.ToRepresentation()
}

func (f *booleanField) Unmarshal(unmarshal func(interface{}) error) error {
	return unmarshal(&f.value)
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
