package django

import "context"

type Repository[M any] interface {
	Save(ctx context.Context, m M) (*M, error)
	Get(ctx context.Context, id uint) (*M, error)
	Find(ctx context.Context, query Specification) ([]M, error)
	Update(ctx context.Context, m M) (*M, error)
	Delete(ctx context.Context, id uint) error
}
