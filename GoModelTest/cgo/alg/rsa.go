package alg

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"os"
)

type PriKeyType uint

//私钥签名
func Sign(data, privateKey []byte) ([]byte, error) {
	h := sha256.New()
	h.Write(data)
	hashed := h.Sum(nil)
	priv, err := getPriKey(privateKey)
	if err != nil {
		return nil, err
	}
	return rsa.SignPKCS1v15(rand.Reader, priv, crypto.SHA256, hashed)
}

//公钥验证
func SignVer(data, signature, publicKey []byte) error {
	hashed := sha256.Sum256(data)
	//获取公钥
	pub, err := getPubKey(publicKey)
	if err != nil {
		return err
	}
	//验证签名
	return rsa.VerifyPKCS1v15(pub, crypto.SHA256, hashed[:], signature)
}

// 公钥加密
func Encrypt(data, publicKey []byte) ([]byte, error) {
	//获取公钥
	pub, err := getPubKey(publicKey)
	if err != nil {
		return nil, err
	}
	//加密
	return rsa.EncryptPKCS1v15(rand.Reader, pub, data)
}

// 私钥解密,privateKey为pem文件里的字符
func Decrypt(encData, privateKey []byte) ([]byte, error) {
	//解析PKCS1a或者PKCS8格式的私钥
	priv, err := getPriKey(privateKey)
	if err != nil {
		return nil, err
	}
	// 解密
	return rsa.DecryptPKCS1v15(rand.Reader, priv, encData)
}

func getPubKey(publicKey []byte) (*rsa.PublicKey, error) {
	//解密pem格式的公钥
	block, _ := pem.Decode(publicKey)
	if block == nil {
		return nil, errors.New("public key error")
	}
	// 解析公钥
	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	// 类型断言
	if pub, ok := pubInterface.(*rsa.PublicKey); ok {
		return pub, nil
	} else {
		return nil, errors.New("public key error")
	}
}

func getPriKey(privateKey []byte) (*rsa.PrivateKey, error) {
	//获取私钥
	block, _ := pem.Decode(privateKey)
	if block == nil {
		return nil, errors.New("private key error!")
	}
	var priKey *rsa.PrivateKey
	var err error
	priKey, err = x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	return priKey, nil
}

//生成 PKCS1私钥、PKCS8私钥和公钥文件
func GenRsaKey(bits int, pub_key_path, pri_key_path string) error {
	//生成私钥文件
	privateKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return err
	}
	derStream := x509.MarshalPKCS1PrivateKey(privateKey)
	block := &pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: derStream,
	}
	file, err := os.Create(pri_key_path)
	if err != nil {
		return err
	}
	err = pem.Encode(file, block)
	if err != nil {
		return err
	}

	//生成公钥文件
	publicKey := &privateKey.PublicKey
	defPkix, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return err
	}
	block = &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: defPkix,
	}
	file, err = os.Create(pub_key_path)
	if err != nil {
		return err
	}
	err = pem.Encode(file, block)
	if err != nil {
		return err
	}
	return nil
}

func GenRsaKeyPair(inBits int) ([]byte, []byte, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, inBits)
	if err != nil {
		return nil, nil, err
	}
	derStream := x509.MarshalPKCS1PrivateKey(privateKey)
	block := &pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: derStream,
	}
	pemPriKey := pem.EncodeToMemory(block)

	//生成公钥文件
	publicKey := &privateKey.PublicKey
	defPkix, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return nil, nil, err
	}
	block = &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: defPkix,
	}
	pemPubKey := pem.EncodeToMemory(block)
	return pemPubKey, pemPriKey, nil
}

func test_rsa() {
	var tests = []string{
		"abasdf中222国asdffffffffffffffffffffffff" +
			"asdfaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa" +
			"2errrrrrrrrrrrrrrrrrrrrrawdfasvvvvvvvvvvvvvvvvv" +
			"adfasdfffffffffffffffffffffffffffffffffweeeeeeeeeeeee" +
			"asvasdfqwetjo234u2kjnv" +
			"asg2erouolnlv" +
			"asgdasdgafs" +
			"asdfasdfnerlqjouopjkllj" +
			"gwgqerttttttttttttttttttttttttttttttttttttttt6666666666666" +
			"3555554hhhhhhhhhhhhhhhhhhhhh" +
			"bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbfsddddddddd" +
			"avcccccccccccccccc" +
			"sadssssssssssqewrqwerqwerrrrrrr" +
			"qwerqwerqwerqw" +
			"qwerqwevfdgerty43y24t15t" +
			"wertwthgwvtyer2t3452342" +
			"adflajvnjljlkaf",
		"12345678",
		"sjgfjvbj",
	}
	pubKey, priKey, err := GenRsaKeyPair(2048)
	fmt.Println(string(pubKey))
	fmt.Println(string(priKey))
	fmt.Println(err)
	for _, test := range tests {
		sign, err := Sign([]byte(test), priKey)
		fmt.Println(string(sign))
		fmt.Println(err)
		err = SignVer([]byte(test), sign, pubKey)
		if err != nil {
			fmt.Println("Failed :", err.Error())
		}
	}
}
