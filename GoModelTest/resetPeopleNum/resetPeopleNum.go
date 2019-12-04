/*
#Time      :  2019/5/12 下午12:38 
#Author    :  chuangangshen@deepglint.com
#File      :  resetPeopleNum.go
#Software  :  GoLand
*/
package main

import (
	"flag"
	"fmt"
	"strings"
	"net/http"
	"encoding/json"
)

func main() {
	var (
		Ip string
	)
	flag.StringVar(&Ip, "ip", "", "")
	flag.Parse()
	fmt.Println("ip is:", Ip)
	resetPeopleNum(Ip)
}

//
func resetPeopleNum(ip string)  {
	postData := RealPeopleNum{
		CashPeopleNum:0,
		BursePeopleNum:0,
	}
	bodyUrl := "http://" + ip + ":1357/api/combination/realPeopleNum"
	data, err:= json.Marshal(postData)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(data))
	bodyBuf := strings.NewReader(string(data))
	contentType := "application/json"
	resp, err := http.Post(bodyUrl , contentType, bodyBuf)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(resp.Status)
}

type RealPeopleNum struct {
	CashPeopleNum  int `json:"cash_people_num"`
	BursePeopleNum int `json:"burse_people_num"`
}


