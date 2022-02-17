package main

/*
#cgo LDFLAGS: -lDECryptoCore
#cgo LDFLAGS: -lpthread -ldl -L.

#include <stdlib.h>

#include "DECryptoCore.h"
*/
import "C"

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"unsafe"
)

func GetUUID() []byte {
	var dest unsafe.Pointer

	C.de_crypto_get_uuid((**C.uint8_t)(unsafe.Pointer(&dest)))
	length := 16
	data := (*[1 << 28]byte)(dest)[:length:length]
	result := make([]byte, length, length)
	copy(result, data)
	C.free(dest)

	return result
}

func ImportLic(fileData []byte) ([]byte, error) {
	var dest unsafe.Pointer

	r := C.de_crypto_lic_import(
		*(**C.uint8_t)(unsafe.Pointer(&fileData)),
		C.uint32_t(len(fileData)),
		(**C.uint8_t)(unsafe.Pointer(&dest)),
	)
	if r < 0 {
		return nil, fmt.Errorf("import lic failed: %d", r)
	}
	data := (*[1 << 28]byte)(dest)[:r:r]
	result := make([]byte, r, r)
	copy(result, data)
	C.free(dest)

	return result, nil
}

func ParseLic(data []byte) (map[string]interface{}, error) {
	var dest unsafe.Pointer

	r := C.de_crypto_lic_parse(
		*(**C.uint8_t)(unsafe.Pointer(&data)),
		C.uint32_t(len(data)),
		(**C.uint8_t)(unsafe.Pointer(&dest)),
	)
	if r < 0 {
		return nil, fmt.Errorf("import lic failed: %d", r)
	}
	data = (*[1 << 28]byte)(dest)[:r:r]
	result := make(map[string]interface{})
	e := json.Unmarshal(data, &result)
	C.free(dest)

	return result, e
}

func main() {
	/* fmt.Println(GetUUID())
	fmt.Println(string(GetUUID()))
	hex.Decode() */
	fmt.Println(hex.EncodeToString(GetUUID()))

	/* f, _ := os.Open("./032e02b4049905fbad061e0700080009.auth.lic")
	data, _ := ioutil.ReadAll(f)

	ilic, e := ImportLic(data)
	fmt.Println(ilic, e)

	plic, e := ParseLic(ilic)
	fmt.Println(plic, e) */
}
