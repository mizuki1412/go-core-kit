package httpkit

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"github.com/mizuki1412/go-core-kit/class/exception"
	"github.com/mizuki1412/go-core-kit/library/jsonkit"
	"github.com/tidwall/gjson"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"strings"
	"time"
)

var client *http.Client

func init() {
	// 忽略证书校验 todo
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client = &http.Client{Transport: tr}
	jar, err := cookiejar.New(nil)
	if err != nil {
		panic(err)
	}
	client.Jar = jar
}

// Req 填写FormData、JsonData时可缺省contentType
type Req struct {
	Method      string
	Url         string
	Header      map[string]string
	ContentType string
	FormData    map[string]string
	JsonData    interface{}
	BinaryData  []byte
	Timeout     int // seconds
}

const ContentTypeForm = "application/x-www-form-urlencoded; charset=utf-8"
const ContentTypeJSON = "application/json; charset=utf-8"

func Request(reqBean Req) (string, int) {
	if reqBean.Method == "" {
		reqBean.Method = http.MethodPost
	}
	var req *http.Request
	var err error
	if reqBean.BinaryData != nil {
		req, err = http.NewRequest(reqBean.Method, reqBean.Url, bytes.NewBuffer(reqBean.BinaryData))
	} else if reqBean.JsonData != nil {
		req, err = http.NewRequest(reqBean.Method, reqBean.Url, bytes.NewBuffer([]byte(jsonkit.ToString(reqBean.JsonData))))
	} else {
		// 自带urlencode转码
		data := make(url.Values)
		for key, val := range reqBean.FormData {
			data.Add(key, val)
		}
		req, err = http.NewRequest(reqBean.Method, reqBean.Url, strings.NewReader(data.Encode()))
	}
	if err != nil {
		panic(exception.New(err.Error()))
	}
	if reqBean.ContentType == "" {
		if reqBean.JsonData != nil {
			req.Header.Set("Content-Type", ContentTypeJSON)
		} else {
			req.Header.Set("Content-Type", ContentTypeForm)
		}
	} else {
		req.Header.Set("Content-Type", reqBean.ContentType)
	}
	for key, val := range reqBean.Header {
		req.Header.Set(key, val)
	}
	if reqBean.Timeout > 0 {
		client.Timeout = time.Duration(reqBean.Timeout) * time.Second
	}
	resp, err := client.Do(req)
	if err != nil {
		panic(exception.New(err.Error()))
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(exception.New(err.Error()))
	}
	return string(body), resp.StatusCode
}

func DownloadToFile(url string, filePath string) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic(err)
	}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	// write
	fout, err := os.Create(filePath)
	if err != nil {
		panic(err)
	}
	defer fout.Close()
	_, err = fout.Write(body)
	if err != nil {
		panic(err)
	}
}

func Demo() {
	ret, _ := Request(Req{
		Url: "https://www.machplat.com/roms-rest-cnc/rest/user/login",
		FormData: map[string]string{
			"username": "admin",
			"pwd":      "123",
		},
	})
	fmt.Println(gjson.Get(ret, "data.user"))
	fmt.Println(gjson.Parse(ret).Value().(map[string]interface{}))
	fmt.Println("----")
	ret, _ = Request(Req{
		Url: "https://www.machplat.com/roms-rest-cnc/rest/user/info",
	})
}
