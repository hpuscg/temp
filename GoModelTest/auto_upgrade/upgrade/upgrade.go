package main

import (
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

type Server struct {
	Ip       string
	Port     int
	Username string
}

var CurServer Server

/*
退出状态代码：
0 	命令成功完成
1	通常的未知错误
2	误用shell命令
126	命令无法执行
127	没有找到命令
128	无效的退出参数
*/
/*
func main() {
	ip := flag.String("ip", "192.168.100.234", "ip of sensor")
	file := flag.String("u", "", "source file at local")
	cmd := flag.String("c", "", "run cmd at remote , after upload file")
	remotefile := flag.String("r", "/home/ubuntu/", "dest path at remote")
	e := flag.Bool("e", false, "show error code")
	s := flag.Bool("s", false, "show sample command")
	//	isofile := flag.String("s", "", "ISO filename of upgrade package")

	flag.Parse()

	if *e {
		fmt.Println("Remote run shell common error code :")
		fmt.Println("0 	命令成功完成")
		fmt.Println("1	通常的未知错误")
		fmt.Println("2	误用shell命令")
		fmt.Println("126	命令无法执行")
		fmt.Println("127	没有找到命令")
		fmt.Println("128	无效的退出参数")
		return
	}

	if *s {
		fmt.Println("Local run shell sample :")
		fmt.Println(`./upgradetk1 -ip 192.168.100.234 -u active_libra_init.sh -c "bash /home/ubuntu/active_libra_init.sh"`)
		return
	}
	err := tryPing(*ip)
	if err != nil {
		fmt.Printf("can't ping %s for %s,please checkout your config ip\n", *ip, err)
		os.Exit(1)
		return
	}

	CurServer.Ip = *ip
	CurServer.Port = 22
	CurServer.Username = "root"
	///////////////////////////////////////////////////////////////////////////////
	if len(*file) > 0 {
		f, e := os.Stat(*file)
		if e != nil || f.IsDir() {
			fmt.Println("please enter a real file")
			return
		}
	}
	if len(*file) == 0 && len(*cmd) == 0 {
		fmt.Println("please enter a real file or a command")
		return
	}
	///////////////////////////////////////////////////////////////////////////////
	remote := *remotefile
	if (remote)[len(remote)-1] != os.PathSeparator {
		remote += string(os.PathSeparator)
	}
	df := filepath.Base(*file)
	fmt.Println(remote)
	fmt.Println(df)
	var command string
	if len(*cmd) > 0 {
		command = *cmd
	} else if len(*file) > 0 {
		command = remote + df
	} else {
		fmt.Println("please input command ")
		return
	}
	//	command += " >/home/ubuntu/debug"
	fmt.Println("command is: ", command)
	///////////////////////////////////////////////////////////////////////////////
	fmt.Printf("connect to %s ...\n", *ip)

	///////////////////////////////////////////////////////////////////////////////
	copySession, err := copyConnect("root", *ip, 22)
	if err != nil {
		fmt.Println("connect sftp fail, ", err)
		return
	}
	///////////////////////////////////////////////////////////////////////////////

	// copyFiles
	if len(*file) > 0 {
		copyFile(copySession, remote, *file)

		runCommand("chmod a+x " + remote + df)
	}

	//	runCommand("cd " + remote)
	runCommand(command)

	if len(*file) > 0 {
		//		runCommand("rm -rf " + remote + df)
	}
	fmt.Println("session end")
}

*/

func main() {
	fmt.Println("2222")
	session, err := connect("root", "192.168.100.232", 22)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("1111")
	// err = session.Run("docker run -i -t -v /libra:/libra 192.168.5.46:5000/armhf-libra-init:1.1.0")
	err = session.Run("docker run -i -t -rm -v /libra:/libra 192.168.5.46:5000/armhf-libra-init:1.1.0")
	/*err = session.Run("docker run -tid --name eventserver --memory=800m --restart=always --net=host -v /tmp:/tmp " +
		"-v /data/tf/eventserver:/data/slice -v /data/eventserver:/data/event -v /etc/localtime:/etc/localtime:ro 192.168" +
		".5.46:5000/armhf-eventserver:1.7.7 ./eventserver.arm -etcdserver=http://127.0.0.1:4001 -mode=fat -report_interval=600 -log_dir_deepglint /tmp/ -mem=750")*/
	fmt.Println("33333")
	fmt.Println("99999", err)
}

// connect for exec cmd
func connect(user, host string, port int) (*ssh.Session, error) {
	var (
		auth         []ssh.AuthMethod
		addr         string
		clientConfig *ssh.ClientConfig
		client       *ssh.Client
		session      *ssh.Session
		err          error
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
		fmt.Println("Unable to parse test key :", err)
	}
	testSingers, _ := ssh.NewSignerFromKey(testPrivateKeys)

	auth = append(auth, ssh.PublicKeys(testSingers))
	clientConfig = &ssh.ClientConfig{
		User: user,
		Auth: auth,
		//		Timeout: 			60 * time.Second,
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	}
	addr = fmt.Sprintf("%s:%d", host, port)
	if client, err = ssh.Dial("tcp", addr, clientConfig); err != nil {
		return nil, err
	}
	if session, err = client.NewSession(); err != nil {
		return nil, err
	}
	return session, nil
}

// connect for copy file
func copyConnect(user, host string, port int) (*sftp.Client, error) {
	var (
		auth         []ssh.AuthMethod
		addr         string
		clientConfig *ssh.ClientConfig
		sshClient    *ssh.Client
		sftpClient   *sftp.Client
		err          error
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
		fmt.Println("Unable to parse test key :", err)
	}
	testSingers, _ := ssh.NewSignerFromKey(testPrivateKeys)
	auth = append(auth, ssh.PublicKeys(testSingers))
	clientConfig = &ssh.ClientConfig{
		User: user,
		Auth: auth,
		//		Timeout: 			60 * time.Second,
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	}
	addr = fmt.Sprintf("%s:%d", host, port)
	if sshClient, err = ssh.Dial("tcp", addr, clientConfig); err != nil {
		return nil, err
	}
	if sftpClient, err = sftp.NewClient(sshClient); err != nil {
		return nil, err
	}
	return sftpClient, nil
}

// copy file and load docker image
func copyFile(copySession *sftp.Client, destpath string, file string) {
	srcFile, err := os.Open(file)
	if err != nil {
		info := fmt.Sprintf("open %s failed :%s", file, err)
		fmt.Println(info)
	}
	df := filepath.Base(file)
	dstFile, err := copySession.Create(destpath + df)
	if err != nil {
		info := fmt.Sprintf("create %s failed :%s", file, err)
		fmt.Println(info)
	}
	// defer dstFile.Close()
	ff, err := ioutil.ReadAll(srcFile)
	if err != nil {
		fmt.Println("readall err: ", err)
	}
	dstFile.Write(ff)
	fmt.Printf("copy file (%s) to (%s) successful\n", file, destpath+df)
}

// rm docker image which need upgrade
func runCommand(cmd string) error {
	session, err := connect(CurServer.Username, CurServer.Ip, CurServer.Port)
	if err != nil {
		info := fmt.Sprintf("connect to %s err : %s", CurServer.Ip, err)
		fmt.Println(info)
		return err
	}
	fmt.Println("root:~# ", cmd)
	err = session.Run(cmd)
	if err != nil {
		fmt.Printf("fail for (%s) \n", err)
		return err
	} else {
		fmt.Printf("run command [%s] successful\n", cmd)
	}
	session.Close()
	return nil
}

//////////////////////////////////
// list files under dirPath
func FilesUnder(dirPath string) ([]string, error) {
	_, err := os.Stat(dirPath)
	if !(err == nil || os.IsExist(err)) {
		return []string{}, nil
	}

	fs, err := ioutil.ReadDir(dirPath)
	if err != nil {
		return []string{}, err
	}

	sz := len(fs)
	if sz == 0 {
		return []string{}, nil
	}

	ret := []string{}
	for i := 0; i < sz; i++ {
		ret = append(ret, fs[i].Name())
	}
	return ret, nil
}

// test the ip is right
func tryPing(IP string) error {
	var cmd *exec.Cmd
	sysInfo := runtime.GOOS
	if "windows" == sysInfo {
		cmd = exec.Command("ping", "-n", "1", IP)
	} else {
		cmd = exec.Command("ping", "-c", "1", IP)
	}
	_, err := cmd.Output()
	return err
}
