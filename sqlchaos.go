package sqlchaos

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

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
	if err != nil || value == nil {
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
		fmt.Fprintln(os.Stderr, "SQLCHAOS:db not specified")
		return ErrArgumentNotSpecified
	}
	if c.RuleReader == nil {
		fmt.Fprintln(os.Stderr, "SQLCHAOS:rule reader not specified")
		return ErrArgumentNotSpecified
	}
	fmt.Fprintln(os.Stderr, "SQLCHAOS:SQLChaos enabled")

	callback := &callback.Callback{
		DBName: c.DBName,
		Rules:  c.RuleReader,
	}
	if err := db.Callback().Create().Before("gorm:create").Register("sqlchaos:before-create", callback.BeforeCreate()); err != nil {
		return err
	}
	if err := db.Callback().Update().Before("gorm:update").Register("sqlchaos:before-update", callback.BeforeUpdate()); err != nil {
		return err
	}
	return nil
}

func WithSimpleHTTPRuleReader() ReadRule {
	address := os.Getenv("SQLCHAOS_HTTP")
	if address == "" {
		fmt.Fprintln(os.Stderr, "SQLCHAOS:http address not specified")
		return nil
	}

	rules := sync.Map{}
	handler := func(w http.ResponseWriter, r *http.Request) {
		key := r.URL.Path[1:]
		switch r.Method {
		case http.MethodDelete:
			rules.Delete(key)
			fmt.Fprint(w, "ok")
		case http.MethodPost:
			defer r.Body.Close()
			data, _ := ioutil.ReadAll(r.Body)
			rules.Store(key, data)
			fmt.Fprint(w, "ok")
		case http.MethodGet:
			var sb strings.Builder
			rules.Range(func(key, value interface{}) bool {
				sb.WriteString(key.(string))
				sb.WriteByte(':')
				sb.Write(value.([]byte))
				sb.WriteByte(',')
				return true
			})
			fmt.Fprint(w, sb.String())
		default:
			w.WriteHeader(401)
		}
	}
	http.HandleFunc("/", handler)

	echan := make(chan error)
	defer close(echan)
	var err error
	go func() { err = http.ListenAndServe(address, nil) }()
	<-time.After(time.Second)
	if err != nil {
		fmt.Fprintf(os.Stderr, "SQLCHAOS:http server starts failed:%s\n", err)
		return nil
	}
	fmt.Fprintf(os.Stderr, "SQLCHAOS:http server listening on %s\n", address)

	return func(ctx context.Context, dbname, table string) ([]byte, error) {
		fmt.Println(dbname + "/" + table)
		if v, ok := rules.Load(dbname + "/" + table); ok {
			return v.([]byte), nil
		}
		return nil, nil
	}
}
