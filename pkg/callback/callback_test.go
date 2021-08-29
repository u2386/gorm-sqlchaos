package callback

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"bou.ke/monkey"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type ff func(context.Context, string, string) string

func (f ff) Get(ctx context.Context, dbname, table string) (*ChaosRule, error) {
	f(ctx, dbname, table)
	return nil, nil
}

func TestGetTableRule(t *testing.T) {
	var (
		c     *Callback
		table string = "bar"
		dml   string
		g     *monkey.PatchGuard
		db    = &gorm.DB{
			Statement: &gorm.Statement{
				Table:  table,
				Schema: &schema.Schema{},
			},
		}
	)
	g = monkey.PatchInstanceMethod(reflect.TypeOf(c), "GetTableRule", func(_ *Callback, ctx context.Context, a, b string) (rule *ChaosRule, err error) {
		assert.Equal(t, table, a)
		assert.Equal(t, dml, b)

		g.Unpatch()
		defer g.Restore()
		return c.GetTableRule(ctx, a, b)
	})

	c = &Callback{
		DBName: "foo",
		Rules: ff(func(c context.Context, s1, s2 string) string {
			assert.Equal(t, "foo", s1)
			assert.Equal(t, table, s2)
			return ""
		}),
	}

	dml = "CREATE"
	c.BeforeCreate()(db)

	dml = "UPDATE"
	c.BeforeUpdate()(db)
}

func TestApplyRule(t *testing.T) {
	r := &ChaosRule{}

	g := monkey.Patch(ParseWhenStatement, func(string) (Matcher, error) {
		return nil, errors.New("")
	})
	assert.False(t, ApplyRule(r, nil))
	g.Unpatch()

	g = monkey.Patch(ParseThenStatement, func(string) (Applier, error) {
		return nil, errors.New("")
	})
	assert.False(t, ApplyRule(r, nil))
	g.Unpatch()
}

func TestApplyValuesIfMatch(t *testing.T) {
	monkey.Patch(Canonical, func(*schema.Schema, map[string]interface{}) {})

	var (
		applied int
		matcher Matcher
		applier Applier
		stmt    *gorm.Statement
	)

	cases := []struct {
		Matched bool
		Dest    interface{}
		Applied int
	}{
		{false, make(map[string]interface{}), 0},
		{true, make(map[string]interface{}), 1},

		{false, make([]map[string]interface{}, 2), 0},
		{true, make([]map[string]interface{}, 2), 2},
	}

	monkey.PatchInstanceMethod(reflect.TypeOf(applier), "Apply", func(Applier, ApplyBy) bool {
		applied++
		return true
	})

	for _, cs := range cases {
		func() {
			defer func() { applied = 0 }()
			g := monkey.PatchInstanceMethod(reflect.TypeOf(matcher), "Match", func(Matcher, MatchBy) bool {
				return cs.Matched
			})
			defer g.Unpatch()

			stmt = &gorm.Statement{Dest: cs.Dest}

			assert.Equal(t, cs.Matched, ApplyValuesIfMatch(stmt, matcher, applier))
			assert.Equal(t, cs.Applied, applied)
		}()
	}
}
