package baidukit

import (
	"encoding/base64"
	"github.com/mizuki1412/go-core-kit/v2/class/exception"
	"github.com/mizuki1412/go-core-kit/v2/library/httpkit"
	"github.com/mizuki1412/go-core-kit/v2/library/timekit"
	"github.com/tidwall/gjson"
	"time"
)

// FaceDetect 人脸检测
// https://cloud.baidu.com/doc/FACE/s/yk37c1u4t
func FaceDetect(image []byte) map[string]gjson.Result {
	checkAccessKey()
	data := base64.StdEncoding.EncodeToString(image)
	res, _ := httpkit.Request(httpkit.Req{
		Url: "https://aip.baidubce.com/rest/2.0/face/v3/detect?access_token=" + accessKey,
		JsonData: map[string]any{
			"image":      data,
			"image_type": "BASE64",
		},
	})
	errCode := gjson.Get(res, "error_code").Int()
	if errCode == 18 {
		// QPS超限额
		timekit.Sleep(500)
		return FaceDetect(image)
	} else if errCode != 0 {
		panic(exception.New("baidukit: " + gjson.Get(res, "error_msg").String()))
	}
	return gjson.Get(res, "result").Map()
}

// FaceAdd 人脸注册：https://aip.baidubce.com/rest/2.0/face/v3/faceset/user/add
// return face_token
// https://cloud.baidu.com/doc/FACE/s/yk37c1u4t
func FaceAdd(image []byte, groupId, userId string) string {
	checkAccessKey()
	data := base64.StdEncoding.EncodeToString(image)
	res, _ := httpkit.Request(httpkit.Req{
		Url: "https://aip.baidubce.com/rest/2.0/face/v3/faceset/user/add?access_token=" + accessKey,
		JsonData: map[string]any{
			"image":      data,
			"image_type": "BASE64",
			"group_id":   groupId,
			"user_id":    userId,
		},
	})
	errCode := gjson.Get(res, "error_code").Int()
	if errCode == 18 {
		// QPS超限额
		timekit.Sleep(500)
		return FaceAdd(image, groupId, userId)
	} else if errCode != 0 {
		panic(exception.New("baidukit: " + gjson.Get(res, "error_msg").String()))
	}
	return gjson.Get(res, "result").Get("face_token").String()
}

// FaceDel 人脸删除：https://aip.baidubce.com/rest/2.0/face/v3/faceset/user/delete
func FaceDel(groupId, userId, faceToken string) {
	checkAccessKey()
	res, _ := httpkit.Request(httpkit.Req{
		Url: "https://aip.baidubce.com/rest/2.0/face/v3/faceset/user/delete?access_token=" + accessKey,
		JsonData: map[string]any{
			"log_id":     time.Now().UnixMilli(),
			"group_id":   groupId,
			"user_id":    userId,
			"face_token": faceToken,
		},
	})
	errCode := gjson.Get(res, "error_code").Int()
	if errCode == 18 {
		// QPS超限额
		timekit.Sleep(500)
		FaceDel(groupId, userId, faceToken)
	} else if errCode != 0 {
		panic(exception.New("baidukit: " + gjson.Get(res, "error_msg").String()))
	}
}

// FaceSearch 人脸搜索：https://aip.baidubce.com/rest/2.0/face/v3/search
// count 返回的匹配个数; return {userId, score}
func FaceSearch(image []byte, groupId string, count int32) []map[string]any {
	checkAccessKey()
	data := base64.StdEncoding.EncodeToString(image)
	res, _ := httpkit.Request(httpkit.Req{
		Url: "https://aip.baidubce.com/rest/2.0/face/v3/search?access_token=" + accessKey,
		JsonData: map[string]any{
			"image":         data,
			"image_type":    "BASE64",
			"group_id_list": groupId,
			"max_user_num":  count,
		},
	})
	//log.Println(res)
	errCode := gjson.Get(res, "error_code").Int()
	if errCode == 18 {
		// QPS超限额
		timekit.Sleep(500)
		return FaceSearch(image, groupId, count)
	} else if errCode == 222207 {
		// 未找到匹配的用户
		return nil
	} else if errCode != 0 {
		panic(exception.New("baidukit: " + gjson.Get(res, "error_msg").String()))
	}
	users := gjson.Get(res, "result").Get("user_list").Array()
	// 返回：groupId,userId,userInfo,score(用户的匹配得分，推荐阈值80分)
	ret := make([]map[string]any, 0, len(users))
	for _, e := range users {
		ee := map[string]any{
			"groupId":  e.Get("group_id").String(),
			"userId":   e.Get("user_id").String(),
			"userInfo": e.Get("user_info"),
			"score":    e.Get("score").Float(),
		}
		ret = append(ret, ee)
	}
	return ret
}
