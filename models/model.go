package models

import (
	"context"

	"github.com/malijoe/djanGo-unchained/db"
)

type DataModel[T any] interface {
	Instance() *T
	Save(ctx context.Context) (*T, error)
	Create(ctx context.Context) (*T, error)
	Update(ctx context.Context) (*T, error)
	Find(ctx context.Context, query db.Specification) ([]T, error)
	Get(ctx context.Context, id uint) (*T, error)
	Delete(ctx context.Context, id uint) error
}

type Model[T any] interface {
	Objects() DataModel[T]
}

type dataModel[T any, R db.BaseRepository[T]] struct {
	instance *T
	repo     R
}

func NewDataModel[T any, R db.BaseRepository[T]](repo R, instance *T) DataModel[T] {
	return &dataModel[T, R]{
		repo:     repo,
		instance: instance,
	}
}

func (m *dataModel[T, R]) Instance() *T {
	return m.instance
}

func (m *dataModel[T, R]) Save(ctx context.Context) (*T, error) {
	panic("save method not implemented")
}

func (m *dataModel[T, R]) Create(ctx context.Context) (*T, error) {
	return m.repo.Save(ctx, *m.instance)
}

func (m *dataModel[T, R]) Update(ctx context.Context) (*T, error) {
	return m.repo.Update(ctx, *m.instance)
}

func (m *dataModel[T, R]) Get(ctx context.Context, id uint) (*T, error) {
	return m.repo.Get(ctx, id)
}

func (m *dataModel[T, R]) Find(ctx context.Context, query db.Specification) ([]T, error) {
	return m.repo.Find(ctx, query)
}

func (m *dataModel[T, R]) Delete(ctx context.Context, id uint) error {
	return m.repo.Delete(ctx, id)
}
