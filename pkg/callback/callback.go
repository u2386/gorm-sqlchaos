package callback

import (
	"context"
	"fmt"
	"os"
	"reflect"
	"strings"

	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type (
	Rules interface {
		Get(ctx context.Context, key string) (*ChaosRule, error)
	}

	ChaosRule struct {
		DML  string
		When string
		Then string
	}

	Callback struct {
		Rules Rules
	}
)

var (
	DEBUG          = os.Getenv("SQLCHAOS_DEBUG") != ""
	SQLCHAOS_DEBUG = func(format string, args ...interface{}) {
		if DEBUG {
			fmt.Fprintf(os.Stdout, format+"\n", args...)
		}
	}
	SQLCHAOS_ERROR = func(format string, args ...interface{}) { fmt.Fprintf(os.Stderr, format, args...) }
)

func (c *Callback) BeforeCreate() func(*gorm.DB) {
	return func(db *gorm.DB) {
		stmt := db.Statement
		if stmt.Schema == nil {
			fmt.Fprintln(os.Stderr, "statement scheme not provided")
			return
		}

		key := db.Statement.Table
		rule, err := c.GetTableRule(context.Background(), key, "CREATE")
		if err != nil {
			SQLCHAOS_ERROR("read rule failed:%+v", err)
			return
		}
		if rule == nil {
			return
		}
		SQLCHAOS_DEBUG("SQLCHAOS:get chaos rule:%#v", rule)

		if ApplyRule(rule, stmt) {
			SQLCHAOS_ERROR("SQLCHAOS:records have been modified")
		}
	}
}

func (c *Callback) BeforeUpdate() func(*gorm.DB) {
	return func(db *gorm.DB) {
		stmt := db.Statement
		if stmt.Schema == nil {
			fmt.Fprintln(os.Stderr, "statement scheme not provided")
			return
		}

		key := db.Statement.Table
		rule, err := c.GetTableRule(context.Background(), key, "UPDATE")
		if err != nil {
			SQLCHAOS_ERROR("read rule failed:%+v", err)
			return
		}
		if rule == nil {
			return
		}
		SQLCHAOS_DEBUG("SQLCHAOS:get chaos rule:%#v", rule)

		if ApplyRule(rule, stmt) {
			SQLCHAOS_ERROR("SQLCHAOS:records have been modified")
		}
	}
}

func ApplyRule(rule *ChaosRule, stmt *gorm.Statement) (applied bool) {
	matcher, err := ParseWhereStatement(rule.When)
	if err != nil {
		SQLCHAOS_ERROR("where statement invalid:%+v", err)
		return
	}

	applier, err := ParseSetStatement(rule.Then)
	if err != nil {
		SQLCHAOS_ERROR("set statement invalid:%+v", err)
		return
	}
	return ApplyValuesIfMatch(stmt, matcher, applier)
}

func (c *Callback) GetTableRule(ctx context.Context, key, dml string) (rule *ChaosRule, err error) {
	rule, err = c.Rules.Get(context.TODO(), key)
	if err != nil || rule == nil {
		return
	}

	if !strings.EqualFold(rule.DML, dml) {
		return
	}
	return
}

// Canonical converts struct field name to database column name
func Canonical(schema *schema.Schema, v map[string]interface{}) {
	for name, value := range v {
		if dbname, ok := schema.FieldsByName[name]; ok {
			v[dbname.DBName] = value
			delete(v, name)
			continue
		}
		v[name] = value
	}
}

func ApplyValuesIfMatch(stmt *gorm.Statement, matcher Matcher, applier Applier) (applied bool) {
	if v, ok := stmt.Dest.(map[string]interface{}); ok {
		Canonical(stmt.Schema, v)
		if matcher.Match(MatchByInterface(v)) {
			applied = applier.Apply(ApplyByInterface(v))
		}

	} else if vs, ok := stmt.Dest.([]map[string]interface{}); ok {
		for _, v := range vs {
			Canonical(stmt.Schema, v)
			if matcher.Match(MatchByInterface(v)) {
				applied = applied || applier.Apply(ApplyByInterface(v))
			}
		}
	} else {
		switch stmt.ReflectValue.Kind() {
		case reflect.Slice, reflect.Array:
			for i := 0; i < stmt.ReflectValue.Len(); i++ {
				one := stmt.ReflectValue.Index(i)
				vs := make(map[string]reflect.Value)
				for name, field := range stmt.Schema.FieldsByName {
					v := one.FieldByName(name)
					if !v.IsValid() || (v.Kind() == reflect.Ptr && v.IsNil()) {
						continue
					}
					vs[field.DBName] = v
				}
				if matcher.Match(MatchByReflectValue(vs)) {
					applied = applied || applier.Apply(ApplyByReflectValue(vs))
				}
			}
		case reflect.Struct:
			vs := make(map[string]reflect.Value)
			for name, field := range stmt.Schema.FieldsByName {
				v := stmt.ReflectValue.FieldByName(name)
				if !v.IsValid() || (v.Kind() == reflect.Ptr && v.IsNil()) {
					continue
				}
				vs[field.DBName] = v
			}
			if matcher.Match(MatchByReflectValue(vs)) {
				applied = applier.Apply(ApplyByReflectValue(vs))
			}
		}
	}
	return
}
