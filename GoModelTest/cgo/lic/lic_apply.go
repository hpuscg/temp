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

const ENC_ORIGIN_KEY = "fgauth123__&^%*^"
const AuthDaysNoLimit = -8888

//const CheckTimePeriod = time.Second * 1
const CheckTimePeriod = 10 //minutes

//lic.apply, not use
type LicApply struct {
	Version          string                 `json:"version"`
	UUID             string                 `json:"uuid"`
	AuthDateGen      string                 `json:"auth_date_gen"`
	AuthDaysTotal    int64                  `json:"auth_days_total"`
	CheckMethod      string                 `json:"check_method"`
	AuthDateDeadline string                 `json:"auth_date_deadline"`
	AuthDaysLeft     int64                  `json:"auth_days_left"`
	AuthCntTotal     int64                  `json:"auth_cnt_total"`
	AuthCntLeft      int64                  `json:"auth_cnt_left"`
	CheckTimePeriod  int64                  `json:"check_period"`
	LicSrc           map[string]interface{} `json:"lic_src"`
}

//format, not use
type FormatInfo struct {
	LicStatus        string                 `json:"lic_status"`
	Version          string                 `json:"version"`
	UUID             string                 `json:"uuid"`
	AuthDateGen      string                 `json:"auth_date_gen"`
	AuthDaysTotal    int64                  `json:"auth_days_total"`
	AuthDaysLeft     int64                  `json:"auth_days_left"`
	AuthDateDeadline string                 `json:"auth_date_deadline"`
	CheckMethod      string                 `json:"check_method"`
	AuthCntLeft      int64                  `json:"auth_cnt_left"`
	LicSrc           map[string]interface{} `json:"lic_src"`
}

//decrease per 10 minutes
func ImportLic(inLicPath, inUUID, inLicApplyPath string) error {
	AuthOne, err := InitAuth("", inLicPath, "", "")
	if err != nil {
		return errors.New(utils.GetRuntimeInfo() + err.Error())
	}
	src_file_info, err := AuthOne.DecLicV2(inUUID)
	if err != nil {
		return errors.New(utils.GetRuntimeInfo() + err.Error())
	}
	if src_file_info.Version == "V1" {
		src_file_info, err := AuthOne.DecLic(inUUID)
		if err != nil {
			return errors.New(utils.GetRuntimeInfo() + err.Error())
		}
		LicApplyOne := LicApply{}
		LicApplyOne.LicSrc = src_file_info
		LicApplyOne.UUID = inUUID
		LicApplyOne.CheckMethod = "countdown"
		//config...
		if auth_date_gen, ok := src_file_info["auth_date_gen"]; ok {
			LicApplyOne.AuthDateGen = auth_date_gen.(string)
		}
		if auth_days_total, ok := src_file_info["auth_days_total"]; ok {
			if val, ok := auth_days_total.(float64); !ok {
				return errors.New("auth_days_total type error! must be integer...")
			} else {
				auth_days_total_int := int64(val)
				if auth_days_total_int == AuthDaysNoLimit {
					LicApplyOne.CheckTimePeriod = CheckTimePeriod
					LicApplyOne.AuthCntTotal = AuthDaysNoLimit
					LicApplyOne.AuthCntLeft = AuthDaysNoLimit
					LicApplyOne.AuthDaysLeft = AuthDaysNoLimit
					LicApplyOne.AuthDaysTotal = AuthDaysNoLimit
				} else {
					LicApplyOne.CheckTimePeriod = CheckTimePeriod
					LicApplyOne.AuthCntTotal = 60 * 24 * auth_days_total_int / int64(LicApplyOne.CheckTimePeriod)
					LicApplyOne.AuthCntLeft = LicApplyOne.AuthCntTotal
					LicApplyOne.AuthDaysLeft = auth_days_total_int
					LicApplyOne.AuthDaysTotal = auth_days_total_int
				}
			}
		}
		//gen check lic file
		err = LicApplyOne.genLicApplyFile(inLicApplyPath, true, inUUID)
		if err != nil {
			return errors.New(utils.GetRuntimeInfo() + err.Error())
		}
	} else {
		LicApplyOne := LicApply{}
		LicApplyOne.Version = src_file_info.Version
		LicApplyOne.LicSrc = src_file_info.LicSrc
		LicApplyOne.CheckMethod = src_file_info.CheckMethod
		LicApplyOne.AuthDateGen = src_file_info.AuthDateGen
		LicApplyOne.AuthDateDeadline = src_file_info.AuthDateDeadline
		LicApplyOne.AuthDaysTotal = src_file_info.AuthDaysTotal
		LicApplyOne.UUID = inUUID

		//config...
		if LicApplyOne.AuthDaysTotal == AuthDaysNoLimit {
			LicApplyOne.CheckTimePeriod = CheckTimePeriod
			LicApplyOne.AuthCntTotal = AuthDaysNoLimit
			LicApplyOne.AuthCntLeft = AuthDaysNoLimit
			LicApplyOne.AuthDaysLeft = AuthDaysNoLimit
		} else {
			LicApplyOne.CheckTimePeriod = CheckTimePeriod
			LicApplyOne.AuthCntTotal = 60 * 24 * LicApplyOne.AuthDaysTotal / int64(LicApplyOne.CheckTimePeriod)
			LicApplyOne.AuthCntLeft = LicApplyOne.AuthCntTotal
			LicApplyOne.AuthDaysLeft = LicApplyOne.AuthDaysTotal
		}
		//gen check lic file
		err = LicApplyOne.genLicApplyFile(inLicApplyPath, true, inUUID)
		if err != nil {
			return errors.New(utils.GetRuntimeInfo() + err.Error())
		}
	}
	return nil
}

