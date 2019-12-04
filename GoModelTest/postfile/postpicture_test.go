/*
#Time      :  2019/2/26 下午6:56 
#Author    :  chuangangshen@deepglint.com
#File      :  postfile_test.go
#Software  :  GoLand
*/
package main

import (
	"testing"
)

func TestPostFile(t *testing.T) {
	tests := []struct {
		fileName string
		url      string
	}{
		{"/Users/hpu_scg/gocode/src/temp/GoModelTest/test/slice.txt",
			"http://192.168.100.235:8008/api/upload"},
	}
	for _, tt := range tests {
		resp := PostFile(tt.fileName, tt.url)
		if resp != nil {
			t.Fail()
		}
	}
}
