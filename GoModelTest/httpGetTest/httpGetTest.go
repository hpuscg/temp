/*
#Time      :  2018/12/21 下午6:55 
#Author    :  chuangangshen@deepglint.com
#File      :  httpGetTest.go
#Software  :  GoLand
*/
package main

import (
	"fmt"
	"strings"
	"encoding/json"
	"io/ioutil"
	"net/http"
	httpS "github.com/deepglint/flowservice/util/http"
	"flag"
	"github.com/deepglint/flowservice/models"
)

func main() {
	var ip string
	flag.StringVar(&ip, "ip", "192.168.5.178:8888", "event server ip")
	flag.Parse()
	// t := Teacher{}
	// t.ShowA()
	// getToken(ip)
	// PostEvent(ip)
	// GetSensorId()
	// GetRealPeopleNum("192.168.19.247")
	// GetRGBValue()
	// getFromT3()
	getPackageVersionFromServer()
}

func getPackageVersionFromServer() {
	serverIp := "192.168.100.235"
	url := "http://" + serverIp + ":8008/api/sensor_version"
	resp, err := http.Get(url)
	if err != nil || resp == nil {
		fmt.Println(err)
		return
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(data))
}

func getFromT3() {
	baseUrl := "http://192.168.100.170/api/l/CapsuleConfig"
	// url := strings.Replace(baseUrl, "IP", ip, -1)
	result, err := httpS.GetWithToken(baseUrl, models.LibraToke, "application/json")
	if err != nil {
		fmt.Println(err)
		return
	}
	var ret models.ResponseData
	err = json.Unmarshal(result, &ret)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("%+v\n", ret)
	switch ret.Data.(type) {
	case string:
		fmt.Printf("%+v\n", ret.Data.(string))
	case map[string]interface{}:
		fmt.Printf("%+v\n", ret.Data.(map[string]interface{}))
	}
}


/*func main() {
	for true {
		GetPeoplePoint()
		time.Sleep(1 * time.Second)
	}
	fmt.Println("1122")
}

func GetPeoplePoint()  {
	url := "http://www.baidu.com"
	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("get people piont from %s err : %s", url, err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("read resp body err:", err)
	}
	fmt.Println(string(body))
}*/

type People struct{}

func (p *People) ShowA() {
	fmt.Println("showA")
	p.ShowB()
}
func (p *People) ShowB() {
	fmt.Println("showB")
}

type Teacher struct {
	People
}

func (t *Teacher) ShowB() {
	fmt.Println("teacher showB")
}

func GetRGBValue() {
	url := "http://192.168.5.251:8008/api/rgb/current"
	result, err := httpS.HTTPGet(url)
	if err != nil {
		// glog.V(0).Infoln(err)
		fmt.Println(err)

	}
	fmt.Println(string(result))
}

func GetRealPeopleNum(ip string) {
	url := "http://" + ip + ":1357/api/realnumber"
	result, err := httpS.HTTPGet(url)
	if err != nil {
		// glog.V(0).Infoln(err)
		fmt.Println(err)

	}
	fmt.Println(string(result))
}

func GetSensorId() {
	url := "http://192.168.100.223:8008/api/sensorid"
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("111", err)
	}
	ret, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("2222", err)
	}
	fmt.Println("4444", string(ret))
}

func CookieTest() {
	/*http.Cookie{

	}*/

	// req, _ := http.NewRequest("POST", "url", "")
	// req.Header.Add("token", "data")
}

var token string

// const Token  = "Bearer eyJhbGciOiJIUzUxMiJ9.eyJzdWIiOiJoYW9tdVQiLCJhdXRob3JpdGllcyI6IltdIiwiZXhwIjoxNTUyNDgyNjgzfQ.JA4N4ASM0h6RX8UtzucJ9EIntHhxnjkkwMdrYl0wGojRpHPvvQ_t5TAaWQVcPccvfbMda1OxhV9Hm5aB6ki8sg"
const Token  = ""

func PostEvent(ip string) {
	code, err := postEvent(ip, Token)
	if err != nil {
		fmt.Println("==113===", err)
		return
	}
	if 401 == code {
		getToken(ip)
		_, err = postEvent(ip, token)
		if err != nil {
			fmt.Println("===120===", err)
		}
	}
}

func getToken(ip string) {
	addr := "http://" + ip +"/api/eventserver/authenticate"
	userName := "haomuT"
	passWord := "abc@Dgsh"

	data := make(map[string]string)
	data["username"] = userName
	data["password"] = passWord
	byteData, err := json.Marshal(data)
	if err != nil {
		fmt.Println("81 err: ", err)
	}
	resp, err := http.Post(addr, "application/json", strings.NewReader(string(byteData)))
	if err != nil {
		fmt.Println("85 err: ", err)
	}
	defer resp.Body.Close()
	respData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("93 err: ", err)
	}
	ret := make(map[string]interface{})
	err = json.Unmarshal(respData, &ret)
	if err != nil {
		fmt.Println("98 err: ", err)
	}
	token = ret["data"].(string)
	fmt.Println(string(respData))
	fmt.Println(token)
}

type HttpEvent struct {
	DeviceId    string      `json:"deviceId"`
	DeviceType  int      `json:"deviceType"`
	EventType   int         `json:"eventType"`
	EventTime   string         `json:"eventTime"`
	EventDetail interface{} `json:"eventDetail"`
}

func postEvent(ip, Tk string) (code float64, err error) {
	addr := "http://" + ip + "/api/eventserver/eventUpload"
	client := &http.Client{}
	postData, err := json.Marshal(HttpEvent{
		DeviceId: "123qweasdzxc",
		DeviceType: 1,
		EventType: 119,
	})
	if err != nil {
		fmt.Println("120 err :", err)
	}
	req, _ := http.NewRequest("POST", addr, strings.NewReader(string(postData)))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("authorization", Tk)
	response, err := client.Do(req)
	if err != nil {
		fmt.Println("127 err :", err)
	}
	defer response.Body.Close()
	respData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println("132 err :", err)
	}
	dataTmp := make(map[string]interface{})
	err = json.Unmarshal(respData, &dataTmp)
	if err != nil {
		fmt.Println(err)
	}
	code = dataTmp["code"].(float64)
	fmt.Println(string(respData))
	fmt.Println(code)
	fmt.Println(code == 401)
	return
}
