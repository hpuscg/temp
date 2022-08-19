package main

type DeviceInfo struct {
	Ip              string
	Name            string
	IotInfo         Iot
	EventServerInfo EventServer
	NtpInfo         Ntp
}

type Iot struct {
	UserName string `json:"user_name"`
	PassWord string `json:"pass_word"`
	Server   string `json:"server"`
	Topic    string `json:"topic"`
}

type EventServer struct {
	UserName string `json:"user_name"`
	PassWord string `json:"pass_word"`
	Url      string `json:"url"`
}

type Ntp struct {
	Mode   int    `json:"Mode"`
	Server string `json:"Server"`
	Ntp    string `json:"Ntp"`
	Date   string `json:"Date"`
}

const (
	BumbleToken = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjQxMTQ4MzIzODUsImlzcyI6ImRlZXBnbGludCIsIlVzZXJ" +
		"JRCI6ImJ1bWJsZSJ9.VWGZm5LkQDoyukekwg6KEG-BbAkP28lcpx8D32t5mLw"
)
