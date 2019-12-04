/*
#Time      :  2019/1/10 下午5:35 
#Author    :  chuangangshen@deepglint.com
#File      :  chmodTest.go
#Software  :  GoLand
*/
package main

import (
	"net"
	"golang.org/x/crypto/ssh"
	"time"
	"fmt"
)

func main() {
	Ips := []string{
		"10.147.15.172", "10.147.16.228", "10.147.17.91", "10.147.17.181", "10.147.13.211", "10.147.13.218",
		"10.147.9.150", "10.147.12.214",  "10.147.19.200", "10.147.11.17", "10.147.8.153", "10.147.8.155",
		"10.147.16.8", "10.147.14.233", "10.147.14.239", "10.147.8.51", "10.147.8.52",
		"10.147.17.3", "10.147.17.180", "10.147.12.215", "10.147.16.186", "10.147.9.90", "10.147.15.118",
		"10.147.15.119", "10.147.17.217",
	}
	DefaultIps := []string{
		"10.147.17.3", "10.147.17.180", "10.147.12.215", "10.147.16.186", "10.147.9.90", "10.147.15.118",
		"10.147.15.119", "10.147.17.217",
	}
	fmt.Println(DefaultIps)
	/*Ips := []string{
		"192.168.100.223",
	}*/
	fmt.Println("=======")
	for _, Ip := range Ips {
		go handlers(Ip)
	}
	time.Sleep(300 * time.Second)
}

func handlers(Ip string) {
	// try get ssh client
	client, err := connect("root", Ip, 22)
	if err != nil {
		info := fmt.Sprintf("connect to %s err : %s\n", Ip, err)
		fmt.Println(info)
		return
	}
	defer client.Close()
	var session *ssh.Session
	// reset config
	session, _ = client.NewSession()
	session.Run(`chmod +x /usr/bin/reload_sensor.sh`)
	session, _ = client.NewSession()
	session.Run(`reload_sensor.sh`)
	fmt.Printf("%s is ok!\n", Ip)

}



// connect for exec cmd
func connect(user, host string, port int) (*ssh.Client, error) {
	var (
		auth 			[]ssh.AuthMethod
		addr 			string
		clientConfig 	*ssh.ClientConfig
		client 			*ssh.Client
		err 			error
	)
	auth = make([]ssh.AuthMethod, 0)
	testPrivateKeys, err := ssh.ParseRawPrivateKey([]byte(`-----BEGIN RSA PRIVATE KEY-----
MIIEowIBAAKCAQEApnIg4Q/g2thAR3vAUw6EPjqgWIEJ7+FZ+AQZtHUc7b920VJI
7JPmZ1xwdUArlCpQIMAt6mAwV9Z/C+Nc9qIpIeQwKaAd6YWVdr3jFyHLC9rqIT2g
VifCDnMkSnV7Lvuu5QTvgURGOYpyUhUDJBdBY4YAu9q1ITy35oB0xLh1vUCwuDxI
oM5lMc+HsPjf4/SyfyAacBuoD7BvAJsxJ6xuXBaIlmWcw8o76O/Y5PGcYKPS9/bI
rN8TrstuWILp2Nvi4WoxVMIQ98i1S6jM47arI+vNGlFrwolrCanH8GBj1NOBh4BF
JwJisi0Z3+RrtxOVRtgZ9S/tKdK73X6EpbN4hwIDAQABAoIBAAuBRAiKgm5eGENY
qHiVPkrW3pJ/iOJN31wnXGd+2NsOKvZZC7Vem8R1PUi9gMWjDxrUbdgPggfwSaPW
uWxK1TEEhte5u5eSpjwo7/N/YHuXTCu0CMsrwFwjVVTYPgWHXBV0e+GhiIEdsr09
upPaD6kDcDWL7o03lzaVlnyqi2jjXT6kUDyEFCbIAGtoxaYf3clT5e30FnyZhiCH
m8/Qqv5M1wcVIVdsItHqMsQXQF34eT/Lg3r/Ui1bQcUldc6yYjGpC08EdDNKhGT2
f2QwAv7UJ+GB8RNl12w3fAh3ReuiW8NEtDQ1nuSahkX5YlIWkqRDOd6Sjrg1ZkfW
u0/zPZECgYEA2m+w90vb3ui7M/Q0AYJivo88YKhT3ismQs2+CkkgWJ7IohJj3VSh
REljeAwEVEKv8G8lXgjTNKQ+B4sPFckIvIWGkwo7cuerIwn9n41K20oGb6gEl0jW
mVbhv0dy6yfp8deBCOZB4YgonXWsuv4lw8DaUoakGxZgFfChjH0VvbUCgYEAwxGj
rmq+RQWYYna9WWn2GPEJoX0SBU39pHQYBKfQ++pMIUrrryCjPvBNnIICng82RjTp
MU8BvudvDCJgj3J79TDetBnwVt8/nAGIkleyuWzDMQwF7khBS9/TqUUqmH88GmOt
40BPThCBx8YgKiPpmGYgPnUww1bqpvxKT9O0IssCgYEAjFH7qKD+mW9/8pwJXH7Z
1/hDnQQE/E9TwM5SKmFXehZmZFbT+DaJckiCsXdmwIomY5nCs2mP490uS8I06pW+
GvzbulF0ZxgTg+rDFl+5mq0u/UM9z8FmuhJp6mqHlDCLxGPf7EuePrctABm74FOr
Btk4ZpM/kHcLOozd+lXQRZECgYBipWr26zgpQ3kaYh3DN9iiKFLMfak9UYFxRtxW
jl8a5hN1yqOBPqoPTAqTmROlxt+VhXBf5Spm1jbMFh5qrGSPTBVzUqK968wJIqVk
DEFvj9bt2LyvEY8jxZ8OPNIbqExGtB3djEoOmj5nPoRJizu4O/0WWME+J5gmtfMG
h3LTHQKBgDlITGqdIM4Pp54X5ppOW9S55yaAMBJUUhgUsJ73vEcQsBCZ8xkJXg/Q
muPfcFzSD/IgeFoWxYrJIk0CBov3ah+14z5YV1JoKIXAlL7V18f7Omaav8/bozOP
x78MQ06CGEFRcD4LPMITxTDj6zDm1h7iPhG4m2c9Shy0rwpFmFdd
-----END RSA PRIVATE KEY-----`))
	if err != nil {
		// fmt.Println("Unable to parse test key :", err)
		return nil, err
	}
	testSingers, _ := ssh.NewSignerFromKey(testPrivateKeys)

	auth = append(auth, ssh.PublicKeys(testSingers))
	clientConfig = &ssh.ClientConfig{
		User: 				user,
		Auth: 				auth,
		Timeout: 			30 * time.Second,
		HostKeyCallback: 	func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	}
	addr = fmt.Sprintf("%s:%d", host, port)
	if client, err = ssh.Dial("tcp", addr, clientConfig); err != nil {
		return nil, err
	}
	return client, nil
}
