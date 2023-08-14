package alivodkit

import (
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/auth/credentials"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/vod"
	"github.com/mizuki1412/go-core-kit/class"
	"github.com/mizuki1412/go-core-kit/class/exception"
	"github.com/mizuki1412/go-core-kit/cli/configkey"
	"github.com/mizuki1412/go-core-kit/library/cryptokit"
	"github.com/mizuki1412/go-core-kit/library/timekit"
	"github.com/mizuki1412/go-core-kit/service/configkit"
	"github.com/spf13/cast"
	"sync"
	"time"
)

var client *vod.Client
var once sync.Once

// InitVodClient accessKeyId string, accessKeySecret string
func InitVodClient(keys ...string) {
	once.Do(func() {
		var accessKeyId string
		var accessKeySecret string
		if len(keys) == 2 {
			accessKeyId = keys[0]
			accessKeySecret = keys[1]
		} else {
			accessKeyId = configkit.GetStringD(configkey.AliAccessKey)
			accessKeySecret = configkit.GetStringD(configkey.AliAccessKeySecret)
		}
		// 点播服务接入地域
		regionId := configkit.GetString(configkey.AliRegionId, "cn-shanghai")
		// 创建授权对象
		credential := &credentials.AccessKeyCredential{
			AccessKeyId:     accessKeyId,
			AccessKeySecret: accessKeySecret,
		}
		// 自定义config
		config := sdk.NewConfig()
		config.AutoRetry = true     // 失败是否自动重试
		config.MaxRetryTime = 3     // 最大重试次数
		config.Timeout = 3000000000 // 连接超时，单位：纳秒；默认为3秒
		var err error
		client, err = vod.NewClientWithOptions(regionId, config, credential)
		if err != nil {
			panic(exception.New(err.Error()))
		}
	})
}

type SearchMediaInfoListParam struct {
	Title string
}

// SearchMediaInfoList return {id, duration/hh:mm:ss, size/MB,cover,status, dt}
// 注意ram账号授权相关api的操作
// https://help.aliyun.com/document_detail/86044.htm?spm=a2c4g.11186623.0.0.3cd418971FU74d#doc-api-vod-SearchMedia
func SearchMediaInfoList(p SearchMediaInfoListParam) []map[string]any {
	InitVodClient()
	request := vod.CreateSearchMediaRequest()
	request.Fields = "Duration,Size,CoverURL,Status"
	request.Match = fmt.Sprintf("Title='%s'", p.Title)
	res, err := client.SearchMedia(request)
	if err != nil {
		panic(exception.New(err.Error()))
	}
	ret := make([]map[string]any, 0, len(res.MediaList))
	for _, e := range res.MediaList {
		ret = append(ret, map[string]any{
			"id":       e.MediaId,
			"dt":       timekit.ParseD(e.CreationTime).Format(timekit.TimeLayout),
			"cover":    e.Video.CoverURL,
			"status":   e.Video.Status,
			"size":     class.NewDecimal(e.Video.Size).Div(class.NewDecimal(1024)).DivRound(class.NewDecimal(1024), 2).Float64(),
			"duration": timekit.FormatSecondHMS(cast.ToInt64(e.Video.Duration), false),
		})
	}
	return ret
}

// GetPlayInfo https://help.aliyun.com/document_detail/56124.html
func GetPlayInfo(videoId string) []map[string]any {
	InitVodClient()
	request := vod.CreateGetPlayInfoRequest()
	request.VideoId = videoId
	// Definition Url
	res, err := client.GetPlayInfo(request)
	if err != nil {
		panic(exception.New(err.Error()))
	}
	ret := make([]map[string]any, 0, len(res.PlayInfoList.PlayInfo))
	for _, e := range res.PlayInfoList.PlayInfo {
		ret = append(ret, map[string]any{
			"definition": getDefinitionName(e.Definition),
			"url":        e.PlayURL,
		})
	}
	return ret
}

func getDefinitionName(val string) string {
	switch val {
	case "FD":
		return "流畅"
	case "LD":
		return "标清"
	case "SD":
		return "高清"
	case "HD":
		return "超清"
	case "OD":
		return "原画"
	case "SQ":
		return "普通音质"
	case "HQ":
		return "高音质"
	}
	return val
}

// GenUrlAuthKey url鉴权中的 auth_key=timestamp-rand-uid-md5hash
func GenUrlAuthKey(uri string, key string) string {
	timestamp := time.Now().Unix()
	rand := 0
	uid := 0
	// URI-timestamp-rand-uid-PrivateKey
	md5hash := cryptokit.MD5(fmt.Sprintf("%s-%d-%d-%d-%s", uri, timestamp, rand, uid, key))
	return fmt.Sprintf("%d-%d-%d-%s", timestamp, rand, uid, md5hash)
}
