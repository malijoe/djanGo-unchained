package fields

import "fmt"

type textField struct {
	value string
	Meta
}

func (f *textField) MetaData() *Meta {
	return &f.Meta
}

func (f *textField) ToInternalValue(value interface{}) error {
	f.value = ""
	if value != nil {
		switch v := value.(type) {
		case string:
			f.value = v
		case []byte:
			f.value = string(v)
		default:
			return NewFieldError(f.Source, fmt.Errorf("%w %v %T", ErrorInvalidValue, value, value))
		}
	}
	return nil
}

func (f textField) Internal() interface{} {
	return f.value
}

func (f textField) ToRepresentation() interface{} {
	return f.value
}

func (f *textField) Unmarshal(unmarshal func(interface{}) error) error {
	if f.ReadOnly {
		return nil
	}

	return unmarshal(&f.value)
}

func (f textField) Marshal() interface{} {
	if f.WriteOnly {
		return nil
	}
	return f.ToRepresentation()
}

func (f *textField) UnmarshalParam(param string) error {
	f.value = param
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
