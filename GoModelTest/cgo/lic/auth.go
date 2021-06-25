package lic

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"temp/GoModelTest/cgo/alg"
	"temp/GoModelTest/cgo/utils"
	"time"
)

const (
	SHELLEXECTIMEOUTDEFAULT = time.Minute * 5
)

/*
主要解决的目标：保证授权文件的可靠性（非明文+不可篡改），主要有加密和签名两种方法。
1.1 加密，相应也要有解密还原功能（客户端生成）
可以约定一种对称加密算法如AES解决，秘钥uuid_key。形如auth.enc = aes.encode(auth.src, uuid_key)
1.2 签名（客户端生成）
可以使用非对称加密算法，如RSA事先用私钥（private_key)生成签名文文件（更安全），或者使用hash算法生成一个签名文件（不安全）
形如auth.sign = hash(auth.enc) 或者 auth.sign = rsa.encode(auth.enc, private_key)
1.3 生成授权文件导入设备（客户端）
FGUpgrade.TransferFile( auth.sign + auth.enc ) --> SysData
1.4 使用授权文件（app）
app启动后执行：验签+解密
验签：1.用公钥解密验签 rsa.decode(auth.enc, pub_key) ?= auth.sign（更安全）2. hash(auth.enc) ?= auth.sign（不安全）
解密：auth.conf = aes.decode(auth.enc, uuid_key)
读取：st_auth = json.unmarsh(auth.conf)
*/
type LicInfo struct {
	//files input and output
	src_file_path string
	lic_file_path string
	pub_key_path  string
	pri_key_path  string
	//generate info
	src_file_info LicGenInfo
	lic_file_info LicFileInfo

	//file read
	file_src_ctx     []byte
	file_lic_ctx     []byte
	file_pub_key_ctx []byte
	file_pri_key_ctx []byte
}

//encoded src info
type LicFileInfo struct {
	UUID   string `json:"uuid"`
	Lic    string `json:"enc_src"`
	PubKey string `json:"pub_key"`
	Sign   string `json:"sign"`
}

type LicGenInfo struct {
	Version          string                 `json:"version"`
	UUID             string                 `json:"uuid"`
	AuthDaysTotal    int64                  `json:"auth_days_total"`
	AuthDateGen      string                 `json:"auth_date_gen"`
	AuthDateDeadline string                 `json:"auth_date_deadline"`
	CheckMethod      string                 `json:"check_method"`
	LicSrc           map[string]interface{} `json:"lic_src"`
}

//init
func InitAuth(inSrcFilePath, inLicFilePath, inPriKeyPath, inPubKeyPath string) (*LicInfo, error) {
	out := &LicInfo{
		src_file_info: LicGenInfo{},
		lic_file_info: LicFileInfo{},
	}
	out.src_file_info.LicSrc = make(map[string]interface{})
	out.src_file_path = inSrcFilePath
	out.lic_file_path = inLicFilePath
	out.pri_key_path = inPriKeyPath
	out.pub_key_path = inPubKeyPath
	//load file, and parse source file info
	if _, err := os.Lstat(inSrcFilePath); err == nil {
		f, err := os.Open(inSrcFilePath)
		if err != nil {
			return nil, errors.New(utils.GetRuntimeInfo() + err.Error())
		}
		fb, err := ioutil.ReadAll(f)
		if err != nil {
			return nil, errors.New(utils.GetRuntimeInfo() + err.Error())
		}
		out.file_src_ctx = fb
		err = json.Unmarshal(fb, &out.src_file_info)
		if err != nil {
			return nil, errors.New(utils.GetRuntimeInfo() + err.Error())
		}
	}

	//load file
	if _, err := os.Lstat(inPubKeyPath); err == nil {
		f, err := os.Open(inPubKeyPath)
		if err != nil {
			return nil, errors.New(utils.GetRuntimeInfo() + err.Error())
		}
		fb, err := ioutil.ReadAll(f)
		if err != nil {
			return nil, errors.New(utils.GetRuntimeInfo() + err.Error())
		}
		out.file_pub_key_ctx = fb
	}
	//load file
	if _, err := os.Lstat(inPriKeyPath); err == nil {
		f, err := os.Open(inPriKeyPath)
		if err != nil {
			return nil, errors.New(utils.GetRuntimeInfo() + err.Error())
		}
		fb, err := ioutil.ReadAll(f)
		if err != nil {
			return nil, errors.New(utils.GetRuntimeInfo() + err.Error())
		}
		out.file_pri_key_ctx = fb
	}
	if _, err := os.Lstat(inLicFilePath); err == nil {
		f, err := os.Open(inLicFilePath)
		if err != nil {
			return nil, errors.New(utils.GetRuntimeInfo() + err.Error())
		}
		fb, err := ioutil.ReadAll(f)
		if err != nil {
			return nil, errors.New(utils.GetRuntimeInfo() + err.Error())
		}
		out.file_lic_ctx = fb
	}
	return out, nil
}