//encode to check_lic file
func (this *LicApply) genLicApplyFile(inLicApplyPath string, inOldLicCheck bool, inUUID string) error {
	if this == nil {
		return errors.New(utils.GetRuntimeInfo() + "LicApply null!")
	}
	//check if file reimport according to AuthDateGen > current.AuthDateGen
	if _, err := os.Lstat(inLicApplyPath); err == nil && inOldLicCheck {
		curLicApply, err := ParseLicApplyFile(inLicApplyPath, inUUID)
		if err != nil {
			return errors.New(utils.GetRuntimeInfo() + err.Error())
		}
		//check auth_date_gen of current applied licence is before one of this new licence
		new_auth_date_gen_str := this.AuthDateGen
		curGenTime, err := time.Parse("2006-01-02", this.AuthDateGen)
		if err != nil {
			return errors.New(utils.GetRuntimeInfo() + err.Error())
		}
		newGenTime, err := time.Parse("2006-01-02", new_auth_date_gen_str)
		if err != nil {
			return errors.New(utils.GetRuntimeInfo() + err.Error())
		}
		if newGenTime.Before(curGenTime) || newGenTime.Equal(curGenTime) {
			return errors.New(utils.GetRuntimeInfo() + "ERROR: DO NOT import old or same licence!")
		}

		//check if lic_type not same, throw errors, ONLY same lic_type can overwrite itself
		if cur_lic_type, ok := curLicApply.LicSrc["lic_type"]; ok {
			if cur_lic_type.(string) != "" {
				if new_lic_type, ok := this.LicSrc["lic_type"]; ok {
					if new_lic_type.(string) != cur_lic_type.(string) {
						return errors.New("ERROR: DO NOT import different type licence to same lic.apply!")
					}
				}
			}
		}
	}

	//gen LicApply file
	jb, err := json.Marshal(*this)
	if err != nil {
		return errors.New(utils.GetRuntimeInfo() + err.Error())
	}
	encrypted, err := alg.AesEncryptCBC(jb, []byte(utils.Gen32ByteAesKey(ENC_ORIGIN_KEY)))
	if err != nil {
		return errors.New(utils.GetRuntimeInfo() + err.Error())
	}
	fd, err := os.Create(inLicApplyPath)
	defer fd.Close()
	if err != nil {
		return errors.New(utils.GetRuntimeInfo() + err.Error())
	}
	_, err = fd.Write([]byte(base64.StdEncoding.EncodeToString(encrypted)))
	if err != nil {
		return errors.New(utils.GetRuntimeInfo() + err.Error())
	}
	return nil
}

//decode from check_lic file
func ParseLicApplyFile(inLicApplyPath, inUUID string) (*LicApply, error) {
	fd, err := os.Open(inLicApplyPath)
	defer fd.Close()
	if err != nil {
		return nil, errors.New(utils.GetRuntimeInfo() + err.Error())
	}
	fdb, err := ioutil.ReadAll(fd)
	if err != nil {
		return nil, errors.New(utils.GetRuntimeInfo() + err.Error())
	}
	origin_encrypt, err := base64.StdEncoding.DecodeString(string(fdb))
	if err != nil {
		return nil, errors.New(utils.GetRuntimeInfo() + err.Error())
	}
	aesDec, err := alg.AesDecryptCBC(origin_encrypt, []byte(utils.Gen32ByteAesKey(ENC_ORIGIN_KEY)))
	if err != nil {
		return nil, errors.New(utils.GetRuntimeInfo() + err.Error())
	}
	CheckLicOne := LicApply{}
	if err = json.Unmarshal(aesDec, &CheckLicOne); err != nil {
		return nil, errors.New(utils.GetRuntimeInfo() + err.Error())
	}
	if strings.Trim(CheckLicOne.UUID, " ") != "" && inUUID != CheckLicOne.UUID {
		return nil, errors.New(utils.GetRuntimeInfo() + " | uuid NOT match!")
	}
	return &CheckLicOne, nil
}

