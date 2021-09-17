package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	sqlchaos "github.com/u2386/gorm-sqlchaos"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

const DSN = "root:toor@tcp(127.0.0.1:3306)/dummy?charset=utf8mb4&parseTime=True&loc=Local"

type (
	User struct {
		ID        int64
		Name      string
		Age       int
		Email     *string
		Balance   int
		CreatedAt *time.Time
		UpdatedAt *time.Time
	}
)

func StaticRuleProvider(_ context.Context, dbname, table string) ([]byte, error) {
	return []byte(`{
		"dml": "UPDATE",
		"when": "name=rick",
		"then": "age=25"
	}`), nil
}

func main() {
	db, err := gorm.Open(mysql.Open(DSN), &gorm.Config{}, &sqlchaos.Config{
		DBName:     "dummy",
		RuleProvider: StaticRuleProvider,
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "connect db failed:%v", err)
		return
	}

	ctx := context.Background()

	user := &User{}
	err = db.WithContext(ctx).Where("name = ?", "rick").First(user).Error
	if err != nil {
		if err != gorm.ErrRecordNotFound {
			fmt.Fprintf(os.Stderr, "select record failed:%v\n", err)
			return
		}

		email := "no-reply@gmail.com"
		user = &User{
			ID: 1,
			Name: "rick",
			Age: 50,
			Email: &email,
			Balance: 1024,
		}
		if err = db.WithContext(ctx).Save(user).Error; err != nil {
			fmt.Fprintf(os.Stderr, "save record failed:%v\n", err)
			return
		}
	}

	data, _ := json.Marshal(user)
	fmt.Println("user:", string(data))

	user.Balance = 2048
	if err := db.Debug().WithContext(ctx).Model(User{}).
		Where("id = ?", user.ID).
		Updates(user).
		Error; err != nil {
		fmt.Fprintf(os.Stderr, "select record failed:%v\n", err)
		return
	}
}
