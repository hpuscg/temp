package main

import "fmt"


func main() {
	a := 3
	if a > 1 {
		fmt.Println("a是大于1的！")
		if a < 4 {
			fmt.Println("a是小于4的！")
		}
	}else {
		fmt.Println("a是小于1的！")
	}
}
