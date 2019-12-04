package main

import (
	"os/exec"
	"fmt"
	"reflect"
)

func main() {
	c := "ps -a |grep main.go |wc -l "
	cmd := exec.Command(`sh`, ` -c `, c)

	out, err := cmd.Output()

	if err != nil {
		fmt.Println("1111")
		fmt.Println(err.Error())
	}

	fmt.Println("2222")
	fmt.Println(reflect.TypeOf(out))

}
