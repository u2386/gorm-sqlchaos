package callback

import (
	"fmt"
	"reflect"
	"strconv"
)

var (
	operators = []struct {
		Sign string
		Fn   func(string) Match
	}{{
		">=", gte,
	}, {
		"<=", lte,
	}, {
		">", gt,
	}, {
		"<", lt,
	}, {
		"=", eq,
	}}
)

func compareAsFloat64(a, b interface{}, cmp func(a, b float64) bool) bool {
	x, err := strconv.ParseFloat(fmt.Sprint(a), 64)
	if err != nil {
		return false
	}
	y, err := strconv.ParseFloat(fmt.Sprint(b), 64)
	if err != nil {
		return false
	}
	return cmp(x, y)
}

func and(ms ...Match) Match {
	return func(actual interface{}) bool {
		for _, match := range ms {
			if !match(actual) {
				return false
			}
		}
		return true
	}
}

func never() Match {
	return func(actual interface{}) bool { return false }
}

func eq(expect string) Match {
	return func(actual interface{}) bool {
		actual = reflect.Indirect(reflect.ValueOf(actual)).Interface()

		switch actual.(type) {
		case uint, uint8, uint16, uint32, uint64, int, int8, int16, int32, int64, float32, float64:
			return compareAsFloat64(actual, expect, func(actual, expect float64) bool {
				return expect == actual
			})
		case string:
			return fmt.Sprint(actual) == expect
		default:
			return false
		}
	}
}

func gt(expect string) Match {
	return func(actual interface{}) bool {
		actual = reflect.Indirect(reflect.ValueOf(actual)).Interface()

		switch actual.(type) {
		case uint, uint8, uint16, uint32, uint64, int, int8, int16, int32, int64, float32, float64:
			return compareAsFloat64(actual, expect, func(actual, expect float64) bool {
				return actual > expect
			})
		case string:
			return fmt.Sprint(actual) > expect
		default:
			return false
		}
	}
}

func gte(expect string) Match {
	return func(actual interface{}) bool {
		actual = reflect.Indirect(reflect.ValueOf(actual)).Interface()

		switch actual.(type) {
		case uint, uint8, uint16, uint32, uint64, int, int8, int16, int32, int64, float32, float64:
			return compareAsFloat64(actual, expect, func(actual, expect float64) bool {
				return actual >= expect
			})
		case string:
			return fmt.Sprint(actual) >= expect
		default:
			return false
		}
	}
}

func lt(expect string) Match {
	return func(actual interface{}) bool {
		actual = reflect.Indirect(reflect.ValueOf(actual)).Interface()

		switch actual.(type) {
		case uint, uint8, uint16, uint32, uint64, int, int8, int16, int32, int64, float32, float64:
			return compareAsFloat64(actual, expect, func(actual, expect float64) bool {
				return actual < expect
			})
		case string:
			return fmt.Sprint(actual) < expect
		default:
			return false
		}
	}
}

func lte(expect string) Match {
	return func(actual interface{}) bool {
		actual = reflect.Indirect(reflect.ValueOf(actual)).Interface()

		switch actual.(type) {
		case uint, uint8, uint16, uint32, uint64, int, int8, int16, int32, int64, float32, float64:
			return compareAsFloat64(actual, expect, func(actual, expect float64) bool {
				return actual <= expect
			})
		case string:
			return fmt.Sprint(actual) <= expect
		default:
			return false
		}
	}
}
