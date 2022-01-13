package main

import "fmt"

func main() {
	var A = []int{1, 2, 3, 4, 5, 6, 7, 16, 16, 3, 5, 16, 16, 16, 16, 16, 16, 16, 16, 16, 16, 16, 16, 16}
	var target = 100
	ret := backPackV(A, target)
	fmt.Println(ret)
}

func backPackV(nums []int, target int) int {
	// write your code here
	count := 0
	length := len(nums)
	for i := 0; i < length; i++ {
		for j := 0; j < length-i; j++ {
			if target == sunCount(nums[j:], i) {
				count++
			}
		}
	}
	return count
}

func sunCount(nums []int, sept int) int {
	if sept == 0 {
		return nums[sept]
	} else {
		return nums[sept] + sunCount(nums[1:], sept-1)
	}
}
