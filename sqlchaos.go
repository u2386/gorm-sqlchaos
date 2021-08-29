package sqlchaos

import (
	"errors"
	"fmt"
	"os"

	"github.com/u2386/gorm-sqlchaos/pkg/callback"
	"gorm.io/gorm"
)

var _ gorm.Option = (*Config)(nil)
var (
	ErrDBNameNotSpecified = errors.New("DBName not specified")
)

type (
	Config struct {
		DBName string
	}
)

func (*Config) Apply(*gorm.Config) error {
	return nil
}

func (c *Config) AfterInitialize(db *gorm.DB) error {
	if c.DBName == "" {
		fmt.Fprintln(os.Stderr, "SQLChaos:db not specified")
		return ErrDBNameNotSpecified
	}
	fmt.Fprintln(os.Stderr, "SQLChaos:SQLChaos enabled")

	callback := &callback.Callback{}
	if err := db.Callback().Create().Before("before-create").Register("sqlchaos:before-create", callback.BeforeCreate()); err != nil {
		return err
	}
	if err := db.Callback().Update().Before("before-update").Register("sqlchaos:before-update", callback.BeforeUpdate()); err != nil {
		return err
	}
	return nil
}
