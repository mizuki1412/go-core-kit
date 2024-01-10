package alismskit

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/services/dysmsapi"
	"github.com/mizuki1412/go-core-kit/v2/class/exception"
	"github.com/mizuki1412/go-core-kit/v2/cli/configkey"
	"github.com/mizuki1412/go-core-kit/v2/library/jsonkit"
	"github.com/mizuki1412/go-core-kit/v2/service/configkit"
	"strings"
	"sync"
)

var client *dysmsapi.Client
var once sync.Once

func InitSMSClient(keys ...string) {
	once.Do(func() {
		var err error
		var accessKeyId string
		var accessKeySecret string
		if len(keys) == 2 {
			accessKeyId = keys[0]
			accessKeySecret = keys[1]
		} else {
			accessKeyId = configkit.GetString(configkey.AliAccessKey)
			accessKeySecret = configkit.GetString(configkey.AliAccessKeySecret)
		}
		client, err = dysmsapi.NewClientWithAccessKey(configkit.GetString(configkey.AliRegionId, "cn-hangzhou"), accessKeyId, accessKeySecret)
		if err != nil {
			panic(exception.New("sms初始化错误: " + err.Error()))
		}
	})
}

type SendParams struct {
	Phone        string
	Phones       []string
	SignName     string
	TemplateCode string
	Data         map[string]any
}

// Send data例如：{"code":"123456"}
func Send(param SendParams) {
	InitSMSClient()
	phones := ""
	if param.Phone != "" {
		phones = param.Phone
	} else if len(param.Phones) > 0 {
		phones = strings.Join(param.Phones, ",")
	} else {
		panic(exception.New("手机号参数未填"))
	}
	if param.SignName == "" {
		param.SignName = configkit.GetString(configkey.AliSMSSign1)
	}
	if param.TemplateCode == "" {
		param.TemplateCode = configkit.GetString(configkey.AliSMSTemplate1)
	}
	request := dysmsapi.CreateSendSmsRequest()
	request.Scheme = "https"
	request.PhoneNumbers = phones
	request.SignName = param.SignName
	request.TemplateCode = param.TemplateCode
	request.TemplateParam = jsonkit.ToString(param.Data)
	response, err := client.SendSms(request)
	if err != nil {
		panic(exception.New(err.Error()))
	}
	if response.Code != "OK" {
		panic(exception.New(response.Message))
	}
}
