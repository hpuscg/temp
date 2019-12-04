package main

import "fmt"

func main() {
	ec()
}



func ca() error {
	fmt.Println("111")
	return nil
}

func ec() error {
	i := 1
	if i != 2 {
		ca()
	}
	fmt.Println("222")
	return nil
}

