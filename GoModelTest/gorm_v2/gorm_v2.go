package main

import (
	"fmt"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var db *gorm.DB
var err error

func main() {
	db, err = gorm.Open(sqlite.Open(""), &gorm.Config{})
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(db)
}
