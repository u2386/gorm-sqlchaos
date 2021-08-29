package callback

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMatchByReflectValue(t *testing.T) {
	x := 1
	m := map[string]reflect.Value{
		"foo": reflect.ValueOf(&x),
	}
	assert.False(t, MatchByReflectValue(m).Match(map[string]Match{"bar":nil}))
	assert.False(t, MatchByReflectValue(m).Match(map[string]Match{"foo": never()}))
	assert.True(t, MatchByReflectValue(m).Match(map[string]Match{"foo": eq("1")}))
}

func TestMatchByInterface(t *testing.T) {
	x := 1
	m := map[string]interface{}{
		"foo": x,
	}
	assert.False(t, MatchByInterface(m).Match(map[string]Match{"bar":nil}))
	assert.False(t, MatchByInterface(m).Match(map[string]Match{"foo": never()}))
	assert.True(t, MatchByInterface(m).Match(map[string]Match{"foo": eq("1")}))
}