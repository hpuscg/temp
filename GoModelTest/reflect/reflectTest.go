/*
#Time      :  2019/12/24 4:59 PM 
#Author    :  chuangangshen@deepglint.com
#File      :  reflectTest.go
#Software  :  GoLand
*/
package main

import (
	"fmt"
	"reflect"
)

func main() {
	ReflectTest()
}

type People struct {
	name string
	age  int
}

func ReflectTest() {
	people := People{
		"zhangSan",
		28,
	}
	fmt.Println(reflect.TypeOf(people))
	fmt.Println(reflect.ValueOf(&people).Elem().CanSet())
}
