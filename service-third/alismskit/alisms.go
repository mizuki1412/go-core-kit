package alismskit

import (
	"errors"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/dysmsapi"
	"github.com/mizuki1412/go-core-kit/class/exception"
	"github.com/mizuki1412/go-core-kit/library/jsonkit"
	"github.com/mizuki1412/go-core-kit/service/configkit"
)

/// phones: xxx,xxxx, data例如：{"code":"123456"}
func Send(phones, signName, templateCode string, data map[string]interface{}) error {
	ak := configkit.GetString(ConfigKeyAliSMSAccessKey, "")
	aks := configkit.GetString(ConfigKeyAliSMSAccessKeySecret, "")
	if ak == "" || aks == "" {
		panic(exception.New("sms accessKey 未设置"))
	}
	client, err := dysmsapi.NewClientWithAccessKey(configkit.GetString(ConfigKeyAliSMSRegionId, "cn-hangzhou"), ak, aks)
	if err != nil {
		panic(exception.New("sms初始化错误: " + err.Error()))
	}
	request := dysmsapi.CreateSendSmsRequest()
	request.Scheme = "https"
	request.PhoneNumbers = phones
	request.SignName = signName
	request.TemplateCode = templateCode
	request.TemplateParam = jsonkit.ToString(data)
	response, err := client.SendSms(request)
	if err != nil {
		return err
	}
	if response.Code != "OK" {
		return errors.New(response.Message)
	} else {
		return nil
	}
}
