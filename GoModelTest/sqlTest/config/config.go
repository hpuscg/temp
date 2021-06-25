/*
#Time      :  2021/1/22 1:55 下午
#Author    :  chuangangshen@deepglint.com
#File      :  config.go
#Software  :  GoLand
*/
package config

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

var (
	DB *gorm.DB
)

func InitDb() (err error) {
	driver := "mysql"
	desc := "root:Deepglint123@(123.59.135.181:3306)/lic_manager?charset=utf8&parseTime=True&loc=UTC"

	DB, err = gorm.Open(driver, desc)
	if err != nil {
		fmt.Println(err.Error())
	}
	DB.SingularTable(true)
	return
}
