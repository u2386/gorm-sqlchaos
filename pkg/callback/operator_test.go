package callback

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOperator(t *testing.T) {
	assert.True(t, eq("a")("a"))
	assert.False(t, eq("a")("b"))
	assert.True(t, eq("10.24")(10.24))
	assert.False(t, eq("10.24")(10.23))

	assert.True(t, gt("a")("b"))
	assert.False(t, gt("a")("a"))
	assert.True(t, gt("10.24")(10.25))
	assert.False(t, gt("10.24")(10.24))

	assert.True(t, gte("b")("b"))
	assert.True(t, gte("b")("c"))
	assert.False(t, gte("b")("a"))
	assert.True(t, gte("10.24")(10.25))
	assert.True(t, gte("10.24")(10.24))
	assert.False(t, gte("10.24")(10.23))

	assert.False(t, lt("b")("b"))
	assert.True(t, lt("b")("a"))
	assert.False(t, lt("10.24")(10.24))
	assert.True(t, lt("10.24")(10.23))

	assert.True(t, lte("b")("b"))
	assert.False(t, lte("b")("c"))
	assert.True(t, lte("b")("a"))
	assert.False(t, lte("10.24")(10.25))
	assert.True(t, lte("10.24")(10.24))
	assert.True(t, lte("10.24")(10.23))

	assert.True(t, and(gt("0"), lte("2"))(2))
}
