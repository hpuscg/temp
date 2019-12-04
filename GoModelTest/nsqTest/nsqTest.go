/*
#Time      :  2019/1/13 下午3:50 
#Author    :  chuangangshen@deepglint.com
#File      :  nsqTest.go
#Software  :  GoLand
*/
package main

import (
	"github.com/deepglint/go-nsq"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"io/ioutil"
)

func main() {
	// nsqPublish()
	httpPublish()
}

// public to nsq topic
func nsqPublish() {
	var (
		err          error
		mqWriterVibo *nsq.Producer = nil
	)
	wcfg := nsq.NewConfig()
	ViboEventTopicServer := "192.168.100.235:4151"

	mqWriterVibo, err = nsq.NewProducer(ViboEventTopicServer, wcfg)
	if err != nil {
		fmt.Println(err)
	}
	defer mqWriterVibo.Stop()
	ei := "test"
	bs, err := json.Marshal(ei)
	pubTopic := "vibo_events"
	if mqWriterVibo != nil {
		err = mqWriterVibo.Publish(pubTopic, bs)
		fmt.Println("err is:", err)
		fmt.Println("yes")
	}
}

func httpPublish() {
	url := "http://192.168.100.235:4151/pub?topic=vibo_events"
	// ret, _ := json.Marshal("")
	// resp, _ := io.Reader().Read(ret)

	request, _ := http.NewRequest("post", url, strings.NewReader("yes"))
	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		fmt.Println("post data err: ", err)
	} else {
		fmt.Println("post data sussess")
		respBody,_ := ioutil.ReadAll(resp.Body)
		fmt.Println("response data is: ", string(respBody))
	}
}
