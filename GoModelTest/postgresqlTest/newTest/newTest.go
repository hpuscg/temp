/*
#Time      :  2018/11/30 下午2:51 
#Author    :  chuangangshen@deepglint.com
#File      :  newTest.go
#Software  :  GoLand
*/
package main

import (
	"github.com/go-pg/pg"
	"fmt"
	"github.com/go-pg/pg/orm"
)

func main() {
	db := connect()
	err := CreateTable(db)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(db)
	}

	user1 := &User{
		Name: "user1",
		Emails: []string{"user1@234.com", "user1@123.com"},
	}

	db.Insert(user1)
}

const (
	addr     = "192.168.100.235:5432"
	user     = "postgres"
	passWard = "deepglint"
	dbName   = "libraT"
)

func connect() *pg.DB {
	db := pg.Connect(&pg.Options{
		User:     user,
		Password: passWard,
		Addr:     addr,
		Database: dbName,
	})
	var n int
	value, err := db.QueryOne(pg.Scan(&n), "SELECT 1")
	if err != nil {
		fmt.Println("clect to postgresql err: ", err)
	} else {
		fmt.Println("test selcet value is: ", value)
	}
	return db
}

type User struct {
	Id        int64
	Name      string
	Emails    []string `sql:"type:text[]"`
	tableName struct{} `sql:"users"`
}

type Story struct {
	Id        int64
	Title     string
	AuthorId  int64
	Author    *User
	tableName struct{} `sql:"story"`
}

func (s Story) String() string {
	return fmt.Sprintf("Story<%d %s %s>", s.Id, s.Title, s.Author)
}

func (u User) String() string {
	return fmt.Sprintf("User<%d %s %v>", u.Id, u.Name, u.Emails)
}

func CreateTable(db *pg.DB) error {
	for _, model := range []interface{}{&User{}, &Story{}} {
		err := db.CreateTable(model, &orm.CreateTableOptions{
			IfNotExists:   true,
			FKConstraints: true,
		})
		if err != nil {
			return err
		}
	}
	return nil
}
