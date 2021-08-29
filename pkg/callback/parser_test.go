package callback

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseSetStatement(t *testing.T) {
	applier, err := ParseThenStatement("a=10.24,b=bar")
	assert.Nil(t, err)

	assert.Equal(t, "10.24", applier["a"])
	assert.Equal(t, "bar", applier["b"])

	_, err = ParseThenStatement("a==10.24,b=bar")
	assert.Error(t, err)
}

func TestParseWhenStatement(t *testing.T) {
	matches, err := ParseWhenStatement("a=foo")
	assert.Nil(t, err)
	assert.True(t, matches["a"]("foo"))

	matches, err = ParseWhenStatement("a=foo AND b=bar")
	assert.Nil(t, err)
	assert.True(t, matches["a"]("foo"))
	assert.True(t, matches["b"]("bar"))

	matches, err = ParseWhenStatement("a>=0 AND a<2")
	assert.Nil(t, err)
	assert.True(t, matches["a"](0))
	assert.True(t, matches["a"](1))
	assert.False(t, matches["a"](2))
	assert.False(t, matches["a"](-1))
}