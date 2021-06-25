package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"temp/GoModelTest/cgo/lic"
	"time"
)

func main() {
	uuid := "00000000000000006333303236313040"
	template := "{\"function\": \"人脸抓拍\",\"lic_type\": \"app_auth\",\"others\": \"\",\"video_chan\": 4}"
	sign := "123"
	fmt.Println(GenLicense(uuid, "", 1, "", template, sign))
}

var (
	// generate temporary files in WorkPath
	WorkPath = ""
	// import license to ApplyPath / ApplyLicName
	ApplyPath = ""
)

const (
	SRC_FILE_PATH     = "/product_config.json"
	LICENCE_FILE_PATH = "/fg.licence"
	PRI_KEY_FILE_PATH = "/pri_key.pem"
	PUB_KEY_FILE_PATH = "/pub_key.pem"
	ApplyLicName      = "fg.applied"
	PriKeySign        = "#_dg2021pri_#"
	LicVersion        = "V2"
)

type RetS struct {
	Code      int         `json:"code"`
	Message   string      `json:"message"`
	CheckSign string      `json:"check_sign"`
	Data      interface{} `json:"data"`
}

const (
	SUCCEED       int = 200
	FAILED        int = 201
	INTERNALERROR int = 5000
)

var MessageMap = map[int]string{
	SUCCEED:       "SUCCEED",
	FAILED:        "FAILED",
	INTERNALERROR: "INTERNALERROR",
}

//export GenRandomSign
// 1. generate random seed and return uint64 as random value
func GenRandomSign() string {
	rand.Seed(time.Now().UnixNano())
	return fmt.Sprintf("%d", rand.Uint64())
}

//export sign_pass
// 2. sign_pass: check if input inRandomSign matches inSignAck
func sign_pass(inRandomSign, inSignAck string) bool {
	return SignAckCheck(inRandomSign, inSignAck)
}

