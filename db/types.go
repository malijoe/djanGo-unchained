package db

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
)

type (
	// helpful type aliases
	NullString = sql.NullString
	NullBool   = sql.NullBool
	NullTime   = sql.NullTime
)

type Null[O any] struct {
	Valid  bool
	Object O
}

func (n Null[O]) Value() (driver.Value, error) {
	var o O
	_, isValuer := any(o).(driver.Valuer)
	if n.Valid {
		if isValuer {
			return any(n.Object).(driver.Valuer).Value()
		}
		return n.Object, nil
	}
	return nil, nil
}

func (n *Null[O]) Scan(value any) error {
	var o O
	_, isScanner := any(o).(sql.Scanner)
	if isScanner {
		return any(n.Object).(sql.Scanner).Scan(value)
	}
	v, ok := value.(O)
	if ok {
		n.Object = v
		return nil
	}
	return fmt.Errorf("Null[%T] received a value of type %T", o, value)
}

func ToNull[O any](o *O) Null[O] {
	n := Null[O]{
		Valid: o != nil,
	}
	if n.Valid {
		n.Object = *o
	}
	return n
}

func FromNull[O any](n Null[O]) *O {
	if n.Valid {
		return &n.Object
	}
	return nil
}