//print message
func GenFormatInfo(inLicApplyPath, inUUID string) (*FormatInfo, error) {
	LicApplyOne, err := ParseLicApplyFile(inLicApplyPath, inUUID)
	if err != nil {
		return nil, errors.New(utils.GetRuntimeInfo() + err.Error())
	}
	FormatInfoOne := FormatInfo{}
	FormatInfoOne.UUID = LicApplyOne.UUID
	FormatInfoOne.Version = LicApplyOne.Version
	FormatInfoOne.AuthDateGen = LicApplyOne.AuthDateGen
	FormatInfoOne.AuthDaysTotal = LicApplyOne.AuthDaysTotal
	FormatInfoOne.AuthDateDeadline = LicApplyOne.AuthDateDeadline
	FormatInfoOne.LicSrc = LicApplyOne.LicSrc
	//added one
	FormatInfoOne.AuthDaysLeft = LicApplyOne.AuthDaysLeft
	FormatInfoOne.AuthCntLeft = LicApplyOne.AuthCntLeft
	FormatInfoOne.AuthDateDeadline = LicApplyOne.AuthDateDeadline
	FormatInfoOne.CheckMethod = LicApplyOne.CheckMethod
	if LicApplyOne.CheckMethod == "deadline" {
		if LicApplyOne.AuthDateDeadline == fmt.Sprintf("%d", AuthDaysNoLimit) {
			FormatInfoOne.LicStatus = "OK"
		} else {
			deadlineDate, err := time.Parse("2006-01-02", LicApplyOne.AuthDateDeadline)
			if err != nil {
				FormatInfoOne.LicStatus = "EXPIRED"
			}
			// before deadlineDate is OK
			if time.Now().UTC().Unix() < deadlineDate.UTC().Unix() {
				FormatInfoOne.LicStatus = "OK"
			} else {
				FormatInfoOne.LicStatus = "EXPIRED"
			}
		}
	} else {
		if LicApplyOne.AuthDaysTotal == AuthDaysNoLimit || LicApplyOne.AuthCntLeft > 0 {
			FormatInfoOne.LicStatus = "OK"
		} else {
			FormatInfoOne.LicStatus = "EXPIRED"
		}
	}
	return &FormatInfoOne, nil
}

func Elapsed(inLicApplyPath, inUUID string) error {
	fmt.Println(fmt.Sprintf("%s -> elaped!", time.Now().Format("2006-01-02 15:04:05")))
	checkLicOne, err := ParseLicApplyFile(inLicApplyPath, inUUID)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return err
	}
	if checkLicOne.CheckMethod == "countdown" || checkLicOne.CheckMethod == "" {
		fmt.Println("check_method is countdown, updating applied license!")
		if checkLicOne.AuthDaysTotal == AuthDaysNoLimit {
			return nil
		}
		if checkLicOne.AuthCntLeft <= 0 {
			return nil
		}
		checkLicOne.AuthCntLeft--
		checkLicOne.AuthDaysLeft = checkLicOne.AuthCntLeft * checkLicOne.AuthDaysTotal / checkLicOne.AuthCntTotal
		//add one day if today's count is not 0
		if checkLicOne.AuthCntLeft*checkLicOne.AuthDaysTotal%checkLicOne.AuthCntTotal > 0 {
			checkLicOne.AuthDaysLeft++
		}
		err = checkLicOne.genLicApplyFile(inLicApplyPath, false, inUUID)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return err
		}
	} else {
		fmt.Println("check_method is deadline, no need to update license applied!")
	}
	return nil
}

//check licence if it is not expired
func CheckLicOK(inLicApplyPath, inUUID string) bool {
	checkLicOne, err := ParseLicApplyFile(inLicApplyPath, inUUID)
	if err != nil {
		fmt.Fprint(os.Stderr, err)
		return false
	}
	if checkLicOne.AuthDaysTotal == AuthDaysNoLimit {
		return true
	}
	if checkLicOne.AuthCntLeft > 0 {
		return true
	}
	return false
}
