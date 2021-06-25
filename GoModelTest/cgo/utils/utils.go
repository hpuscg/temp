package utils

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	log "github.com/xiaomi-tc/log15"
	"io/ioutil"
	"math"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"syscall"
	"temp/GoModelTest/cgo/alg"
	"time"
)

const SHELLEXECTIMEOUTDEFAULT = time.Minute * 5
const PriKeyChallenge = "#_dg2021pri_#"

func GetRuntimeInfo() string {
	_, file, line, _ := runtime.Caller(1)
	//funcName := runtime.FuncForPC(pc)
	if os.Getenv("ENV_PRODUCTION") == "true" {
		return ""
	} else {
		return "[" + filepath.Base(file) + ":" + strconv.Itoa(line) + "]"
	}
}

func ExecShell(s string) (string, error) {
	cmd := exec.Command("/bin/sh", "-c", s)
	out, err := cmd.CombinedOutput()
	return string(out), err
}

func ExecShellTimeout(s string, d time.Duration) (string, error) {
	type st_cmdOut struct {
		out string
		err error
	}

	resultChan := make(chan st_cmdOut, 1)
	go func(in chan st_cmdOut) {
		cmd := exec.Command("/bin/sh", "-c", s)
		out, err := cmd.CombinedOutput()
		in <- st_cmdOut{
			out: string(out),
			err: err,
		}
	}(resultChan)

	select {
	case <-time.After(d):
		return "", errors.New("ExecShellTimeout: " + d.String() + " | " + "cmd: " + s)
	case res := <-resultChan:
		return res.out, res.err
	}
}

func ExecShellTimeoutDefault(s string) (string, error) {
	return ExecShellTimeout(s, SHELLEXECTIMEOUTDEFAULT)
}

func GenMD5(path string) (string, error) {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}
	md5Ctx := md5.New()
	md5Ctx.Write(file)
	cipherStr := md5Ctx.Sum(nil)
	return hex.EncodeToString(cipherStr), nil
}

func GenMD5Str(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

// 获取正在运行的函数名
func RunFuncName() string {
	pc := make([]uintptr, 1)
	runtime.Callers(2, pc)
	f := runtime.FuncForPC(pc[0])
	return f.Name()
}

func GetPKGNames() []string {
	cmd := fmt.Sprintf("ls /home/deepglint/PKG")
	outstr, err := ExecShellTimeoutDefault(cmd)
	if err != nil {
		return []string{}
	}
	outstr = strings.Trim(outstr, " \n")
	outList := strings.Split(outstr, "\n")
	outListClean := make([]string, 0)
	for _, val := range outList {
		tmpVal := strings.Trim(val, " ")
		if tmpVal != "" {
			outListClean = append(outListClean, tmpVal)
		}
	}
	return outListClean
}

func HttpGetFile(inURL string) (*http.Response, error) {
	req, err := http.NewRequest("GET", inURL, nil)
	if err != nil {
		return nil, err
	}
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if strings.Trim(res.Header.Get("Content-Disposition"), " ") == "" {
		rb, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}
		return nil, errors.New(string(rb))
	}
	return res, nil
}

//生成秘钥
func GetDeviceUUID() string {
	cmd := "fgdevinfo dev-uuid"
	cmdOut, err := ExecShellTimeoutDefault(cmd)
	if err != nil {
		cmdOut = ""
	}
	return cmdOut
}

func Gen32ByteAesKey(inKey string) string {
	h := md5.New()
	h.Write([]byte("deepglint" + inKey + "123_%&~*^"))
	return hex.EncodeToString(h.Sum(nil))
}

func Sync() {
	cmdstr := "sync"
	log.Debug("cmd", "debug", fmt.Sprintf("%s -> %s", cmdstr, "start"))
	outstr, err := ExecShell(cmdstr)
	if err != nil {
		log.Error("cmd", "error", fmt.Sprintf("%s -> %s", cmdstr, outstr+err.Error()))
	}
}
func FSync(inPath string) {
	cmdstr := fmt.Sprintf("fsync %s", inPath)
	log.Debug("cmd", "debug", fmt.Sprintf("%s -> %s", cmdstr, "start"))
	outstr, err := ExecShell(cmdstr)
	if err != nil {
		log.Error("cmd", "error", fmt.Sprintf("%s -> %s", cmdstr, outstr+err.Error()))
	}
}

// -------------------------------------------------------------for api challenge to avoid api fake -----

// 1. generate random seed and return uint64 as random value
func GenRandomSign() string {
	rand.Seed(time.Now().UnixNano())
	return fmt.Sprintf("%d", rand.Uint64())
}

// input a random sed, return a hash sed, to avoid api fake
func SignAck(inSed, inPriKey string) string {
	return GenMD5Str(inSed + inPriKey)
}

// check if api reutrn inChallengeResult is match inSed
func SignAckCheck(inRandomSign, inSignAck string) bool {
	return inSignAck == SignAck(inRandomSign, PriKeyChallenge)
}

// 2. generate seed by timestamp
func GenSign() (string, error) {
	tsU := time.Now().Unix()
	origin := fmt.Sprintf("%s%d%s", PriKeyChallenge, tsU, PriKeyChallenge)
	encrypted, err := alg.AesEncryptCBC([]byte(origin), []byte(Gen32ByteAesKey(PriKeyChallenge)))
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(encrypted), nil
}

func CheckSign(inSign string) bool {
	tsTimeout := 2
	tsNow := time.Now().Unix()
	signEncrypted, err := base64.StdEncoding.DecodeString(inSign)
	signDecrypted, err := alg.AesDecryptCBC(signEncrypted, []byte(Gen32ByteAesKey(PriKeyChallenge)))
	if err != nil {
		return false
	}
	pkLen := len(PriKeyChallenge)
	totalLen := len(signDecrypted)
	if (pkLen * 2) >= totalLen {
		return false
	}

	inTs := signDecrypted[pkLen:(totalLen - pkLen)]
	inSignTsNum, err := strconv.ParseInt(string(inTs), 10, 64)
	if err != nil {
		return false
	}

	if math.Abs(float64(tsNow-inSignTsNum)) <= float64(tsTimeout) {
		return true
	}
	return false
}

func CheckPid(pid int) bool {
	out, _ := exec.Command("kill", "-s", "0", strconv.Itoa(pid)).CombinedOutput()
	if string(out) == "" {
		return true // pid exist
	}
	return false
}

var globalFileLock *os.File

func FileLock(filePath string) error {
	var err error
	globalFileLock, err = os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return err
	}
	return syscall.Flock(int(globalFileLock.Fd()), syscall.LOCK_EX|syscall.LOCK_NB)
}
