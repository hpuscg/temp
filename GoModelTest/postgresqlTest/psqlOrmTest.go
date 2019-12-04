package main

import (
	"fmt"
	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
)

func main() {
	ExampleDBModel()
}

type User struct {
	Id     int64
	Name   string
	Emails []string `sql:",type:text[]" pg:",array"`
}

func (u User) String() string {
	return fmt.Sprintf("User<%d %s %v>", u.Id, u.Name, u.Emails)
}

type Story struct {
	Id       int64
	Title    string
	AuthorId int64
	Author   *User
}

func (s Story) String() string {
	return fmt.Sprintf("Story<%d %s %s>", s.Id, s.Title, s.Author)
}

func ExampleDBModel() {
	db := pg.Connect(&pg.Options{
		User: "postgres",
		Addr: "192.168.100.235:5432",
	})
	fmt.Println(db.String())
	defer db.Close()

	err := createSchema(db)
	if err != nil {
		fmt.Println("create schema err: ", err)
	}

	user1 := &User{
		Name:   "admin",
		Emails: []string{"admin1@admin", "admin2@admi"},
	}
	err = db.Insert(user1)
	if err != nil {
		fmt.Println("insert user1 err :", err)
	}

	err = db.Insert(&User{
		Name:   "root",
		Emails: []string{"root1@root", "root2@root"},
	})
	if err != nil {
		fmt.Println("User err :", err)
	}

	story1 := &Story{
		Title:    "cool story",
		AuthorId: user1.Id,
	}
	err = db.Insert(story1)
	if err != nil {
		fmt.Println("story1 err :", err)
	}

	user := &User{Id: user1.Id}
	err = db.Select(user)
	if err != nil {
		fmt.Println("selcet user err :", err)
	}

	var users []User
	err = db.Model(&users).Select()
	if err != nil {
		fmt.Println("select users err: ", err)
	}

	story := new(Story)
	err = db.Model(story).Relation("Author").Where("story.id = ?", story1.Id).Select()
	if err != nil {
		fmt.Println("select story err :", err)
	}

	fmt.Println("user is: ", user)
	fmt.Println("users is: ", users)
	fmt.Println("story is: ", story)
}

func createSchema(db *pg.DB) error {
	for _, model := range []interface{}{(*User)(nil), (*Story)(nil)} {
		err := db.CreateTable(model, &orm.CreateTableOptions{
			Temp: true,
		})
		if err != nil {
			fmt.Println("55", err)
			return err
		}
	}
	return nil
}
