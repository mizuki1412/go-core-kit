package httpkit

import (
	"fmt"
	"github.com/mizuki1412/go-core-kit/class/exception"
	"github.com/tidwall/gjson"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"strings"
)

var client *http.Client

func init() {
	client = &http.Client{}
	jar, err := cookiejar.New(nil)
	if err != nil {
		panic(err)
	}
	client.Jar = jar
}

type Req struct {
	Method      string
	Url         string
	Header      map[string]string
	ContentType string
	// FormData
	FormData map[string]string
}

const ContentTypeForm = "application/x-www-form-urlencoded"

func Request(reqBean Req) (string, int) {
	data := make(url.Values)
	for key, val := range reqBean.FormData {
		data.Add(key, val)
	}
	if reqBean.Method == "" {
		reqBean.Method = "POST"
	}
	req, err := http.NewRequest(reqBean.Method, reqBean.Url, strings.NewReader(data.Encode()))
	if err != nil {
		panic(exception.New(err.Error()))
	}
	if reqBean.ContentType == "" {
		req.Header.Set("Content-Type", ContentTypeForm)
	}
	for key, val := range reqBean.Header {
		req.Header.Set(key, val)
	}
	resp, err := client.Do(req)
	if err != nil {
		panic(exception.New(err.Error()))
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
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
	body, err := ioutil.ReadAll(resp.Body)
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
