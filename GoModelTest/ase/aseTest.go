/*
#Time      :  2020/3/19 5:53 PM
#Author    :  chuangangshen@deepglint.com
#File      :  aseTest.go
#Software  :  GoLand
*/
package main

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"flag"
	"gitlab.deepglint.com/junkaicao/glog"
	"strings"
)

const (
	LogDir     = "log"
	LogLevel   = "info"
)

var (
	aeskey = []byte("t3.deepglint.com")
	aesiv  = []byte("86-10-62950616-9")
	ps     = []byte("Dg1304!@")
	AlsoToFile = true
)

func main() {
	flag.BoolVar(&AlsoToFile, "alsoToFile", true, "")
	flag.Parse()
	glog.Config(glog.WithAlsoToStd(AlsoToFile), glog.WithFilePath(LogDir), glog.WithLevel(LogLevel))
	EncodePassWord()
	// DecodePassWord()
}

// 加密
func EncodePassWord() {
	ciphertext := AESEncrypt(ps, aeskey, aesiv)
	ret := hex.EncodeToString(ciphertext)
	glog.Infoln(strings.ToUpper(ret))
}

// 解密
func DecodePassWord() {
	en := "B32FA6018771638F277F0BE418708C10"
	data, err := hex.DecodeString(en)
	if err != nil {
		glog.Infoln(err)
	}
	dnData, err := aesCBCDecrypt(data, aeskey, aesiv)
	if err != nil {
		glog.Infoln(err)
	}
	glog.Infoln(string(dnData))
}

// 解密
func aesCBCDecrypt(encryptData, key, iv []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}
	blockSize := block.BlockSize()
	if len(encryptData) < blockSize {
		panic("ciphertext too short")
	}

	if len(encryptData)%blockSize != 0 {
		panic("ciphertext is not a multiple of the block size")
	}

	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(encryptData, encryptData)

	// 解填充
	encryptData = PKCS7UnPadding(encryptData)
	return encryptData, nil
}

//去补码
func PKCS7UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:length-unpadding]
}

//加密
func AESEncrypt(origData, key, iv []byte) []byte {
	//获取block块
	block, _ := aes.NewCipher(key)
	//补码
	origData = PKCS7Padding(origData, block.BlockSize())

	//创建明文长度的数组
	crypted := make([]byte, len(origData))
	//加密模式，
	blockMode := cipher.NewCBCEncrypter(block, iv)
	//加密明文
	blockMode.CryptBlocks(crypted, origData)

	return crypted
}

//补码
func PKCS7Padding(origData []byte, blockSize int) []byte {
	//计算需要补几位数
	padding := blockSize - len(origData)%blockSize
	//在切片后面追加char数量的byte(char)
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(origData, padtext...)
}
