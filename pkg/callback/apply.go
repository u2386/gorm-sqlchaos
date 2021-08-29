package callback

import (
	"reflect"
	"strconv"
)

type (
	ApplyBy interface {
		Apply(values map[string]string) bool
	}

	ApplyByInterface    map[string]interface{}
	ApplyByReflectValue map[string]reflect.Value
	Applier             map[string]string
)

func (applier Applier) Apply(by ApplyBy) bool {
	return by.Apply(applier)
}

func (m ApplyByInterface) Apply(values map[string]string) bool {
	for name, value := range values {
		m[name] = value
	}
	return true
}

func (vm ApplyByReflectValue) Apply(values map[string]string) (applied bool) {
	for name, value := range values {
		v, ok := vm[name]
		if !ok {
			SQLCHAOS_ERROR("can not set value %s: not found", name)
			return
		}

		v = reflect.Indirect(reflect.Value(v))
		switch v.Kind() {
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			x, err := strconv.ParseUint(value, 10, 64)
			if err != nil {
				SQLCHAOS_ERROR("can not set value %s: %+v", name, err)
				return
			}
			v.SetUint(x)

		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			x, err := strconv.ParseInt(value, 10, 64)
			if err != nil {
				SQLCHAOS_ERROR("can not set value %s: %+v", name, err)
				return
			}
			v.SetInt(x)

		case reflect.Float32, reflect.Float64:
			x, err := strconv.ParseFloat(value, 64)
			if err != nil {
				SQLCHAOS_ERROR("can not set value %s: %+v", name, err)
				return
			}
			v.SetFloat(x)

		case reflect.String:
			v.SetString(value)

		default:
			SQLCHAOS_ERROR("can not set value %s: unsupported type", name)
			return
		}
	}
	return true
}
