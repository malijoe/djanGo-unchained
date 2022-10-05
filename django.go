package django

type Serializer interface {
	Metadata() []Field
}

type Field struct {
	ReadOnly  bool
	WriteOnly bool
	Required  bool
	AllowNull bool
	Default   any
	Label     string
	HelpText  string
	Name      string
}
