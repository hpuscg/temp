package main

import (
	"fmt"

	"github.com/spf13/viper"
)

func main() {
	var w worker = person{}
	fmt.Println(w)
	readConfig()
}

type worker interface {
	work()
}

type person struct {
	name string
	worker
}

func readConfig() {
	config := viper.New()
	config.AddConfigPath("./")
	config.SetConfigName("fgConfig")
	config.SetConfigType("json")
	if err := config.ReadInConfig(); err != nil {
		fmt.Println(err)
	}
	config.SetDefault("deviceinfo.deviceName", "123")
	// fmt.Println(config.AllSettings())
	fmt.Println(config.Get("deviceinfo"))
	config.Set("deviceinfo.deviceName", "F57")
	if err := config.WriteConfig(); err != nil {
		fmt.Println(err)
	}
	fmt.Println(config.Get("deviceinfo"))
	fmt.Println(config.AllSettings())
}
