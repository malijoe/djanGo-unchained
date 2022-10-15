package db

import (
	"fmt"
	"testing"
	"time"
)

func TestNull(t *testing.T) {
	type stringTest struct {
		s     *string
		value *string
	}

	testLen := 10
	start := time.Now()
	tests := make([]stringTest, testLen)
	for i := range tests {
		var (
			s, value *string
		)
		if time.Now().Sub(start)%2 == 0 {
			v := fmt.Sprintf("s%d", i)
			s, value = &v, &v
		}
		tests[i] = stringTest{
			s:     s,
			value: value,
		}
	}

	for _, test := range tests {
		var ns = &Null[string]{}
		if test.s != nil {
			if err := ns.Scan(test.s); err != nil {
				*ns = ToNull(test.s)
				t.Errorf("Null[string].Scan(%s) got error %v", *test.s, err)
			}
		} else {
			*ns = ToNull(test.s)
		}
		v, _ := ns.Value()
		if (test.value == nil && v != nil) || (test.value != nil && *test.value != v) {
			t.Errorf("Null[string].Value() got %v, want %v", v, test.value)
		}
		if test.value == nil && ns.Valid {
			t.Error("Null[string].Valid got true, want false")
		}
		if test.value != nil && *test.value != ns.Object {
			t.Errorf("Null[string].Object got %s, want %s", ns.Object, *test.value)
		}
	}
}
