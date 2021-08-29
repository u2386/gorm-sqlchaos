package callback

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestApplyByReflectValue(t *testing.T) {
	i := int(1)
	i8 := int8(1)
	i16 := int16(1)
	i32 := int32(1)
	i64 := int64(1)
	ui := uint(1)
	ui8 := uint8(1)
	ui16 := uint16(1)
	ui32 := uint32(1)
	ui64 := uint64(1)
	f32 := float32(1.0)
	f64 := float64(1.0)
	s := "sql"

	vm := map[string]reflect.Value{
		"i":    reflect.ValueOf(&i),
		"i8":   reflect.ValueOf(&i8),
		"i16":  reflect.ValueOf(&i16),
		"i32":  reflect.ValueOf(&i32),
		"i64":  reflect.ValueOf(&i64),
		"ui":   reflect.ValueOf(&ui),
		"ui8":  reflect.ValueOf(&ui8),
		"ui16": reflect.ValueOf(&ui16),
		"ui32": reflect.ValueOf(&ui32),
		"ui64": reflect.ValueOf(&ui64),
		"f32":  reflect.ValueOf(&f32),
		"f64":  reflect.ValueOf(&f64),
		"s":    reflect.ValueOf(&s),
	}

	applier := Applier(map[string]string{
		"i":    "1024",
		"i8":   "1024",
		"i16":  "1024",
		"i32":  "1024",
		"i64":  "1024",
		"ui":   "1024",
		"ui8":  "1024",
		"ui16": "1024",
		"ui32": "1024",
		"ui64": "1024",
		"f32":  "10.24",
		"f64":  "10.24",
		"s":    "sqlchaos",
	})

	applied := applier.Apply(ApplyByReflectValue(vm))
	assert.True(t, applied)

	assert.Equal(t, int(1024), i)
	assert.Equal(t, int8(0), i8) // overflow
	assert.Equal(t, int16(1024), i16)
	assert.Equal(t, int32(1024), i32)
	assert.Equal(t, int64(1024), i64)
	assert.Equal(t, uint(1024), ui)
	assert.Equal(t, uint8(0), ui8) // overflow
	assert.Equal(t, uint16(1024), ui16)
	assert.Equal(t, uint32(1024), ui32)
	assert.Equal(t, uint64(1024), ui64)
	assert.Equal(t, float32(10.24), f32)
	assert.Equal(t, float64(10.24), f64)
	assert.Equal(t, "sqlchaos", s)
}
