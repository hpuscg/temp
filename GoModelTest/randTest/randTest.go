package main

import (
	"math/rand"
	"fmt"
	"time"
)

func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	tempInt := GetRandNum(intNum)
	tempFloat := GetRandFloat()
	fmt.Println("tempInt is : ", tempInt)
	fmt.Println("tempFloat is : ", tempFloat)
}

const intNum  = 100

func GetRandNum(num int)  int {
	randNum := rand.Intn(num)
	return randNum
}

func GetRandFloat() float64 {
	randFloat := rand.Float64()
	return randFloat
}
