package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	sqlchaos "github.com/u2386/gorm-sqlchaos"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

const DSN = "root:toor@tcp(127.0.0.1:3306)/dummy?charset=utf8mb4&parseTime=True&loc=Local"

type (
	User struct {
		ID        int64      `json:"id"`
		Name      string     `json:"name"`
		Age       int        `json:"age"`
		Email     *string    `json:"email"`
		Balance   int        `json:"balance"`
		CreatedAt *time.Time `json:"created_at"`
		UpdatedAt *time.Time `json:"updated_at"`
	}
)

func main() {
	db, err := gorm.Open(mysql.Open(DSN), &gorm.Config{}, &sqlchaos.Config{
		DBName:     "dummy",
		RuleReader: sqlchaos.WithSimpleHTTPRuleReader(),
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "connect db failed:%v", err)
		return
	}

	handler := func(w http.ResponseWriter, r *http.Request) {
		userID := r.URL.Path[len("/users/"):]
		switch r.Method {
		case http.MethodGet:
			user := &User{}
			if err := db.Debug().Model(User{}).Where("id = ?", userID).First(user).Error; err != nil {
				w.WriteHeader(400)
				fmt.Fprintf(w, "err:%v", err)
				return
			}
			data, _ := json.Marshal(user)
			fmt.Fprint(w, string(data))

		case http.MethodPut:
			user := &User{}
			defer r.Body.Close()
			data, _ := ioutil.ReadAll(r.Body)
			if err := json.Unmarshal(data, user); err != nil {
				w.WriteHeader(400)
				fmt.Fprintf(w, "err:%v", err)
				return
			}

			if err := db.Debug().Model(User{}).Where("id = ?", userID).Updates(user).Error; err != nil {
				w.WriteHeader(400)
				fmt.Fprintf(w, "err:%v", err)
				return
			}
			fmt.Fprint(w, "ok")

		case http.MethodPost:
			defer r.Body.Close()
			data, _ := ioutil.ReadAll(r.Body)
			names := strings.Split(string(data), ",")

			users := []User{}
			for _, name := range names {
				users = append(users, User{
					Name: name,
					Age: 24,
				})
			}
			if err := db.Debug().Model(User{}).Create(&users).Error; err != nil {
				w.WriteHeader(400)
				fmt.Fprintf(w, "err:%v", err)
				return
			}
			fmt.Fprint(w, "ok")

		default:
			w.WriteHeader(401)
			return
		}
	}
	http.HandleFunc("/users/", handler)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}
