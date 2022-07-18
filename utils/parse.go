package utils

import "strconv"

func ParseInt(arg interface{}) (num int, ok bool) {
	switch x := arg.(type) {
	case string:
		temp, err := strconv.Atoi(x)
		if ok = err == nil; ok {
			num = temp
		}
	case int:
		ok = true
		num = x
	case int32:
		ok = true
		num = int(x)
	case int64:
		ok = true
		num = int(x)
	case float64:
		ok = true
		num = int(x)
	}
	return
}