func GenMD5Str(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

// input a random sed, return a hash sed, to avoid api fake
func SignAck(inRandomSign string) string {
	return GenMD5Str(inRandomSign + PriKeySign)
}

// check if api reutrn inChallengeResult is match inSed
func SignAckCheck(inRandomSign, inSignAck string) bool {
	return inSignAck == SignAck(inRandomSign)
}

func GenRetChar(code int, checksign string, data interface{}) string {
	return GenRet(code, checksign, data)
}

func GenRet(code int, checksign string, data interface{}) string {
	ret := RetS{}
	ret.Code = code
	ret.CheckSign = checksign
	if messge, ok := MessageMap[code]; ok {
		ret.Message = messge
	}
	ret.Data = data
	retJ, err := json.Marshal(ret)
	if err != nil {
		return err.Error()
	}
	return string(retJ)
}

// =====amd64
// 1. GenLicense
// 2. InfoOrigin

// =====arm64
// 3. GetUUID_hi3559 + GetUUID_x86 + GetUUID_safenet + GetUUID_encryptchip
// 4. ImportLicense
// 5. InfoAppled
// 6. InfoOrigin
// 7. EncryptFile
// 8. DecryptFile

func init() {
	licWorkPathEnv := os.Getenv("LICWORKPATH")
	licApplyPathEnv := os.Getenv("LICAPPLYPATH")
	if licWorkPathEnv != "" {
		WorkPath = licWorkPathEnv + "/.licworkpath"
	} else {
		WorkPath, _ = filepath.Abs(filepath.Dir(os.Args[0]))
		WorkPath = WorkPath + "/.licworkpath"
	}
	if licApplyPathEnv != "" {
		ApplyPath = licApplyPathEnv
	} else {
		ApplyPath = WorkPath
	}
	if _, err := os.Lstat(WorkPath); err != nil && os.IsNotExist(err) {
		os.MkdirAll(WorkPath, os.ModePerm)
	}
	if _, err := os.Lstat(ApplyPath); err != nil && os.IsNotExist(err) {
		os.MkdirAll(ApplyPath, os.ModePerm)
	}
	fmt.Println("workpath = ", WorkPath)
	fmt.Println("apply path = ", ApplyPath)
}

//export GenLicense
// 3. 生成证书
// 输入: uuid, lic_path, auth_days_total, check_method, template, sign
// 输出: code, message, checksign, data
func GenLicense(uuid, lic_path string, auth_days_total int64, check_method string, template, sign string) string {
	// 0. gen sign_ack
	retSignAck := ""
	if sign == "" {
		retSignAck = ""
	} else {
		retSignAck = SignAck(sign)
	}

	// 1.update product_config.json
	origin := lic.LicGenInfo{}
	origin.UUID = uuid
	origin.Version = LicVersion
	origin.LicSrc = make(map[string]interface{})
	if err := json.Unmarshal([]byte(template), &origin.LicSrc); err != nil {
		return err.Error()
	}
	// 1. update AuthDateGen + AuthDaysTotal + AuthDateDeadline
	origin.AuthDateGen = time.Now().Format("2006-01-02")
	origin.AuthDaysTotal = auth_days_total
	if origin.AuthDaysTotal == lic.AuthDaysNoLimit {
		origin.AuthDateDeadline = fmt.Sprintf("%d", lic.AuthDaysNoLimit)
	} else {
		origin.AuthDateDeadline = time.Now().AddDate(0, 0, int(auth_days_total)).Format("2006-01-02")
	}
	// 2. update CheckMethod
	default_check_method := "countdown"
	check_method_go := check_method
	if len(check_method_go) == 0 || check_method_go == "countdown" {
		origin.CheckMethod = default_check_method
	} else if check_method_go == "deadline" {
		origin.CheckMethod = check_method_go
	} else {
		return GenRetChar(FAILED, retSignAck, fmt.Sprintf("check_method illegal:%s, must be deadline or countdown!", check_method_go))
	}
	// 3. write back
	if jb, err := json.MarshalIndent(origin, "", "	"); err != nil {
		return GenRetChar(INTERNALERROR, retSignAck, err.Error())
	} else {
		fd, err := os.OpenFile(WorkPath+SRC_FILE_PATH, os.O_RDWR|os.O_TRUNC|os.O_CREATE, os.ModePerm)
		if err != nil {
			return GenRetChar(INTERNALERROR, retSignAck, err.Error())
		}
		if _, err := fd.Write(jb); err != nil {
			return GenRetChar(INTERNALERROR, retSignAck, err.Error())
		}
		fd.Close()
	}

	// 4.generate license
	lic_path_go := lic_path
	if lic_path_go == "" {
		lic_path_go = WorkPath + LICENCE_FILE_PATH
	}
	authOne, err := lic.InitAuth(WorkPath+SRC_FILE_PATH, lic_path_go, WorkPath+PRI_KEY_FILE_PATH, WorkPath+PUB_KEY_FILE_PATH)
	if err != nil {
		return GenRetChar(INTERNALERROR, retSignAck, fmt.Sprintf("ERROR: InitAuth-> %s", err.Error()))
	}
	if err := authOne.GenLicenceInfo(uuid); err != nil {
		return GenRetChar(INTERNALERROR, retSignAck, fmt.Sprintf("ERROR: GenLicenceInfo-> %s", err.Error()))
	}
	return GenRetChar(SUCCEED, retSignAck, authOne.GetLicFileInfo())
}

//export InfoOrigin
// 4. 查看原始证书
// 输入: lic_path, sign
// 输出: code, message, checksign, data
func InfoOrigin(lic_path, sign string) string {
	// 0. gen sign_ack
	retSignAck := ""
	if sign == "" {
		retSignAck = ""
	} else {
		retSignAck = SignAck(sign)
	}

	lic_path_go := lic_path
	if lic_path_go == "" {
		lic_path_go = WorkPath + LICENCE_FILE_PATH
	}

	authOne, err := lic.InitAuth("", lic_path_go, "", "")
	if err != nil {
		return GenRetChar(FAILED, retSignAck, fmt.Sprintf("ERROR: InitAuth-> %s", err.Error()))
	}
	// 1. decode lic using V2 license src format
	srcInfo, err := authOne.DecLicV2(authOne.GetLicUUID())
	if err != nil {
		return GenRetChar(FAILED, retSignAck, fmt.Sprintf("ERROR: DecLicV2-> %s", err.Error()))
	}
	//2. check if lic format is V1
	if srcInfo.Version == "" {
		srcinfo_map, err := authOne.DecLic(authOne.GetLicUUID())
		if err != nil {
			return GenRetChar(FAILED, retSignAck, fmt.Sprintf("ERROR: DecLic-> %s", err.Error()))
		}
		return GenRetChar(SUCCEED, retSignAck, srcinfo_map)
	}
	return GenRetChar(SUCCEED, retSignAck, srcInfo)
}

//export EncryptFile
func EncryptFile() string {
	return ""
}

//export DecryptFile
func DecryptFile() string {
	return ""
}
