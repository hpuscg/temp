package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
)

func main() {
	// test1()
	// strSplit()
	// bufferToInterface()
	// byteToString()
	// justTest()
	// strIntTest()
	stringSplit()
	// preTest()
	// stringToIntTest()
	// FilePathJoinTest()
	// intTest()
}

func intTest() {
	a := 1
	b := 2
	fmt.Println(float64(a) / float64(b))
}

func FilePathJoinTest() {
	path := filepath.Join("/data/", "123")
	path2 := filepath.Join("/data", "123")
	fmt.Println(path)
	fmt.Println(path2)
}

func stringToIntTest() {
	var data interface{}
	var ret int
	var err error
	var tmp int32 = 123
	data = tmp
	switch data.(type) {
	case string:
		ret, err = strconv.Atoi(data.(string))
		if err != nil {
			fmt.Println(err)
		}
	case int:
		ret = data.(int)
	case int32:
		ret = int(data.(int32))
	}
	fmt.Println(ret)
}

func preTest() {
	// str := "http://123"
	fmt.Println(strings.HasPrefix("brightchanged_legacy", "brightchanged"))
}

func stringSplit() {
	s := "rtsp://192.168.12.12/libra/www"
	// var ret string
	// ss := strings.SplitN(s, "/", -1)
	n := strings.LastIndex(s, "/")
	pre := s[:n]
	fmt.Println(pre)
	ss := s[n+1:]
	fmt.Println(ss)
	ss1 := strings.SplitAfterN(s, "/", 2)
	fmt.Println(ss1)
	//strings.SplitAfterN()
	/* for i, data := range ss {
		fmt.Printf("%d=%s\n", i, data)
		ret = strings.Join([]string{
			ret,
			data,
		}, "/")
	}
	ret1 := strings.Trim(ret, "/")
	fmt.Println(ret1) */
}

func strIntTest() {
	// strTime := strconv.Itoa(-24 * 2) + "h"
	// fmt.Println(strTime)
	fmt.Println(strings.HasPrefix("1234", "34"))
}

func strSplit() {
	a := "shen.chuan"
	b := strings.Split(a, ".")[0]
	c := strings.Split(a, b+".")
	fmt.Println(c[1])
}

func test1() {
	fmt.Println("Contains:", strings.Contains("shenchuangang", "sh"))
	fmt.Println("Contains:", strings.Contains("shenchuangang", "df"))
	fmt.Println("Count:", strings.Count("shenchaugnang", "c"))
	fmt.Println("HasPrefix:", strings.HasPrefix("shenchuangang", "sh"))
	fmt.Println("HasSuffix:", strings.HasSuffix("shenchuangang", "ang"))
	fmt.Println("Index:", strings.Index("等待", "待"))
	fmt.Println("Replace:", strings.Replace("shenchuangang", "n", "N", -1))
	fmt.Println("Replace:", strings.Replace("shenchuangang", "n", "N", 1))
	// a := "[552, 551, 139, 130, 129, 120, 740, 741, 742, 743, 720, 721, 722, 723, 630, 631, 632, 633, 550, 551, 552, 553, 540, 541, 542, 543, 1000, 1001, 1002, 1003]"
	// fmt.Println(a)
	// fmt.Println(len(a))
	/*
		// a := "[552:false 742:false 543:false 721:false 743:false 130:false 1001:false 631:false 722:false 120:false 550:false 632:false 129:false 139:false 630:false 540:false 553:false 633:false 723:false 1000:false 740:false 1003:false 541:false 542:false 1002:false 720:false 741:false 551:false]"
		a := "[\"552\":false, \"742\":false]"
		var sv map[string]bool
		err := json.Unmarshal([]byte(a), &sv)
		if err != nil {
			fmt.Println(err.Error())
		} else {
			fmt.Println(sv)
		}
	*/

	/*
			m := map[string]string{
		        "LOG_LEVEL": "DEBUG",
		        "API_KEY":   "12345678-1234-1234-1234-1234-123456789abc",
		    }
			println(createKeyValuePairs(m))
	*/
	mapvalue := make(map[string]string)
	mapvalue["wer"] = "123"
	mapvalue["we"] = "234"
	// map2str := mapvalue.(string)
	// fmt.Println("wer is :", mapvalue["wer"])
	// fmt.Println("none is :", mapvalue["none"])
	fmt.Println(mapvalue)
	// fmt.Println(map2str)
	fmt.Println("===============")
	reStr := toLower("Post")
	fmt.Println(reStr)
	fmt.Println("--------")
	upStr := toUp("post")
	fmt.Println(upStr)
	fmt.Println("&&&&&&&&&&&&&&&&")
	allUpStr := toAllUp("post")
	fmt.Println(allUpStr)
}

func createKeyValuePairs(m map[string]string) string {
	b := new(bytes.Buffer)
	for key, value := range m {
		fmt.Fprintf(b, "%s=\"%s\"\n", key, value)
	}
	return b.String()
}

func toLower(str string) string {
	resultStr := strings.ToLower(str)
	return resultStr
}

func toUp(str string) string {
	resultStr := strings.Title(str)
	return resultStr
}

func toAllUp(str string) string {
	resultStr := strings.ToUpper(str)
	return resultStr
}

func bufferToInterface() {
	f, err := os.Open("data.json")
	if err != nil {
		fmt.Println(err.Error())
	}
	defer f.Close()
	buf := new(bytes.Buffer)
	buf.ReadFrom(f)
	fmt.Println(buf.String())
	var ret map[string]interface{}
	json.Unmarshal(buf.Bytes(), &ret)
	fmt.Println(ret["data"])
	ret2 := ret["data"]
	// str := ret2.(string)
	fmt.Println(reflect.TypeOf(ret2))
	rut, _ := json.Marshal(ret2)
	fmt.Println(reflect.TypeOf(rut))
	// data := *(*[]byte)(unsafe.Pointer(ret2))
	// fmt.Println(str)
	// buf2 := bytes.NewBuffer(ret["data"])
}

func byteToString() {
	type data1 struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}
	data := data1{
		Name: "lisi",
		Age:  12,
	}
	fmt.Println(data)
	rel, err := json.Marshal(data)
	if err != nil {
		fmt.Println("err is ", err)
	}
	relStr := string(rel)
	fmt.Println("rel is ", relStr)
	fmt.Println(reflect.TypeOf(relStr))
	var data2 data1
	err = json.Unmarshal([]byte(string(rel)), &data2)
	fmt.Println("err2 is", err)
	fmt.Println(data2)
}

func justTest() {
	data := 1 == 2
	fmt.Println(data)
}
