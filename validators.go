package django

import (
	"errors"
	"strings"
)

func MinLengthValidator(min int) func(any) error {
	return func(a any) error {
		switch v := a.(type) {
		case string:
			if len(v) > min {
				return nil
			}
		default:
			return errors.New("unsupported data type for min_length validator")
		}
		return errors.New("does not meet the min_length requirement")
	}
}

func MaxLengthValidator(max int) func(any) error {
	return func(a any) error {
		switch v := a.(type) {
		case string:
			if len(v) < max {
				return nil
			}
		default:
			return errors.New("unsupported data type for min_length validator")
		}
		return errors.New("does not meet the max_length requirement")
	}
}

func ProhibitNullCharactersValidator(value any) error {
	v, ok := value.(string)
	if !ok {
		return errors.New("unsupported datatype for prohibit_null_characters validator")
	}
	if strings.Contains(v, "\x00") {
		return errors.New("null characters are not allowed")
	}
	return nil
}

func MaxValueValidator(max int) func(any) error {
	return func(a any) error {
		v, ok := a.(int)
		if !ok {
			return errors.New("unsupported datatype for max_value validator")
		}
		if v > max {
			return errors.New("does not meet the max_value requirement")
		}
		return nil
	}
}

func MinValueValidator(min int) func(any) error {
	return func(a any) error {
		v, ok := a.(int)
		if !ok {
			return errors.New("unsupported datatype for min_value validator")
		}
		if v < min {
			return errors.New("does not meet the min_value requirement")
		}
		return nil
	}
}
