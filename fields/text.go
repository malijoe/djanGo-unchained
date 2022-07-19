package fields

import "fmt"

type textField struct {
	InternalValue string
	Meta
}

func (f *textField) MetaData() *Meta {
	return &f.Meta
}

func (f *textField) ToInternalValue(value interface{}) error {
	f.InternalValue = ""
	if value != nil {
		switch v := value.(type) {
		case string:
			f.InternalValue = v
		case []byte:
			f.InternalValue = string(v)
		default:
			return NewFieldError(f.Source, fmt.Errorf("%w %v %T", ErrorInvalidValue, value, value))
		}
	}
	return nil
}

func (f textField) Internal() interface{} {
	return f.InternalValue
}

func (f textField) ToRepresentation() interface{} {
	return f.InternalValue
}

func (f *textField) Unmarshal(unmarshal func(interface{}) error) error {
	if f.ReadOnly {
		return nil
	}

	return unmarshal(&f.InternalValue)
}

func (f textField) Marshal() interface{} {
	if f.WriteOnly {
		return nil
	}
	return f.ToRepresentation()
}

func (f *textField) UnmarshalParam(param string) error {
	f.InternalValue = param
	return nil
}

type TextField struct {
	*textField
	FullSerializer
}

func NewTextField(meta Meta) TextField {
	field := textField{
		Meta: meta,
	}

	return TextField{
		textField:      &field,
		FullSerializer: DefaultFullSerializer(&field),
	}
}
