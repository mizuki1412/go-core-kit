package baidukit

import (
	"github.com/mizuki1412/go-core-kit/v2/class/exception"
	"github.com/mizuki1412/go-core-kit/v2/library/httpkit"
	"github.com/tidwall/gjson"
	"time"
)

var accessKey string
var accessKeyExpire time.Time
var appKey string
var secretKey string

func Init(key, secret string) {
	appKey = key
	secretKey = secret
	// 'https://aip.baidubce.com/oauth/2.0/token?grant_type=client_credentials&client_id=【百度云应用的AK】&client_secret=【百度云应用的SK】'
	res, _ := httpkit.Request(httpkit.Req{
		Url: "https://aip.baidubce.com/oauth/2.0/token",
		FormData: map[string]string{
			"grant_type":    "client_credentials",
			"client_id":     key,
			"client_secret": secret,
		},
	})
	if gjson.Get(res, "error").String() != "" {
		panic(exception.New(gjson.Get(res, "error_description").String()))
	} else {
		accessKey = gjson.Get(res, "access_token").String()
		accessKeyExpire = time.Now().Add(time.Duration(gjson.Get(res, "expires_in").Int()-24*3600) * time.Second)
	}
	//log.Println(accessKey, accessKeyExpire)
}

func checkAccessKey() {
	if time.Now().After(accessKeyExpire) {
		Init(appKey, secretKey)
	}
}
