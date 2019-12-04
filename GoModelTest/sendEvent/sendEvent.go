package main

import (
	"fmt"
	"io/ioutil"
	"strings"
	"time"
	"net"
	"net/http"
	"encoding/json"
	"bytes"
)

func main()  {
	// num := 1
	sendEvent()
}

func sendEvent()  {
	// num := 1
	for true {
		var ip string
		contents, err := ioutil.ReadFile("config.txt")
		if err != nil {
			fmt.Println(err.Error())
		}else {
			ip = strings.Replace(string(contents),"\n","",1)
		}
		conn, err := net.DialTimeout("tcp", ip + ":8008", 2*time.Second)
		if (err != nil) || (conn == nil) {
			fmt.Println("this is a error ip, please checkout config.txt")
			return
		}

		time.Sleep(2 * time.Second)
		url := `http://` + ip + `:8008/api/serviceconfig?keys=["/config/libra/data/enable_color_tracking","/config/eventserver/full_video_storage_ttl_days","/config/eventserver/tss_ttl_days","/config/eventserver/pub_vibo_url","/config/eventserver/pub_db_url","/config/libra/sensor/depth_MoG_factor"]`
		respGet, err := http.Get(url)
		defer respGet.Body.Close()
		result, err := ioutil.ReadAll(respGet.Body)
		var f interface{}
		json.Unmarshal(result, &f)
		m := f.(map[string]interface{})

		var vidoUrl string
		for k, v := range m {
			if "/config/eventserver/pub_vibo_url" == k {
				urlBefore := v.(string)
				if strings.HasPrefix(urlBefore, "http://") {
					vidoUrl = urlBefore
				}else {
					vidoUrl = `http://` + urlBefore
				}
			}
		}
		eventStruct := [...]string {
			1: `{"AlarmLevel":1,"EventType":611,"EventTypeProbability":0,"HotspotId":"","Id":"5b3c3dff6c385200010007c5","PeopleId":"","PlanetId":"undefined","SceneId":"undefined","SensorId":"a1f20f445035323133000004006f0118","StartTime":1530674686147,"TimeLength":1000,"UserId":""}`,
			2: `{"AlarmLevel":5,"EventType":329,"EventTypeProbability":1,"HotspotId":"undefined","Id":"5b3c3e2f6c385200010007ee","PeopleId":"049d7110f8fa4f86bb7cfcbb038a446c","PlanetId":"undefined","SceneId":"undefined","SensorId":"a1f20f445035323133000004006f0118","StartTime":1530674731868,"TimeLength":0,"UserId":""}`,
			3: `{"AlarmLevel":5,"EventType":319,"EventTypeProbability":0.99,"HotspotId":"undefined","Id":"5b3c3f4a6c38520001000819","PeopleId":"","PlanetId":"undefined","SceneId":"undefined","SensorId":"a1f20f445035323133000004006f0118","StartTime":1530675018509,"TimeLength":0,"UserId":""}`,
			4: `{"AlarmLevel":1,"EventType":633,"EventTypeProbability":0,"HotspotId":"","Id":"5b3c3e016c385200010007cd","PeopleId":"","PlanetId":"undefined","SceneId":"undefined","SensorId":"a1f20f445035323133000004006f0118","StartTime":1530674688147,"TimeLength":16692,"UserId":""}`,
		}
		for _, v := range eventStruct {
			respPost, err := doBytesPost(vidoUrl, []byte(v))
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