//1.
func (this *LicInfo) GenLicenceInfo(inDevUUID string) error {
	//conflict check
	if inDevUUID == "" {
		return errors.New(utils.GetRuntimeInfo() + "input params must all pass!!")
	}
	this.lic_file_info.UUID = inDevUUID

	if this.file_src_ctx == nil || len(this.file_src_ctx) == 0 {
		return errors.New(utils.GetRuntimeInfo() + "file_src_ctx empty!")
	}
	//1. gen lic
	encrypted, err := alg.AesEncryptCBC(this.file_src_ctx, []byte(utils.Gen32ByteAesKey(this.lic_file_info.UUID)))
	if err != nil {
		return errors.New(utils.GetRuntimeInfo() + err.Error())
	}
	this.lic_file_info.Lic = base64.StdEncoding.EncodeToString(encrypted)

	//2. gen pub_pri key if not exist
	if this.file_pri_key_ctx != nil && len(this.file_pri_key_ctx) != 0 {
		fmt.Println("pri_key already generated...")
	} else {
		err := alg.GenRsaKey(2048, this.pub_key_path, this.pri_key_path)
		if err != nil {
			return errors.New(utils.GetRuntimeInfo() + err.Error())
		}
		f, err := os.Open(this.pub_key_path)
		if err != nil {
			return errors.New(utils.GetRuntimeInfo() + err.Error())
		}
		fb, err := ioutil.ReadAll(f)
		if err != nil {
			return errors.New(utils.GetRuntimeInfo() + err.Error())
		}
		//update key_ctx
		this.file_pub_key_ctx = fb

		f, err = os.Open(this.pri_key_path)
		if err != nil {
			return errors.New(utils.GetRuntimeInfo() + err.Error())
		}
		fb, err = ioutil.ReadAll(f)
		if err != nil {
			return errors.New(utils.GetRuntimeInfo() + err.Error())
		}
		//update key_ctx
		this.file_pri_key_ctx = fb
	}
	this.lic_file_info.PubKey = string(this.file_pub_key_ctx)

	//3. gen sign
	if this.file_pri_key_ctx == nil || len(this.file_pri_key_ctx) == 0 {
		return errors.New(utils.GetRuntimeInfo() + "file_pri_key_ctx empty!")
	}
	sign, err := alg.Sign([]byte(this.lic_file_info.Lic), this.file_pri_key_ctx)
	if err != nil {
		return errors.New(utils.GetRuntimeInfo() + err.Error())
	}
	this.lic_file_info.Sign = base64.StdEncoding.EncodeToString(sign)
	//4. gen lic file
	fd, err := os.Create(this.lic_file_path)
	if err != nil {
		return errors.New(utils.GetRuntimeInfo() + err.Error())
	}
	jb, err := json.MarshalIndent(this.lic_file_info, "", "  ")
	if err != nil {
		return errors.New(utils.GetRuntimeInfo() + err.Error())
	}
	_, err = fd.Write(jb)
	if err != nil {
		return errors.New(utils.GetRuntimeInfo() + err.Error())
	}
	this.file_lic_ctx = jb
	return nil
}

