package main

import (
	"fmt"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	db *gorm.DB
)

func main() {
	InitDb()
	OneTooMany()
}

type User struct {
	ID    uint   `json:"id" gorm:"primaryKey"`
	Cards []Card `json:"cards" gorm:"foreignKey:UserId"`
	Age   uint
}

type Card struct {
	UserId uint
	CardId uint `json:"cardId" gorm:"primaryKey"`
	Name   int
}

func InitDb() {
	config := gorm.Config{}
	config.PrepareStmt = true
	config.QueryFields = true
	config.AllowGlobalUpdate = true
	config.FullSaveAssociations = false
	if con, err := gorm.Open(sqlite.Open("/Users/hpu_scg/gocode/src/temp/GoModelTest/gorm_v2/oneTooMany/config/test.db"), &config); err != nil {
		panic(err)
	} else {
		db = con
		db.Logger = db.Logger.LogMode(4)
	}
	db.AutoMigrate(User{}, Card{})
}

func OneTooMany() {
	user := User{
		Age: 18,
	}
	db.Create(&user)
	for i := 0; i < 2; i++ {
		card1 := Card{
			Name: i,
		}
		user.Cards = append(user.Cards, card1)
	}
	db.Save(&user)
	var tempUser User
	db.Where("id = ?", user.ID).Preload("Cards").Find(&tempUser)
	fmt.Printf("%+v", tempUser)
}
