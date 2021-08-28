package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

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

func main() {
	db, err := gorm.Open(mysql.Open(DSN), &gorm.Config{})
	if err != nil {
		fmt.Fprintf(os.Stderr, "connect db failed:%v", err)
		return
	}

	ctx := context.Background()

	email := "no-reply@gmail.com"
	now := time.Now()
	user := &User{
		ID:      1,
		Name:    "rick",
		Age:     50,
		Balance: 1024,
		Email:   &email,
		CreatedAt: &now,
	}
	if err := db.WithContext(ctx).Save(user).Error; err != nil {
		fmt.Fprintf(os.Stderr, "save record failed:%v\n", err)
		return
	}

	if err := db.WithContext(ctx).Where("name = ?", "rick").First(user).Error; err != nil {
		fmt.Fprintf(os.Stderr, "select record failed:%v\n", err)
		return
	}

	data, _ := json.Marshal(user)
	fmt.Println("user:", string(data))

	user.Balance = 2048
	if err := db.WithContext(ctx).Updates(user).Error; err != nil {
		fmt.Fprintf(os.Stderr, "select record failed:%v\n", err)
		return
	}
}
