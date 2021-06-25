/*
#Time      :  2020/12/17 5:21 下午
#Author    :  chuangangshen@deepglint.com
#File      :  main.go
#Software  :  GoLand
*/
package main

import (
	"fmt"
	"github.com/casbin/casbin/v2"
	"math"
)

func main() {
	// CasbinTest()
	CeilTest()
}

func CeilTest() {
	a := 0.3 - 0.5
	b := 0.3
	fmt.Println(math.Ceil(a), math.Ceil(b))
}

func CasbinTest() {
	e, err := casbin.NewEnforcer("./config/perm.conf", "./config/policy.csv")
	if err != nil {
		fmt.Println(err)
	}

	subs := []string{"bob", "zeta"}
	objs := []string{"data1", "data2"}
	acts := []string{"read", "write"}

	for _, sub := range subs {
		for _, obj := range objs {
			for _, act := range acts {
				ok, _ := e.Enforce(sub, obj, act)
				fmt.Println(sub, obj, act, "=", ok)
			}
		}
	}
}
