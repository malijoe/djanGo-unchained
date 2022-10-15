package db

import (
	"fmt"
	"strings"

	"github.com/uptrace/bun"
)

type SaveOption func(*bun.InsertQuery)

func WithConflictHandler(conflict string, fields []string) SaveOption {
	on := fmt.Sprintf("CONFLICT (%s) DO UPDATE", conflict)
	set := make([]string, len(fields))
	for i, field := range fields {
		set[i] = fmt.Sprintf("%s = EXCLUDED.%s", field, field)
	}
	return func(iq *bun.InsertQuery) {
		iq.On(on).Set(strings.Join(set, ","))
	}
}

type SelectOption func(*bun.SelectQuery)

func WithSelectCondition(specification Specification) SelectOption {
	return func(sq *bun.SelectQuery) {
		sq.Where(specification.Query(), specification.Values()...)
	}
}

type DeleteOption func(*bun.DeleteQuery)

func WithDeleteCondition(specification Specification) DeleteOption {
	return func(dq *bun.DeleteQuery) {
		dq.Where(specification.Query(), specification.Values()...)
	}
}

type UpdateOption func(*bun.UpdateQuery)

func WithUpdateCondition(specification Specification) UpdateOption {
	return func(uq *bun.UpdateQuery) {
		uq.Where(specification.Query(), specification.Values()...)
	}
}
