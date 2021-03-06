package aliosskit

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/services/sts"
	"github.com/mizuki1412/go-core-kit/class/exception"
	"github.com/mizuki1412/go-core-kit/library/jsonkit"
	"github.com/mizuki1412/go-core-kit/service/configkit"
)

type STSData struct {
	AccessKeyId     string `json:"accessKey"`
	AccessKeySecret string `json:"accessKeySecret"`
	StsToken        string `json:"stsToken"`
	Expiration      string `json:"expiration"`
	Region          string `json:"region"`
	Bucket          string `json:"bucket"`
}

// path是在bucket下的相对路径(eg:*), paths将用于resource
func GetSTS(roleSession, bucket string, paths ...string) STSData {
	ak := configkit.GetString(ConfigKeyAliSTSAccessKey, "")
	aks := configkit.GetString(ConfigKeyAliSTSAccessKeySecret, "")
	role := configkit.GetString(ConfigKeyAliSTSRoleArn, "")
	if ak == "" || aks == "" || role == "" {
		panic(exception.New("sts 必要参数未设置"))
	}
	region := configkit.GetString(ConfigKeyAliSTSRegionId, "cn-hangzhou")
	client, err := sts.NewClientWithAccessKey(region, ak, aks)
	if err != nil {
		panic(exception.New("sts初始化错误: " + err.Error()))
	}
	request := sts.CreateAssumeRoleRequest()
	request.Scheme = "https"
	request.RoleArn = role
	request.RoleSessionName = roleSession
	resource := make([]string, len(paths))
	for i, p := range paths {
		if p[0] != '/' {
			p = "/" + p
		}
		resource[i] = "acs:oss:*:*:" + bucket + p
	}
	request.Policy = jsonkit.ToString(map[string]interface{}{
		"Version": "1",
		"Statement": []map[string]interface{}{
			{
				"Action":   []string{"oss:GetObject", "oss:PutObject", "oss:DeleteObject", "oss:PutObjectAcl", "oss:GetObjectAcl"},
				"Resource": resource,
				"Effect":   "Allow",
			},
		},
	})
	response, err := client.AssumeRole(request)
	if err != nil {
		panic(exception.New("sts request error: " + err.Error()))
	}
	if response.IsSuccess() && response.Credentials.SecurityToken != "" {
		return STSData{
			AccessKeyId:     response.Credentials.AccessKeyId,
			AccessKeySecret: response.Credentials.AccessKeySecret,
			StsToken:        response.Credentials.SecurityToken,
			Expiration:      response.Credentials.Expiration,
			Region:          "oss-" + region,
			Bucket:          bucket,
		}
	} else {
		panic(exception.New("sts response error"))
	}
}
