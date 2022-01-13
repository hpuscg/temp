package main

import "fmt"

func main() {
	var A = []int{1, 1, 2, 2, 3, 4, 4}
	ret := singleNumber(A)
	fmt.Println(ret)
}

func singleNumber(A []int) int {
	// write your code here
	var B []int
	var ret int
	n := len(A) - 1
	i := 0
	for i <= n {
		var count = 0
		data := A[i]
		for _, data2 := range B {
			if data2 == data {
				count++
				break
			}
		}
		if count == 0 {
			fmt.Println("data: ", data)
			j := i + 1
			for j <= n {
				if A[j] == data {
					B = append(B, data)
					count++
					break
				}
				j++
			}
		}
		fmt.Println("count: ", count)
		if count == 0 {
			fmt.Println("data: ", data)
			ret = data
			break
		}
		i++
	}
	return ret
}
