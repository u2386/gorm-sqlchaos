package sqlchaos

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/u2386/gorm-sqlchaos/pkg/callback"
	"gorm.io/gorm"
)

var _ gorm.Option = (*Config)(nil)

type (
	ReadRule func(ctx context.Context, dbname, table string) ([]byte, error)
	Config   struct {
		DBName     string
		RuleReader ReadRule
	}
)

var (
	ErrArgumentNotSpecified = errors.New("argument not specified")
)

func (f ReadRule) Get(ctx context.Context, dbname, table string) (*callback.ChaosRule, error) {
	value, err := f(ctx, dbname, table)
	if err != nil {
		return nil, err
	}
	r := &callback.ChaosRule{}
	if err := json.Unmarshal(value, r); err != nil {
		return nil, err
	}
	return r, nil
}

func (*Config) Apply(*gorm.Config) error {
	return nil
}

func (c *Config) AfterInitialize(db *gorm.DB) error {
	if c.DBName == "" {
		fmt.Fprintln(os.Stderr, "SQLChaos:db not specified")
		return ErrArgumentNotSpecified
	}
	if c.RuleReader == nil {
		fmt.Fprintln(os.Stderr, "SQLChaos:rule reader not specified")
		return ErrArgumentNotSpecified
	}
	fmt.Fprintln(os.Stderr, "SQLChaos:SQLChaos enabled")

	callback := &callback.Callback{
		DBName: c.DBName,
		Rules: c.RuleReader,
	}
	if err := db.Callback().Create().Before("gorm:create").Register("sqlchaos:before-create", callback.BeforeCreate()); err != nil {
		return err
	}
	if err := db.Callback().Update().Before("gorm:update").Register("sqlchaos:before-update", callback.BeforeUpdate()); err != nil {
		return err
	}
	return nil
}
