package django

import "context"

type test struct{}

type BaseSerializer[M any, Meta Metadata[M]] struct {
	meta           Meta
	DeclaredFields map[string]Field
	Instance       any
	Partial        bool
	Context        context.Context
	ValidatedData  map[string]any
	InitialData    map[string]any
	errs           []error
}

func (s *BaseSerializer[M, FM]) Update(m M) (*M, error) {
	return s.meta.DB().Update(s.Context, m)
}

func (s *BaseSerializer[M, FM]) Create(m M) (*M, error) {
	return s.meta.DB().Save(s.Context, m)
}

func (s *BaseSerializer[M, FM]) Save() (any, error) {
	panic("not implemented")
}

func (s *BaseSerializer[M, FM]) Get(id uint) (*M, error) {
	return s.meta.DB().Get(s.Context, id)
}

func (s *BaseSerializer[M, FM]) Delete(id uint) error {
	return s.meta.DB().Delete(s.Context, id)
}

func (s *BaseSerializer[M, FM]) IsValid() bool {
	return false
}

func (s *BaseSerializer[M, FM]) Errors() []error {
	return s.errs
}
