package main

import "github.com/spf13/viper"

var (
	ViperConfig *viper.Viper
)

// {"Id":"latch_legacy","Enabled":true,"TimeRange":[0,0],"UpperBound":0,"LowerBound":0}
// {"Id":"brightchanged_legacy","Enabled":true,"TimeRange":[0,0],"UpperBound":0,"LowerBound":0}

func main() {
	writeFile()
}

type AlertFilterData struct {
	Id        string `json:"Id,omitempty"`
	Enabled   bool   `json:"Enabled"`
	TimeRange [2]int
}

type Data struct {
	AlertFilterData
	UpperBound int
	LowerBound int
	EventType  string
}

func writeFile() {
	ViperConfig = viper.New()
	ViperConfig.AddConfigPath("/Users/hpu_scg/gocode/src/temp/GoModelTest/yamlTest/writeFile")
	ViperConfig.SetConfigName("default")
	ViperConfig.SetConfigType("yaml")
	if err := ViperConfig.ReadInConfig(); err != nil {
		panic(err)
	}
	data1 := []Data{{
		AlertFilterData: AlertFilterData{
			Id:        "latch_legacy",
			Enabled:   true,
			TimeRange: [2]int{0, 0},
		},
		UpperBound: 0,
		LowerBound: 0,
		EventType:  "latch_legacy",
	},
	}
	ViperConfig.Set("rule.latch_legacy", data1)
	data2 := []Data{
		{
			AlertFilterData: AlertFilterData{
				Id:        "brightchanged_legacy",
				Enabled:   true,
				TimeRange: [2]int{0, 0},
			},
			UpperBound: 0,
			LowerBound: 0,
			EventType:  "brightchanged_legacy",
		},
	}
	ViperConfig.Set("rule.brightchanged_legacy", data2)

	ViperConfig.Set("config.event_status_ttl", 10)
	ViperConfig.Set("config.persistence_cycle", 60)
	ViperConfig.Set("config.clear_cycle", 3600)
	ViperConfig.Set("config.tss_ttl_days", 30)
	ViperConfig.Set("config.reset_time", 10)
	ViperConfig.Set("config.post_video", false)

	ViperConfig.Set("eventserver.username", "haomuT")
	ViperConfig.Set("eventserver.password", "abc@Dgsh")
	ViperConfig.Set("eventserver.pub_http_url", "http://IP:8888")
	if err := ViperConfig.WriteConfig(); err != nil {
		panic(err)
	}

}
