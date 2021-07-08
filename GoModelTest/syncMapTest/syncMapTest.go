package main

import (
	"fmt"
	"strconv"
	"strings"
	"sync"
)

func main() {
	// syncMap()
	stringTest()
}

func syncMap() {
	var count int
	var mapTest = sync.Map{}

	mapTest.Store("1", "3")
	mapTest.Store("2", "4")
	mapTest.Range(func(key, value interface{}) bool {
		count += 1
		return true
	})
	fmt.Println(count)
}

func strByXOR(message string, keywords string) string {
	messageLen := len(message)
	keywordsLen := len(keywords)

	result := ""

	for i := 0; i < messageLen; i++ {
		result += string(message[i] ^ keywords[i%keywordsLen])
	}
	return result
}

func stringTest() {
	// a := "a"
	// fmt.Printf("%#x\n", a)

	// str := "02 01"
	baseStr := "FE 0B 02 01 02 00 00 00 04 3C CC"
	strList := strings.Split(baseStr, " ")
	cutStrList := strList[:len(strList)-1]
	fmt.Println(cutStrList)
	var result uint64
	for _, value := range cutStrList {
		n, err := strconv.ParseUint(value, 16, 64)
		if err != nil {
			fmt.Println(err)
			return
		}
		result = result ^ n
	}

	fmt.Println(result)
	fmt.Printf("%X\n", result)
	retData := fmt.Sprintf("%X", result)
	fmt.Println(retData)
	var postData string
	for _, value := range strList {
		retValue := fmt.Sprintf("%X", value)
		strings.Join([]string{
			postData,
			retValue,
		}, "")
	}

	/*hexA := hex.EncodeToString([]byte(a))

	fmt.Println(hexA)*/

	// fmt.Printf("%v", b)

	/*for key, value := range a {
		fmt.Println(key, string(value))
	}*/
}