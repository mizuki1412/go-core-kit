package httpkit

import (
	"bufio"
	"bytes"
	"crypto/tls"
	"github.com/mizuki1412/go-core-kit/v2/class/const/httpconst"
	"github.com/mizuki1412/go-core-kit/v2/class/exception"
	"github.com/mizuki1412/go-core-kit/v2/library/filekit"
	"github.com/mizuki1412/go-core-kit/v2/library/jsonkit"
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
	// 忽略证书校验
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
	Method         string
	Url            string
	Header         map[string]string
	ContentType    string // header: Content-Type
	OutputFilePath string // 下载输出至文件

	QueryData  map[string]string // 拼接为url参数
	FormData   map[string]string // body转为表单数据
	JsonData   any               // body转为json格式
	BinaryData []byte

	Timeout       int               // 超时 seconds
	Stream        bool              // 流式处理-字节流
	StreamHandler func(data []byte) // 流式处理逻辑
}

// Request 请求最后阻塞处理
func Request(reqParams Req) (string, int) {
	req, err := genRequest(reqParams)
	if reqParams.Timeout > 0 {
		client.Timeout = time.Duration(reqParams.Timeout) * time.Second
	}
	resp, err := client.Do(req)
	if err != nil {
		panic(exception.New(err.Error()))
	}
	defer resp.Body.Close()
	if reqParams.Stream {
		reader := bufio.NewReader(resp.Body)
		// 这里必须初始0值
		buffer := make([]byte, 4*1024)
		var fout *os.File
		if reqParams.OutputFilePath != "" {
			fout, err = openFile(reqParams.OutputFilePath)
			if err != nil {
				panic(exception.New(err.Error()))
			}
			defer fout.Close()
		}
		for {
			n, err := reader.Read(buffer)
			if err != nil && err != io.EOF {
				panic(exception.New(err.Error()))
			}
			if n > 0 {
				if reqParams.OutputFilePath != "" && fout != nil {
					_, err = fout.Write(buffer[:n])
					if err != nil {
						panic(exception.New(err.Error()))
					}
				} else {
					reqParams.StreamHandler(buffer[:n])
				}
			}
			if err == io.EOF {
				break
			}
			if n == 0 {
				time.Sleep(time.Duration(10) * time.Millisecond)
			}
		}
	} else {
		if reqParams.OutputFilePath != "" {
			fout, err := openFile(reqParams.OutputFilePath)
			if err != nil {
				panic(exception.New(err.Error()))
			}
			defer fout.Close()
			_, err = io.Copy(fout, resp.Body)
			if err != nil {
				panic(exception.New(err.Error()))
			}
		} else {
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				panic(exception.New(err.Error()))
			}
			return string(body), resp.StatusCode
		}
	}
	return "", resp.StatusCode
}

func genRequest(reqParams Req) (*http.Request, error) {
	var req *http.Request
	var err error
	if reqParams.Method == "" {
		reqParams.Method = http.MethodPost
	}
	reqParams.Method = strings.ToUpper(reqParams.Method)
	// 数据体处理
	if reqParams.BinaryData != nil {
		req, err = http.NewRequest(reqParams.Method, reqParams.Url, bytes.NewBuffer(reqParams.BinaryData))
	} else if reqParams.JsonData != nil {
		req, err = http.NewRequest(reqParams.Method, reqParams.Url, bytes.NewBuffer([]byte(jsonkit.ToString(reqParams.JsonData))))
	} else if reqParams.FormData != nil {
		// 自带urlencode转码
		data := make(url.Values)
		for key, val := range reqParams.FormData {
			data.Add(key, val)
		}
		req, err = http.NewRequest(reqParams.Method, reqParams.Url, strings.NewReader(data.Encode()))
	} else if reqParams.QueryData != nil {
		req, err = http.NewRequest(reqParams.Method, reqParams.Url, nil)
		query := req.URL.Query()
		for key, val := range reqParams.QueryData {
			query.Add(key, val)
		}
		req.URL.RawQuery = query.Encode()
	} else {
		req, err = http.NewRequest(reqParams.Method, reqParams.Url, nil)
	}
	if err != nil {
		panic(exception.New(err.Error()))
	}
	if reqParams.ContentType == "" {
		if reqParams.JsonData != nil {
			req.Header.Set("Content-Type", httpconst.ContentTypeJSON)
		} else if reqParams.FormData != nil {
			req.Header.Set("Content-Type", httpconst.ContentTypeForm)
		}
	} else {
		req.Header.Set("Content-Type", reqParams.ContentType)
	}
	for key, val := range reqParams.Header {
		req.Header.Set(key, val)
	}
	return req, err
}

func openFile(fp string) (*os.File, error) {
	if filekit.Exists(fp) {
		return os.OpenFile(fp, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0755)
	} else {
		return os.OpenFile(fp, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0755)
	}
}
