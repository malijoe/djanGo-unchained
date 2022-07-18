package fields

import "fmt"

type FieldType uint32

const (
	Invalid FieldType = iota
	String
	Boolean
	Choice
	Integer
	ID
	DateTime
	Object
)

func (t FieldType) String() string {
	switch t {
	case String:
		return "string"
	case Boolean:
		return "boolean"
	case Choice:
		return "choice"
	case Integer:
		return "integer"
	case ID:
		return "id"
	case DateTime:
		return "datetime"
	case Object:
		return "object"
	default:
		return "invalid"
	}
}

func ParseFieldType(s string) FieldType {
	switch s {
	case String.String():
		return String
	case Boolean.String():
		return Boolean
	case Choice.String():
		return Choice
	case Integer.String():
		return Integer
	case ID.String():
		return ID
	case DateTime.String():
		return DateTime
	case Object.String():
		return Object
	default:
		return Invalid
	}
}

func (t FieldType) MarshalJSON() ([]byte, error) {
	return []byte(t.String()), nil
}

func (t *FieldType) UnmarshalJSON(data []byte) error {
	ft := ParseFieldType(string(data))
	if ft == Invalid {
		return fmt.Errorf("%s is not a valid field_type", string(data))
	}
	*t = ft
	return nil
}
