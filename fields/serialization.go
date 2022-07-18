package fields

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"

	"gopkg.in/yaml.v2"
)

type Serializable interface {
	Marshal() interface{}
	IsWriteOnly() bool
}

type Serializer interface {
	json.Marshaler
	yaml.Marshaler
}

type defaultFieldSerializer struct {
	Serializable
}

func DefaultFieldSerializer(s Serializable) *defaultFieldSerializer {
	return &defaultFieldSerializer{
		Serializable: s,
	}
}

func (s defaultFieldSerializer) MarshalJSON() ([]byte, error) {
	if s.IsWriteOnly() {
		return json.Marshal(nil)
	}
	obj := s.Marshal()
	return json.Marshal(obj)
}

func (s defaultFieldSerializer) MarshalYAML() (interface{}, error) {
	if s.IsWriteOnly() {
		return nil, nil
	}
	return s.Marshal(), nil
}

type Deserializable interface {
	Unmarshal(func(interface{}) error) error
	IsReadOnly() bool
}

type Deserializer interface {
	json.Unmarshaler
	yaml.Unmarshaler
}

type defaultFieldDeserializer struct {
	Deserializable
}

func DefaultFieldDeserializer(d Deserializable) *defaultFieldDeserializer {
	return &defaultFieldDeserializer{
		Deserializable: d,
	}
}

func (s *defaultFieldDeserializer) UnmarshalJSON(data []byte) error {
	if s.IsReadOnly() {
		return nil
	}
	unmarshal := func(i interface{}) error {
		return json.Unmarshal(data, i)
	}
	return s.Unmarshal(unmarshal)
}

func (s *defaultFieldDeserializer) UnmarshalYAML(unmarshal func(interface{}) error) error {
	if s.IsReadOnly() {
		return nil
	}
	return s.Unmarshal(unmarshal)
}

type DBSerializer interface {
	driver.Valuer
	sql.Scanner
}

type defaultDBSerializer struct {
	Internalizable
}

func DefaultDBSerializer(i Internalizable) *defaultDBSerializer {
	return &defaultDBSerializer{
		Internalizable: i,
	}
}

func (s defaultDBSerializer) Value() (driver.Value, error) {
	i := s.Internal()
	if v, ok := i.(driver.Valuer); ok {
		return v.Value()
	}
	return i, nil
}

func (s *defaultDBSerializer) Scan(value interface{}) error {
	i := s.Internal()
	if scanner, ok := i.(sql.Scanner); ok {
		return scanner.Scan(value)
	}
	return s.ToInternalValue(value)
}

type FullySerializable interface {
	Serializable
	Deserializable
	Internalizable
}

type FullSerializer interface {
	Serializer
	Deserializer
	DBSerializer
}

type defaultFullSerializer struct {
	Serializer
	Deserializer
	DBSerializer
}

func DefaultFullSerializer(f FullySerializable) *defaultFullSerializer {
	return &defaultFullSerializer{
		Serializer:   DefaultFieldSerializer(f),
		Deserializer: DefaultFieldDeserializer(f),
		DBSerializer: DefaultDBSerializer(f),
	}
}
