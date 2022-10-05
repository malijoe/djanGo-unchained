package django

import "reflect"

type (
	InstanceBuilder interface {
		AddField(name string, typ any, tag string) InstanceBuilder
		RemoveField(name string) InstanceBuilder
		HasField(name string) bool
		GetField(name string) InstanceField
		Build() Instance
	}

	InstanceField interface {
		SetType(typ any) InstanceField
		SetTag(tag string) InstanceField
	}

	Instance interface {
		New() any
	}
	builder struct {
		fields []*field
	}
	field struct {
		name string
		typ  any
		tag  string
	}
	instance struct {
		definition reflect.Type
	}
)

func NewInstance() InstanceBuilder {
	return &builder{
		fields: []*field{},
	}
}

func (b *builder) AddField(name string, typ any, tag string) InstanceBuilder {
	b.fields = append(b.fields, &field{
		name: name,
		typ:  typ,
		tag:  tag,
	})
	return b
}

func (b *builder) RemoveField(name string) InstanceBuilder {
	for i := range b.fields {
		if b.fields[i].name == name {
			b.fields = append(b.fields[:i], b.fields[i+1:]...)
			break
		}
	}
	return b
}

func (b *builder) HasField(name string) bool {
	return b.GetField(name) != nil
}

func (b *builder) GetField(name string) InstanceField {
	for i := range b.fields {
		if b.fields[i].name == name {
			return b.fields[i]
		}
	}
	return nil
}

func (b *builder) Build() Instance {
	structFields := make([]reflect.StructField, len(b.fields))
	for i, f := range b.fields {
		structFields[i] = reflect.StructField{
			Name: f.name,
			Type: reflect.TypeOf(f.typ),
			Tag:  reflect.StructTag(f.tag),
		}
	}
	return &instance{
		definition: reflect.StructOf(structFields),
	}
}

func (f *field) SetType(typ any) InstanceField {
	f.typ = typ
	return f
}

func (f *field) SetTag(tag string) InstanceField {
	f.tag = tag
	return f
}

func (i *instance) New() any {
	return reflect.New(i.definition).Interface()
}
