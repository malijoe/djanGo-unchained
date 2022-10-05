package db

import (
	"context"
	"database/sql"
)

type Opt[T any] func(*MockRepository[T])

func WithSaveFn[T any](saveFn func(*T) (*T, error)) Opt[T] {
	return func(r *MockRepository[T]) {
		r.saveFn = saveFn
	}
}

func WithGetFn[T any](getFn func(uint) (*T, error)) Opt[T] {
	return func(r *MockRepository[T]) {
		r.getFn = getFn
	}
}

func WithFindFn[T any](findFn func(Specification) ([]T, error)) Opt[T] {
	return func(r *MockRepository[T]) {
		r.findFn = findFn
	}
}

func WithUpdateFn[T any](updateFn func(*T) (*T, error)) Opt[T] {
	return func(r *MockRepository[T]) {
		r.updateFn = updateFn
	}
}

func WithDeleteFn[T any](deleteFn func(uint) error) Opt[T] {
	return func(r *MockRepository[T]) {
		r.deleteFn = deleteFn
	}
}

type MockRepository[T any] struct {
	saveFn      func(t *T) (*T, error)
	SaveInvoked bool

	getFn      func(id uint) (*T, error)
	GetInvoked bool

	findFn      func(query Specification) ([]T, error)
	FindInvoked bool

	updateFn      func(t *T) (*T, error)
	UpdateInvoked bool

	deleteFn      func(id uint) error
	DeleteInvoked bool

	WithTxInvoked   bool
	CommitInvoked   bool
	RollbackInvoked bool
}

func NewMockRepository[T any](opts ...Opt[T]) *MockRepository[T] {
	repo := new(MockRepository[T])
	for _, opt := range opts {
		opt(repo)
	}
	return repo
}

func (r *MockRepository[T]) WithTx(_ ...*sql.Tx) (TxRepository[T], error) {
	r.WithTxInvoked = true
	return r, nil
}

func (r *MockRepository[T]) Tx() *sql.Tx {
	return nil
}

func (r *MockRepository[T]) Commit() error {
	r.CommitInvoked = true
	return nil
}

func (r *MockRepository[T]) Rollback() error {
	r.RollbackInvoked = true
	return nil
}

func (r *MockRepository[T]) Save(_ context.Context, t T) (*T, error) {
	r.SaveInvoked = true
	if r.saveFn != nil {
		return r.saveFn(&t)
	}
	return &t, nil
}

func (r *MockRepository[T]) Get(_ context.Context, id uint) (*T, error) {
	r.GetInvoked = true
	if r.getFn != nil {
		return r.getFn(id)
	}
	return nil, nil
}

func (r *MockRepository[T]) Find(_ context.Context, query Specification) ([]T, error) {
	r.FindInvoked = true
	if r.findFn != nil {
		return r.findFn(query)
	}
	return nil, nil
}

func (r *MockRepository[T]) Update(_ context.Context, t T) (*T, error) {
	r.UpdateInvoked = true
	if r.updateFn != nil {
		return r.updateFn(&t)
	}
	return &t, nil
}

func (r *MockRepository[T]) Delete(_ context.Context, id uint) error {
	r.DeleteInvoked = true
	if r.deleteFn != nil {
		return r.deleteFn(id)
	}
	return nil
}
