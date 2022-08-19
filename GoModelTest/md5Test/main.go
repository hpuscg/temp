package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"os"
)

// md5('Deepglint'+md5(this.ruleForm.password))

func main() {
	SumFileMd5()
}

func StringMd5() {
	m1 := md5.Sum([]byte("Dg1304!@"))
	b1 := m1[:]
	h1 := hex.EncodeToString(b1)
	fmt.Println(h1)

	m2 := md5.Sum([]byte("Deepglint" + h1))
	b2 := m2[:]
	h2 := hex.EncodeToString(b2)
	fmt.Println(h2)

	fmt.Printf("%x\n", md5.Sum([]byte("Deepglint"+fmt.Sprintf("%x", md5.Sum([]byte("Dg1304!@"))))))
}

func SumFileMd5() {
	f, err := os.Open("/Users/hpu_scg/gocode/src/temp/GoModelTest/md5Test/test.txt")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()

	md5Handle := md5.New()
	_, err = io.Copy(md5Handle, f)
	if err != nil {
		fmt.Println(err)
		return
	}
	md := md5Handle.Sum(nil)
	fmt.Printf("%x", md)
}
