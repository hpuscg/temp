package main

import (
	"strconv"
	"fmt"
	"github.com/pierrec/xxHash/xxHash64"
)

func main() {
	// str2int("admin")
	// ret := encode("admin12")
	// res := sfuint(ret)
	// fmt.Println(res)
	// weiYi()
	fmtTest()
}

func str2int(s string) {
	temp, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println(temp)
	}
}

func encode(s string) uint64 {
	ret := xxHash64.Checksum([]byte(s), 15120059285)
	return ret
}

func sfuint(u uint64) string {
	ret := strconv.FormatUint(u, 10)
	return ret
}

func weiYi() {
	x := 32 << 25
	fmt.Println(x / 8 / 1024 / 1024)
	fmt.Println(x / 8 / 1024)
	fmt.Println(128 * 1024)
}

type people struct {
	total int
	in    int
	out   int
}

func fmtTest() {
	var peopletmp people
	peopletmp.total = 21
	peopletmp.in = 2
	peopletmp.out = 3
	fmt.Printf("%+v", peopletmp)
	testAddress(&peopletmp)
}

func testAddress(peopleTemp *people) {
	fmt.Println(peopleTemp)
}
