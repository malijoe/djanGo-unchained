package django

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

var (
	SkipField            = errors.New("skip field")
	Required             = errors.New("required")
	InvalidParent        = errors.New("invalid parent")
	InvalidFieldType     = errors.New("invalid field type")
	NotReadOnlyWriteOnly = errors.New("may not set both `read_only` and `write_only`")
	NotReadOnlyRequired  = errors.New("may not set both `read_only` and `required`")
	NotRequiredDefault   = errors.New("may not set both `required` and `default`")
)

// Empty used to represent no data being provided for a given value
// This is required because nil may be a valid value
type Empty any

type ErrorMessage string

func (e ErrorMessage) Format(args ...any) error {
	if !strings.Contains(string(e), "%") || len(args) == 0 {
		return e
	}
	return fmt.Errorf(string(e), args...)
}

func (e ErrorMessage) Error() string {
	return string(e)
}

var (
	MissingErrorMessage ErrorMessage = "error raised by %s, but error key %s does not exist in the `error_messages` dictionary"
)

func isZero(v any) bool {
	return reflect.Indirect(reflect.ValueOf(v)).IsZero()
}

type FieldMeta interface {
	GetFields() []Field
}

type Field interface {
	Init() error
	Bind(field_name string, parent any) error
	Unmarshal(unmarshal func(any) error) error
	Marshal() any
	Fail(key string, args ...any) error
}

type BaseField[T any] struct {
	ReadOnly, WriteOnly, Required bool
	Default, Initial              T
	Label, HelpText               string
	AllowNull                     bool
	_Validators                   []func(T) error
	FieldName                     string
	Source                        string
	value                         *T
	Parent                        any
	DefaultErrorMessages          map[string]ErrorMessage
}

func (f *BaseField[T]) zero() T {
	var zero T
	return zero
}

func (f *BaseField[T]) Init() error {
	if f.ReadOnly && f.WriteOnly {
		return NotReadOnlyWriteOnly
	}
	if f.ReadOnly && f.Required {
		return NotReadOnlyRequired
	}
	if f.ReadOnly && !isZero(f.Default) {
		return NotRequiredDefault
	}
	default_error_messages := map[string]ErrorMessage{
		"required": "this field is required",
		"null":     "this field may not be null",
	}
	if f.DefaultErrorMessages == nil {
		f.DefaultErrorMessages = default_error_messages
	}
	for k, v := range default_error_messages {
		if _, ok := f.DefaultErrorMessages[k]; !ok {
			f.DefaultErrorMessages[k] = v
		}
	}
	return nil
}

func (f *BaseField[T]) Bind(field_name string, parent any) error {
	p := reflect.Indirect(reflect.ValueOf(parent))
	if p.Kind() != reflect.Struct {
		return InvalidParent
	}
	f.Parent = parent
	f.FieldName = field_name
	if f.Source == "" {
		f.Source = f.FieldName
	}

	v, ok := p.FieldByName(field_name).Interface().(T)
	if !ok {
		return InvalidFieldType
	}
	f.value = &v
	return nil
}

func (f *BaseField[T]) Validators() []func(T) error {
	return f._Validators
}

// GetValue given the incoming primitive data, return the value for this field that should be validated and transformed
// to a native value
func (f *BaseField[T]) GetValue(dictionary map[string]any) any {
	return dictionary[f.FieldName]
}

func (f *BaseField[T]) GetDefault() (T, error) {
	var zero T
	if reflect.ValueOf(f.Default).IsZero() {
		return zero, SkipField
	}
	return f.Default, nil
}

func (f *BaseField[T]) Unmarshal(unmarshal func(any) error) error {
	if f.Parent == nil {
		panic("unmarshalling unbound field")
	}
	if f.ReadOnly {
		return nil
	}
	if err := unmarshal(f.Parent); err != nil {
		return err
	}
	return nil
}

func (f *BaseField[T]) Marshal() any {
	if f.WriteOnly {
		return nil
	}
	return *f.value
}

func (f *BaseField[T]) Fail(key string, args ...any) error {
	msg, ok := f.DefaultErrorMessages[key]
	if !ok {
		return MissingErrorMessage.Format(f.FieldName, key)
	}
	return msg.Format(args...)
}

type BooleanField struct {
	BaseField[bool]
}

func (f *BooleanField) Init() error {
	f.Initial = false
	f.DefaultErrorMessages = map[string]ErrorMessage{
		"invalid": "must be a valid boolean",
	}
	return f.BaseField.Init()
}

type CharField struct {
	trimWhiteSpace, trimSet bool
	MaxLength, MinLength    int
	BaseField[string]
}

func (f *CharField) TrimWhiteSpace(b bool) *CharField {
	f.trimWhiteSpace = b
	f.trimSet = true
	return f
}

func (f *CharField) Init() error {
	f.Initial = ""
	if !f.trimSet && !f.trimWhiteSpace {
		f.trimWhiteSpace = true
	}
	f.DefaultErrorMessages = map[string]ErrorMessage{
		"invalid":    "not a valid string",
		"blank":      "this field may not be blank",
		"max_length": "ensure this field has no more than %d characters",
		"min_length": "ensure this field has at least %d characters",
	}
	return f.BaseField.Init()
}

type IntegerField struct {
	BaseField[int]
	MaxValue, MinValue int
}

func (f *IntegerField) Init() error {
	f.Initial = 0
	f.DefaultErrorMessages = map[string]ErrorMessage{
		"invalid":   "a valid integer is required",
		"max_value": "ensure this value is less than or equal to %d",
		"min_value": "ensure this value is greater tha or equal to %d",
	}
	return f.BaseField.Init()
}

type ChoiceField[T any] struct {
	BaseField[T]
	Choices []T
}
