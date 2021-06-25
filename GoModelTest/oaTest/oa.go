package main

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/xml"
	"github.com/hooklift/gowsdl/soap"
	"gitlab.deepglint.com/junkaicao/glog"
	"io/ioutil"
	"net"
	"net/http"
	"temp/GoModelTest/oaTest/myservice"
	"time"
)

var (
	client *soap.Client
	// clientType myservice.WorkflowServicePortType
)

func main() {
	glog.Config(glog.WithAlsoToStd(true), glog.WithLevel("info"))
	client = soap.NewClient("http://192.168.2.4:8899/services/WorkflowService")
	// clientType = myservice.NewWorkflowServicePortType(client)
	// GetUserId()
	// GetTodoWorkList()
	// GetCreteTypeList()
	// GetTodoWorkListCount()
	getUserIdWithSoap()
}

/*func GetCreteTypeList() {
	arrayData := myservice.ArrayOfString{Astring: []string{" typename like '%人%' "}}
	requestData := myservice.GetCreateWorkflowList{
		In0: 0,
		In1: 20,
		In2: 10,
		In3: 22,
		In4: 0,
		In5: &arrayData,
	}
	responseData, err := workFlowType.GetCreateWorkflowList(&requestData)
	if err != nil {
		glog.Infoln(err)
	}
	glog.Infof("%+v", *responseData)
}*/

func GetUserId() {

	requestData := myservice.GetUserId{
		In0: "chuangangshen",
		In1: "901010",
	}
	// responseData := myservice.GetUserIdResponse{}
	var (
		responseData interface{}
	)
	err := client.CallContext(context.Background(), "urn:weaver.workflow.webservices."+
		"WorkflowService.getUserid", requestData, &responseData)

	if err != nil {
		glog.Infoln(err)
	}
	glog.Infof("%+v", responseData)
}

/*func GetTodoWorkList() {
	requestData := myservice.GetToDoWorkflowRequestList{
		In0: 0,
		In1: 10,
		In2: 10,
		In3: 20,
		In4: nil,
	}
	responseData, err := workFlowType.GetToDoWorkflowRequestList(&requestData)
	if err != nil {
		glog.Infoln(err)
	}
	glog.Infof("%+v", *responseData)
}*/

func GetTodoWorkListCount() {
	requestData := myservice.GetToDoWorkflowRequestCount{
		In0: 100,
		In1: &myservice.ArrayOfString{Astring: []string{}},
	}
	responseData := myservice.GetToDoWorkflowRequestCountResponse{}
	err := client.CallContext(context.Background(), "urn:weaver.workflow.webservices.WorkflowService."+
		"getToDoWorkflowRequestCount", requestData, &responseData)
	if err != nil {
		glog.Errorln(err)
		return
	}
	glog.Infof("%+v", responseData)
}

func getUserIdWithSoap() {
	envelope := soap.SOAPEnvelope{}
	requestData := myservice.GetUserId{
		In0: "chuangangshen",
		In1: "901010",
	}
	envelope.Body.Content = requestData
	buffer := new(bytes.Buffer)
	var encoder soap.SOAPEncoder
	encoder = xml.NewEncoder(buffer)

	if err := encoder.Encode(envelope); err != nil {
		glog.Errorln(err)
	}

	if err := encoder.Flush(); err != nil {
		glog.Errorln(err)
	}

	req, err := http.NewRequest("POST", "http://192.168.2.4:8899/services/WorkflowService", buffer)
	if err != nil {
		glog.Errorln(err)
	}

	req = req.WithContext(context.Background())

	req.Header.Add("Content-Type", "text/xml; charset=\"utf-8\"")
	req.Header.Add("SOAPAction", "urn:weaver.workflow.webservices.WorkflowService.getUserid")
	req.Header.Set("User-Agent", "gowsdl/0.1")

	req.Close = true

	var hClient *http.Client
	if hClient == nil {
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{},
			DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
				d := net.Dialer{Timeout: 30 * time.Second}
				return d.DialContext(ctx, network, addr)
			},
			TLSHandshakeTimeout: 15 * time.Second,
		}
		hClient = &http.Client{Timeout: 90 * time.Second, Transport: tr}
	}

	res, err := hClient.Do(req)
	if err != nil {
		glog.Errorln(err)
	}
	defer res.Body.Close()
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		glog.Errorln(err)
	}
	glog.Infof("%+v", string(data))
}
