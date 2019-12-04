package main

import (
	"net/smtp"
	"strings"
	"fmt"
)

const (
	HOST       = "smtp.163.com"
	ServerAddr = "smtp.163.com:25"
	USER       = "hpu_scg@163.com"
	PASSWORD   = "hpuscg123"
)

type Email struct {
	to      string
	subject string
	msg     string
}

func NewEmail(to, subject, msg string) *Email {
	return &Email{to: to, subject: subject, msg: msg}
}

func SendEmail(email *Email) error {
	auth := smtp.PlainAuth("", USER, PASSWORD, HOST)
	sendTo := strings.Split(email.to, ";")
	done := make(chan error, 1024)

	go func() {
		defer close(done)
		for _, v := range sendTo {
			str := strings.Replace("From: " + "aManLoveYou@163.com" + "~To: " + v + "~Subject: " + email.subject + "~~", "~", "\r\n", -1) + email.msg
			err := smtp.SendMail(
				ServerAddr,
				auth,
				USER,
				[]string{v},
				[]byte(str),
			)
			done <- err
		}
	}()
	for i := 0; i < len(sendTo); i++ {
		<-done
	}
	return nil
}



func main() {
	myContent := "Dear Li:\r\n	I Love You Forever!"
	email := NewEmail("m18715179281@163.com", "六一儿童节快乐！", myContent)

	err := SendEmail(email)
	if err != nil {
		fmt.Println(err)
	}else{
		fmt.Println("send email success")
	}
}
