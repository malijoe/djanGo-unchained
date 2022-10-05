package db

import (
	"context"
	"database/sql"
	"os"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/extra/bundebug"
	"github.com/uptrace/bun/schema"
)

type InsertModifier interface {
	OnInsert(query *bun.InsertQuery) *bun.InsertQuery
}

type DTO[E any] interface {
	FromEntity(e E)
	ToEntity() E
}

type BaseRepository[T any] interface {
	Save(ctx context.Context, t T) (*T, error)
	Get(ctx context.Context, id uint) (*T, error)
	Find(ctx context.Context, query Specification) ([]T, error)
	Update(ctx context.Context, t T) (*T, error)
	Delete(ctx context.Context, id uint) error
}

type Repository[T any] interface {
	WithTx(tx ...*sql.Tx) (TxRepository[T], error)
	BaseRepository[T]
}

type TxRepository[T any] interface {
	Tx() *sql.Tx
	Commit() error
	Rollback() error
	BaseRepository[T]
}

type conn interface {
	NewInsert() *bun.InsertQuery
	NewSelect() *bun.SelectQuery
	NewUpdate() *bun.UpdateQuery
	NewDelete() *bun.DeleteQuery
}

type baseRepository[D DTO[E], E any] struct {
	db conn
}

func newBaseRepository[D DTO[E], E any](conn conn) BaseRepository[E] {
	return &baseRepository[D, E]{
		db: conn,
	}
}

func (r *baseRepository[D, E]) Save(ctx context.Context, e E) (*E, error) {
	var dto D
	dto.FromEntity(e)

	stmt := r.db.NewInsert().Model(&dto).Returning("id")
	if inserter, ok := any(dto).(InsertModifier); ok {
		stmt = inserter.OnInsert(stmt)
	}

	if _, err := stmt.Exec(ctx); err != nil {
		return nil, err
	}
	entity := dto.ToEntity()
	return &entity, nil
}

func (r *baseRepository[D, E]) Get(ctx context.Context, id uint) (*E, error) {
	var dto D
	stmt := r.db.NewSelect().Model(&dto).Where("id = ?", id)
	if err := stmt.Scan(ctx); err != nil {
		return nil, err
	}
	entity := dto.ToEntity()
	return &entity, nil
}

func (r *baseRepository[D, E]) Find(ctx context.Context, query Specification) ([]E, error) {
	var dto []D
	stmt := r.db.NewSelect().Model(&dto)
	if query != nil {
		stmt = stmt.Where(query.Query(), query.Values()...)
	}

	if err := stmt.Scan(ctx); err != nil {
		return nil, err
	}

	response := make([]E, len(dto))
	for i := range dto {
		response[i] = dto[i].ToEntity()
	}
	return response, nil
}

func (r *baseRepository[D, E]) Update(ctx context.Context, e E) (*E, error) {
	var dto D
	dto.FromEntity(e)

	stmt := r.db.NewUpdate().Model(&dto).OmitZero().WherePK()

	_, err := stmt.Exec(ctx)
	if err != nil {
		return nil, err
	}

	entity := dto.ToEntity()
	return &entity, nil
}

func (r *baseRepository[D, E]) Delete(ctx context.Context, id uint) error {
	var dto D
	stmt := r.db.NewDelete().Model(&dto).Where("id = ?", id)
	_, err := stmt.Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}

type repository[D DTO[E], E any] struct {
	db *bun.DB
	BaseRepository[E]
}

func NewRepository[D DTO[E], E any](conn *sql.DB, dialect schema.Dialect) Repository[E] {
	db := bun.NewDB(conn, dialect)
	db.AddQueryHook(
		bundebug.NewQueryHook(
			bundebug.WithEnabled(true),
			bundebug.WithVerbose(true),
			bundebug.WithWriter(os.Stdout),
		))
	r := repository[D, E]{
		db:             db,
		BaseRepository: newBaseRepository[D, E](db),
	}
	return &r
}

func (r *repository[D, E]) WithTx(tx ...*sql.Tx) (TxRepository[E], error) {
	return newTxRepository[D, E](r.db, tx...)
}

type txRepository[D DTO[E], E any] struct {
	tx bun.Tx
	BaseRepository[E]
}

func newTxRepository[D DTO[E], E any](db *bun.DB, tx ...*sql.Tx) (TxRepository[E], error) {
	ttx, err := db.Begin()
	if err != nil {
		return nil, err
	}
	if len(tx) > 0 {
		// close the tx started by Begin()
		_ = ttx.Tx.Rollback()
		ttx.Tx = tx[0]
	}
	return &txRepository[D, E]{
		tx:             ttx,
		BaseRepository: newBaseRepository[D, E](ttx),
	}, nil
}

func (r *txRepository[D, E]) Commit() error {
	return r.tx.Commit()
}

func (r *txRepository[D, E]) Rollback() error {
	return r.tx.Rollback()
}

func (r *txRepository[D, E]) Tx() *sql.Tx {
	return r.tx.Tx
}
