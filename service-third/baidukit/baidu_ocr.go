package baidukit

import (
	"encoding/base64"
	"github.com/mizuki1412/go-core-kit/class/exception"
	"github.com/mizuki1412/go-core-kit/library/httpkit"
	"github.com/mizuki1412/go-core-kit/library/jsonkit"
	"github.com/tidwall/gjson"
)

type OCRVatInvoiceParam struct {
	Image []byte
	Pdf   []byte
}

// OCRVatInvoice 增值税发票识别
// https://cloud.baidu.com/doc/OCR/s/nk3h7xy2t
func OCRVatInvoice(param OCRVatInvoiceParam) map[string]interface{} {
	checkAccessKey()
	data := map[string]string{}
	if len(param.Image) > 0 {
		data["image"] = base64.StdEncoding.EncodeToString(param.Image)
	} else if len(param.Pdf) > 0 {
		data["pdf_file"] = base64.StdEncoding.EncodeToString(param.Pdf)
	}
	res, _ := httpkit.Request(httpkit.Req{
		Url:      "https://aip.baidubce.com/rest/2.0/ocr/v1/vat_invoice?access_token=" + accessKey,
		FormData: data,
	})
	errCode := gjson.Get(res, "error_code").Int()
	if errCode != 0 {
		panic(exception.New("baidukit ocr: " + gjson.Get(res, "error_msg").String()))
	}
	return jsonkit.ParseMap(gjson.Get(res, "words_result").String())
}
