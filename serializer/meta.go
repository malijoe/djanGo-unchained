package serializer

type ModFieldTime int32

const (
	Invalid ModFieldTime = iota
	PreRead
	PreWrite
	PostRead
	PostWrite
)

type Meta struct {
	Fields          []string
	Modifiers       map[ModFieldTime][]FieldModifierSpec
	ReadOnlyFields  []string
	WriteOnlyFields []string
}

type FieldModifierSpec struct {
	Fields   []string
	Modifier FieldModifier
}

func initMeta(meta *Meta) {
	if meta.Modifiers == nil {
		meta.Modifiers = make(map[ModFieldTime][]FieldModifierSpec)
	}
	if len(meta.ReadOnlyFields) > 0 {
		meta.Modifiers[PreWrite] = append([]FieldModifierSpec{setModifierForFields(SetFieldReadOnly, meta.ReadOnlyFields)}, meta.Modifiers[PreWrite]...)
	}

	if len(meta.WriteOnlyFields) > 0 {
		meta.Modifiers[PreRead] = append([]FieldModifierSpec{setModifierForFields(SetFieldWriteOnly, meta.WriteOnlyFields)}, meta.Modifiers[PreRead]...)
	}
}

func setModifierForFields(modifier FieldModifier, fields []string) FieldModifierSpec {
	return FieldModifierSpec{
		Modifier: modifier,
		Fields:   fields,
	}
}
