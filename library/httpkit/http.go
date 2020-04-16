package httpkit

import (
	"fmt"
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
	Method string
	Url    string
	Header map[string]string
	// FormData
	FormData map[string]string
}

func Do(reqBean *Req) (string, int, error) {
	data := make(url.Values)
	for key, val := range reqBean.FormData {
		data.Add(key, val)
	}
	if reqBean.Method == "" {
		reqBean.Method = "POST"
	}
	req, err := http.NewRequest(reqBean.Method, reqBean.Url, strings.NewReader(data.Encode()))
	if err != nil {
		return "", 0, err
	}
	if reqBean.Method == "POST" {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}
	for key, val := range reqBean.Header {
		req.Header.Set(key, val)
	}
	resp, err := client.Do(req)
	if err != nil {
		return "", 0, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	return string(body), resp.StatusCode, nil
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
	ret, _, err := Do(&Req{
		Url: "https://www.machplat.com/roms-server-cnc/rest/user/login",
		FormData: map[string]string{
			"username": "@staff",
			"pwd":      "666666",
			"schema":   "cnc",
		},
	})
	if err != nil {
		return
	}
	fmt.Println(gjson.Get(ret, "data.user"))
	fmt.Println(gjson.Parse(ret).Value().(map[string]interface{}))
	fmt.Println("----")
	ret, _, err = Do(&Req{
		Url: "https://www.machplat.com/roms-server-cnc/rest/user/info",
	})
	if err != nil {
		return
	}
}
