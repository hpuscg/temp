/*
#Time      :  2021/1/22 1:52 下午
#Author    :  chuangangshen@deepglint.com
#File      :  sqlTest.go
#Software  :  GoLand
*/
package main

import (
	"fmt"
	"temp/GoModelTest/sqlTest/config"
	"temp/GoModelTest/sqlTest/model"
)

func main() {
	err := config.InitDb()
	if err != nil {
		fmt.Println(err)
	}

	var products []model.Product
	query := config.DB

	// query = query.Where("product_id in (select id from product where product.name like ?)", "%haomuT%")
	// query = query.Where("id = ?", 3).Joins("inner join license on license.product_id = product.id").Find(&products)
	query.Joins("left join license as l on product.id = l.product_id").Find(&products)
	fmt.Println(len(products))
	for _, pro := range products {
		fmt.Printf("%+v\n", pro)
		return
	}
}
