package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

func main() {
	sendEvent()
}

var (
	eventUrl string
)

func sendEvent() {
	for {
		contents, err := ioutil.ReadFile("config.txt")
		if err != nil {
			panic(err.Error())
		} else {
			eventUrl = strings.Replace(string(contents), "\n", "", 1)
		}
		eventStruct := [...]string{
			`{"AlarmLevel":1,"EventType":611,"EventTypeProbability":0,"HotspotId":"","Id":"5b3c3dff6c385200010007c5","PeopleId":"","PlanetId":"undefined","SceneId":"undefined","SensorId":"a1f20f445035323133000004006f0118","StartTime":1530674686147,"TimeLength":1000,"UserId":""}`,
			`{"AlarmLevel":5,"EventType":329,"EventTypeProbability":1,"HotspotId":"undefined","Id":"5b3c3e2f6c385200010007ee","PeopleId":"049d7110f8fa4f86bb7cfcbb038a446c","PlanetId":"undefined","SceneId":"undefined","SensorId":"a1f20f445035323133000004006f0118","StartTime":1530674731868,"TimeLength":0,"UserId":""}`,
			`{"AlarmLevel":5,"EventType":319,"EventTypeProbability":0.99,"HotspotId":"undefined","Id":"5b3c3f4a6c38520001000819","PeopleId":"","PlanetId":"undefined","SceneId":"undefined","SensorId":"a1f20f445035323133000004006f0118","StartTime":1530675018509,"TimeLength":0,"UserId":""}`,
			`{"AlarmLevel":1,"EventType":633,"EventTypeProbability":0,"HotspotId":"","Id":"5b3c3e016c385200010007cd","PeopleId":"","PlanetId":"undefined","SceneId":"undefined","SensorId":"a1f20f445035323133000004006f0118","StartTime":1530674688147,"TimeLength":16692,"UserId":""}`,
		}
		for k, v := range eventStruct {
			fmt.Println(k)
			fmt.Println(v)
			respPost, err := doBytesPost(eventUrl, []byte(v))
			if err != nil {
				fmt.Println(err.Error())

			} else {
				res := string(respPost)
				fmt.Println(res)
			}
			time.Sleep(2 * time.Second)
		}
	}
}

func doBytesPost(url string, data []byte) ([]byte, error) {
	body := bytes.NewReader(data)
	request, err := http.NewRequest("POST", url, body)
	if err != nil {
		fmt.Println(err.Error())
		return []byte(""), err
	}
	request.Header.Set("Connection", "Keep-Alive")
	var resp *http.Response
	resp, err = http.DefaultClient.Do(request)
	if err != nil {
		fmt.Println(err.Error())
		return []byte(""), err
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err.Error())
		return []byte(""), err
	}
	return b, err
}
