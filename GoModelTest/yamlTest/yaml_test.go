package yaml

import (
	"testing"
	"fmt"
	"reflect"
	"time"
	"encoding/json"
)

var yamlCli *YamlConfig = nil

func TestYamlConfig_GetValue(t *testing.T) {
	initYamlCli()
	/*ret, err := yamlcli.GetValue("/hotspot/door")
	if err != nil {
		fmt.Println("get value err:", err)
	}
	data := interface2string(ret)
	fmt.Println(data)
	ret, _ = json.Marshal(data)

	fmt.Println("get vuale result is : ", ret)
	// fmt.Println("marshal result is : ", marter)
	err = yamlcli.DeleteValue("/hotspot/door/door_1")
	if err != nil {
		fmt.Println("err is: ", err)
	} else {
		fmt.Println("no err !")
	}*/
	/*data := make(map[string]interface{})
	data["dwellingtime_legacy"] = {"Id": "dwellingtime_legacy","Enabled": true, "TimeRange": [0,0],"UpperBound": 1800,"LowerBound": 0}
	data = {"dwellingtime_legacy": {"Id": "dwellingtime_legacy","Enabled": true, "TimeRange": [0,0],"UpperBound": 1800,"LowerBound": 0},
		"population_legacy": {"Id": "population_legacy","Enabled": true,"TimeRange": [0,0],"UpperBound": 40,"LowerBound": 0},
		"velocity_legacy": {"Id": "velocity_legacy","Enabled": true,"TimeRange": [0,0],"UpperBound": 2200,"LowerBound": 0},
		"approaching_legacy": {"Id": "approaching_legacy","Enabled": true,"TimeRange": [0,0],"UpperBound": 0.9,"LowerBound": 0},
		"distance_legacy": {"Id": "distance_legacy","Enabled": true,"TimeRange": [0,0],"UpperBound": 100000,"LowerBound": 0},
		"singlepeoplein_legacy": {"Enabled": false,"Id": "singlepeoplein_legacy","LowerBound": 0,"TimeRange": [0,0],"UpperBound": 10,"WeekdayRange": 0}}
	setHotSpot(data, yamlcli)*/
	setYamlValue()
	deleteYamlValue()
	// yamlDelTest(yamlcli)
	// writeStr(yamlcli)
	// readStr(yamlcli)
	// writeUpper(yamlCli)
	ret, err := yamlCli.GetValue("/config/test")
	if ret == nil {
		fmt.Println(1111)
	}
	fmt.Println(ret, err)

}

func setHotSpot(data map[string]interface{}, yamlcli *YamlConfig)  {
	for key, value := range data {
		_, err := yamlcli.SetValue("/hotspot/door/" + key, value)
		fmt.Printf("key is: %s, value is: %s", key, value)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func initYamlCli() {
	var err error
	yamlCli, err = NewYamlConfig("data.yaml")
	if err != nil {
		fmt.Println("client to yaml err: ", err)
	}
}

func setYamlValue() {
	// rs, _ := json.Marshal("hello world")
	yamlCli.SetValue("/config/test", nil)
}

func deleteYamlValue() {
	yamlCli.DeleteValue("/config/test")
	readYaml()
}

func interface2string(data interface{}) map[string]interface{} {
	ret := make(map[string]interface{})
	switch vv := data.(type) {
	case map[interface{}]interface{}:
		for key, value := range vv {
			switch kk := value.(type) {
			case map[interface{}]interface{}:
				result := interface2string(kk)
				ret[key.(string)] = result
			default:
				ret[key.(string)] = value
			}
		}
	default:
		fmt.Println("=====vv is:====", vv)
		fmt.Println("=====vv type is:===", reflect.TypeOf(vv))
	}
	return ret
}

func forTest(yamlcli *YamlConfig)  {
	for {
		flg := true
		if flg {
			yamlcli.SetValue("data", "yes")
			time.Sleep(60 * time.Second)
			flg = false
		} else {
			flg = true
		}

	}
}

func readYaml()  {
	rult, _ := yamlCli.GetValue("/config/test")
	switch vv := rult.(type) {
	case interface{}:
		switch kk := vv.(type) {
		case []uint8:
			var s string
			json.Unmarshal(kk, &s)
			fmt.Println("s", s)
		case []interface{}:
			ret := make([]uint8, 0)
			for _, value := range kk {
				switch value.(type) {
				case interface{}:
					ret = append(ret, uint8(value.(int)))
				}
			}
			var s string
			json.Unmarshal(ret, &s)
			fmt.Println("ss ", s)
		}
	}
}

func yamlDelTest(yamlcli *YamlConfig)  {
	err := yamlcli.DeleteValue("/config/test")
	fmt.Println(err)
}

type data1 struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func writeStr(yamlcli *YamlConfig) {
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
	yamlcli.SetValue("/config/test2", relStr)
}

func readStr(yamlcli *YamlConfig) {
	var data2 data1
	ret, err := yamlcli.GetValue("/config/test2")
	if err != nil {
		fmt.Println(err)
	}
	switch vv := ret.(type) {
	case interface{}:
		switch kk := vv.(type) {
		case string:
			json.Unmarshal([]byte(kk), &data2)
			fmt.Println("111: ", data2)
		}
	}
}


/*'{"Id":"singlepeoplein_legacy","Query":"(.Attribute.DeveloperId==\"DeepGlint\")\u0026\u0026(.Attribute.ProjectId==\"DeepGlint\")\u0026\u0026(.Attribute.QueryId==\"SinglePeopleIn\")\u0026\u0026((.Duration
\u003e 10)||(.Duration \u003c 0))\u0026\u0026(.Population == 1)","TimeRange":[0,0],"EventBaseType":1000,"UpperBound":10,"LowerBound":0}'
'{"Id":"singlepeoplein_legacy","Query":"(.Attribute.DeveloperId==\"DeepGlint\")\u0026\u0026(.Attribute.ProjectId==\"DeepGlint\")\u0026\u0026(.Attribute.QueryId==\"SinglePeopleIn\")\u0026\u0026((.Duration \u003e 10)||(.Duration \u003c 0))\u0026\u0026(.Population == 1)","TimeRange":[0,0],"EventBaseType":1000,"UpperBound":10,"LowerBound":0}'
*/

func writeUpper(yamlcli *YamlConfig)  {
	str := `{"Id":"singlepeoplein_legacy","Query":"(.Attribute.DeveloperId==\"DeepGlint\")\u0026\u0026(.Attribute.ProjectId==\"DeepGlint\")\u0026\u0026(.Attribute.QueryId==\"SinglePeopleIn\")\u0026\u0026((.Duration \u003e 10)||(.Duration \u003c 0))\u0026\u0026(.Population == 1)","TimeRange":[0,0],"EventBaseType":1000,"UpperBound":10,"LowerBound":0}`
	// str := `{"yes == no"}`
	yamlcli.SetValue("/config/test3", str)
}