//3. rsa(sign.rsa, sign.rsa.pub_key) -> cksum(enc) ?= cksum(enc)
func (this *LicInfo) checkSign() error {
	if this.file_lic_ctx == nil || len(this.file_lic_ctx) == 0 {
		return errors.New(utils.GetRuntimeInfo() + "file_lic_ctx empty!")
	}
	//unmashal sign file
	err := json.Unmarshal(this.file_lic_ctx, &this.lic_file_info)
	if err != nil {
		return errors.New(utils.GetRuntimeInfo() + err.Error())
	}
	originSign, err := base64.StdEncoding.DecodeString(this.lic_file_info.Sign)
	if err != nil {
		return errors.New(utils.GetRuntimeInfo() + err.Error())
	}
	err = alg.SignVer([]byte(this.lic_file_info.Lic), originSign, []byte(this.lic_file_info.PubKey))
	if err != nil {
		return errors.New(utils.GetRuntimeInfo() + err.Error())
	}
	return nil
}

func (this *LicInfo) GetLicUUID() string {
	if this.file_lic_ctx == nil {
		return ""
	}
	err := json.Unmarshal(this.file_lic_ctx, &this.lic_file_info)
	if err != nil {
		return ""
	}
	return this.lic_file_info.UUID
}

//4. enc + uuid_key -> src.uuid ?= dev.uuid
func (this *LicInfo) DecLic(inUUIDKey string) (map[string]interface{}, error) {
	if err := this.checkSign(); err != nil {
		return nil, errors.New("checkSign error:" + err.Error())
	}
	if this.file_lic_ctx == nil || len(this.file_lic_ctx) == 0 {
		return nil, errors.New(utils.GetRuntimeInfo() + "file_lic_ctx empty!")
	}
	err := json.Unmarshal(this.file_lic_ctx, &this.lic_file_info)
	if err != nil {
		return nil, errors.New(utils.GetRuntimeInfo() + err.Error())
	}
	if this.lic_file_info.UUID != inUUIDKey && strings.Trim(inUUIDKey, " ") != "" {
		return nil, errors.New(inUUIDKey + " not support this device!")
	}
	originLic, err := base64.StdEncoding.DecodeString(this.lic_file_info.Lic)
	if err != nil {
		return nil, errors.New(utils.GetRuntimeInfo() + err.Error())
	}
	file_src_ctx, err := alg.AesDecryptCBC(originLic, []byte(utils.Gen32ByteAesKey(inUUIDKey)))
	if err != nil {
		return nil, errors.New(utils.GetRuntimeInfo() + err.Error())
	}
	out := make(map[string]interface{})
	err = json.Unmarshal(file_src_ctx, &out)
	if err != nil {
		return nil, errors.New(utils.GetRuntimeInfo() + err.Error())
	}
	return out, nil
}

func (this *LicInfo) DecLicV2(inUUIDKey string) (*LicGenInfo, error) {
	if err := this.checkSign(); err != nil {
		return nil, errors.New("checkSign error:" + err.Error())
	}
	if this.file_lic_ctx == nil || len(this.file_lic_ctx) == 0 {
		return nil, errors.New(utils.GetRuntimeInfo() + "file_lic_ctx empty!")
	}
	err := json.Unmarshal(this.file_lic_ctx, &this.lic_file_info)
	if err != nil {
		return nil, errors.New(utils.GetRuntimeInfo() + err.Error())
	}
	if this.lic_file_info.UUID != inUUIDKey && strings.Trim(inUUIDKey, " ") != "" {
		return nil, errors.New(inUUIDKey + " not support this device!")
	}
	originLic, err := base64.StdEncoding.DecodeString(this.lic_file_info.Lic)
	if err != nil {
		return nil, errors.New(utils.GetRuntimeInfo() + err.Error())
	}
	file_src_ctx, err := alg.AesDecryptCBC(originLic, []byte(utils.Gen32ByteAesKey(inUUIDKey)))
	if err != nil {
		return nil, errors.New(utils.GetRuntimeInfo() + err.Error())
	}
	out := LicGenInfo{}
	out.LicSrc = make(map[string]interface{})
	err = json.Unmarshal(file_src_ctx, &out)
	if err != nil {
		return nil, errors.New(utils.GetRuntimeInfo() + err.Error())
	}
	if out.Version == "" {
		out.Version = "V1"
	}
	return &out, nil
}

func (this *LicInfo) GetLicFileInfo() LicFileInfo {
	return this.lic_file_info
}
