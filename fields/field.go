package fields

import (
	"github.com/hashicorp/go-multierror"
	"github.com/malijoe/djanGo-unchained/utils"
)

type Internalizable interface {
	ToInternalValue(interface{}) error
	Internal() interface{}
}

type Representable interface {
	ToRepresentation() interface{}
}

type Field interface {
	Internalizable
	Representable
	MetaData() *Meta
}

func Validate(field Field) error {
	var (
		meta = field.MetaData()
	)

	var errs error
	for _, validator := range meta.Validators {
		if err := validator(field); err != nil {
			errs = multierror.Append(errs, err)
		}
	}
	if errs != nil {
		return errs
	}

	return nil
}

func FieldIsNullOrZero(f Field) bool {
	var (
		isZero  bool
		current = f.Internal()
	)

	switch z, isZeroer := current.(Zeroer); {
	case current == nil:
		isZero = true
	case isZeroer:
		isZero = z.IsZero()
	default:
		isZero = utils.IsZero(current)
	}
	return isZero
}

func TrySetDefault(f Field) bool {
	ok := FieldIsNullOrZero(f)
	if ok {
		var meta = f.MetaData()
		dv, hasDefault := meta.HasDefault()
		if hasDefault {
			_ = f.ToInternalValue(dv)
		}

		return hasDefault
	}
	return ok
}
