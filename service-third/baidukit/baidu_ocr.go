package baidukit

import (
	"encoding/base64"
	"github.com/mizuki1412/go-core-kit/class/exception"
	"github.com/mizuki1412/go-core-kit/library/httpkit"
	"github.com/mizuki1412/go-core-kit/library/jsonkit"
	"github.com/mizuki1412/go-core-kit/library/timekit"
	"github.com/tidwall/gjson"
)

type OCRParam struct {
	Image []byte
	Pdf   []byte
	Flag  int // 内部使用计数
}

func _getOCRData(param OCRParam) map[string]string {
	data := map[string]string{}
	if len(param.Image) > 0 {
		data["image"] = base64.StdEncoding.EncodeToString(param.Image)
	} else if len(param.Pdf) > 0 {
		data["pdf_file"] = base64.StdEncoding.EncodeToString(param.Pdf)
	}
	return data
}

// OCRVatInvoice 增值税发票识别
// https://cloud.baidu.com/doc/OCR/s/nk3h7xy2t
func OCRVatInvoice(param OCRParam) map[string]any {
	checkAccessKey()
	data := _getOCRData(param)
	res, _ := httpkit.Request(httpkit.Req{
		Url:      "https://aip.baidubce.com/rest/2.0/ocr/v1/vat_invoice?access_token=" + accessKey,
		FormData: data,
	})
	errCode := gjson.Get(res, "error_code").Int()
	if errCode == 18 {
		// QPS超限额
		timekit.Sleep(500)
		if param.Flag == 3 {
			panic(exception.New("qps超时"))
		}
		param.Flag += 1
		return OCRVatInvoice(param)
	} else if errCode != 0 {
		panic(exception.New("baidukit ocr: " + gjson.Get(res, "error_msg").String()))
	}
	return jsonkit.ParseMap(gjson.Get(res, "words_result").String())
}

// OCRTrainTicket 火车票
func OCRTrainTicket(param OCRParam) map[string]any {
	checkAccessKey()
	data := _getOCRData(param)
	res, _ := httpkit.Request(httpkit.Req{
		Url:      "https://aip.baidubce.com/rest/2.0/ocr/v1/train_ticket?access_token=" + accessKey,
		FormData: data,
	})
	errCode := gjson.Get(res, "error_code").Int()
	if errCode == 18 {
		// QPS超限额
		timekit.Sleep(500)
		if param.Flag == 3 {
			panic(exception.New("qps超时"))
		}
		param.Flag += 1
		return OCRTrainTicket(param)
	} else if errCode != 0 {
		panic(exception.New("baidukit ocr: " + gjson.Get(res, "error_msg").String()))
	}
	return jsonkit.ParseMap(gjson.Get(res, "words_result").String())
}

// OCRTollTicket 过路费
func OCRTollTicket(param OCRParam) map[string]any {
	checkAccessKey()
	data := _getOCRData(param)
	res, _ := httpkit.Request(httpkit.Req{
		Url:      "https://aip.baidubce.com/rest/2.0/ocr/v1/toll_invoice?access_token=" + accessKey,
		FormData: data,
	})
	errCode := gjson.Get(res, "error_code").Int()
	if errCode == 18 {
		// QPS超限额
		timekit.Sleep(500)
		if param.Flag == 3 {
			panic(exception.New("qps超时"))
		}
		param.Flag += 1
		return OCRTrainTicket(param)
	} else if errCode != 0 {
		panic(exception.New("baidukit ocr: " + gjson.Get(res, "error_msg").String()))
	}
	return jsonkit.ParseMap(gjson.Get(res, "words_result").String())
}
