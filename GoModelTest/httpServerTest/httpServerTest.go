package main

import (
	"crypto/md5"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/julienschmidt/httprouter"
)

func main() {
	/*
		// 方法二
		// router.StartTest(":5678")
	*/
	// run()
	// runApp()
	runGin()
}

func runApp() {
	router := httprouter.New()
	router.PUT("/upload_file", handleUploadFile)
	http.ListenAndServe("127.0.0.1:19999", router)
}

func handleUploadFile(rw http.ResponseWriter, req *http.Request, p httprouter.Params) {
	var err error
	const _24K = (1 << 20) * 24
	if err = req.ParseMultipartForm(_24K); err != nil {
		http.Error(rw, err.Error(), http.StatusForbidden)
		return
	}

	// get md5 value passed with request
	values, ok := req.MultipartForm.Value["md5"]
	if !ok || len(values) < 1 {
		http.Error(rw, "no md5 value passed", http.StatusBadRequest)
		return
	}
	md5Value := values[0]

	// get tar file passed with PostForm request
	var fileAbsPath, md5Hash string
	for _, fileHeaders := range req.MultipartForm.File {
		if len(fileHeaders) != 1 {
			http.Error(rw, "file error", http.StatusBadRequest)
			return
		}
		for _, fileHeader := range fileHeaders {
			file, _ := fileHeader.Open()
			defer file.Close()
			filename := fmt.Sprintf("%s-%d.tar.gz", strings.TrimSuffix(fileHeader.Filename, ".tar.gz"), time.Now().Unix())
			fileAbsPath = filepath.Join("/Users/hpu_scg/gocode/src/temp/GoModelTest/httpServerTest", filename)
			f, err := os.Create(fileAbsPath)
			defer f.Close()
			if err != nil {
				fmt.Println(err)
				http.Error(rw, "file error:"+err.Error(), 501)
				return
			}
			s := make([]byte, 4096)
			for {
				switch nr, err := file.Read(s[:]); true {
				case nr < 0:
					fmt.Fprintf(os.Stderr, "cat: error reading: %s\n", err.Error())
					os.Exit(1)
				case nr == 0: // EOF
					goto L
				case nr > 0:
					fmt.Println(nr)
					f.Write(s[0:nr])
				}
			}
		L:
			data, err := ioutil.ReadFile(fileAbsPath)
			md5Hash = fmt.Sprintf("%x", md5.Sum(data))
			fmt.Println(md5Hash)
			if md5Hash != md5Value {
				http.Error(rw, "bad md5 value caculated", http.StatusBadRequest)
				return
			}
		}
	}
}

func runGin() {
	gin.SetMode(gin.DebugMode)
	e := gin.New()
	server := &http.Server{
		Addr:    ":8686",
		Handler: e,
	}
	e.GET("/api/download", DownloadFile)
	server.ListenAndServe()

}

func DownloadFile(c *gin.Context) {
	c.File("/Users/hpu_scg/Desktop/HaomuT+L-fdu-21.08.27.01.tar.gz")
}

func run() {
	http.HandleFunc("/event", HandedRequest)
	err := http.ListenAndServe(":0", nil)
	if err != nil {
		return
	}
}

func HandedRequest(rw http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	result, err := ioutil.ReadAll(req.Body)
	defer req.Body.Close()
	if err != nil {
		io.WriteString(rw, err.Error())
		return
	}
	ret := string(result)
	fmt.Println(ret)
	io.WriteString(rw, ret)
}
