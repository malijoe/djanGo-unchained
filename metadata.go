package django

type Metadata[M any] interface {
	Model() M
	GetFields() []Field
	DB() Repository[M]
}
