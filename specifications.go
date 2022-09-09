package django

import (
	"fmt"
	"strings"
)

type Specification interface {
	Query() string
	Values() []any
}

type compositeSpecification struct {
	specifications []Specification
	separator      string
}

func (s compositeSpecification) Query() string {
	queries := make([]string, len(s.specifications))
	for i, spec := range s.specifications {
		queries[i] = spec.Query()
	}
	return strings.Join(queries, fmt.Sprintf(" %s ", s.separator))
}

func (s compositeSpecification) Values() []any {
	var values []any
	for _, spec := range s.specifications {
		values = append(values, spec.Values()...)
	}
	return values
}

func And(specifications ...Specification) Specification {
	return compositeSpecification{
		specifications: specifications,
		separator:      "AND",
	}
}

func Or(specifications ...Specification) Specification {
	return compositeSpecification{
		specifications: specifications,
		separator:      "OR",
	}
}

type notSpecification struct {
	Specification
}

func (s notSpecification) Query() string {
	return fmt.Sprintf(" NOT (%s)", s.Specification.Query())
}

func Not(specification Specification) Specification {
	return notSpecification{
		specification,
	}
}

type binaryOperatorSpecification[T any] struct {
	field    string
	operator string
	value    T
}

func (s binaryOperatorSpecification[T]) Query() string {
	return fmt.Sprintf("%s %s ?", s.field, s.operator)
}

func (s binaryOperatorSpecification[T]) Values() []any {
	return []any{s.value}
}

func Equal[T any](field string, value T) Specification {
	return binaryOperatorSpecification[T]{
		field:    field,
		operator: "=",
		value:    value,
	}
}

type invertedBinaryOperatorSpecification[T any] struct {
	binaryOperatorSpecification[T]
}

func (s invertedBinaryOperatorSpecification[T]) Query() string {
	return fmt.Sprintf("? %s %s", s.operator, s.field)
}

func EqualInverted[T any](field string, value T) Specification {
	return invertedBinaryOperatorSpecification[T]{
		binaryOperatorSpecification[T]{
			field:    field,
			operator: "=",
			value:    value,
		},
	}
}

type inSpecification[T any] struct {
	field  string
	values []T
}

func (s inSpecification[T]) Query() string {
	qStr := make([]string, len(s.values))
	for i := range s.values {
		qStr[i] = "?"
	}
	return fmt.Sprintf("%s IN (%s)", s.field, strings.Join(qStr, ","))
}

func (s inSpecification[T]) Values() []any {
	values := make([]any, len(s.values))
	for i := range s.values {
		values[i] = s.values[i]
	}
	return values
}

func In[T any](field string, values []T) Specification {
	return inSpecification[T]{
		field:  field,
		values: values,
	}
}

func QueryString(query Specification) string {
	var queryString string
	if query != nil {
		queryString = fmt.Sprintf("%s %v", query.Query(), query.Values())
	}
	return queryString
}
