/*
#Time      :  2020/5/7 3:19 下午
#Author    :  chuangangshen@deepglint.com
#File      :  gifGraph.go
#Software  :  GoLand
*/
package main

import (
	"github.com/zxbit2011/gifCaptcha"
	"image/color"
	"image/gif"
	"net/http"
)

var captcha = gifCaptcha.New()

func main() {
	//设置颜色
	captcha.SetFrontColor(color.Black, color.RGBA{255, 0, 0, 255}, color.RGBA{0, 0, 255, 255}, color.RGBA{0, 153, 0, 255})
	captcha.SetSize(256,96)
	http.HandleFunc("/img", func(w http.ResponseWriter, r *http.Request) {
		gifData, code := captcha.RangCaptcha()
		println(code)
		_ = gif.EncodeAll(w, gifData)
	})
	_ = http.ListenAndServe(":8080", nil)
}
