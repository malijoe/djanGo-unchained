package fields

type Meta struct {
	ReadOnly    bool             `json:"read_only"`
	WriteOnly   bool             `json:"write_only"`
	Required    bool             `json:"required"`
	Default     interface{}      `json:"default"`
	AllowNull   bool             `json:"allow_null"`
	AllowBlank  bool             `json:"allow_blank"`
	Source      string           `json:"source"`
	Label       string           `json:"label"`
	HelpText    string           `json:"help_text"`
	Type        FieldType        `json:"type"`
	PlaceHolder string           `json:"place_holder"`
	Validators  []FieldValidator `json:"-"`
}

func (m Meta) HasDefault() (interface{}, bool) {
	return m.Default, m.Default != nil
}

func (m Meta) IsReadOnly() bool {
	return m.ReadOnly
}

func (m Meta) IsWriteOnly() bool {
	return m.WriteOnly
}
