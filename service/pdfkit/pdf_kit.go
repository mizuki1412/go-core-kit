package pdfkit

import (
	"github.com/SebastiaanKlippert/go-wkhtmltopdf"
	"github.com/mizuki1412/go-core-kit/v2/class/exception"
	"strings"
)

// 格式化完的html数据，destFile-目标pdf文件
func Gen2File(html, destFile string) {
	pdfg, err := wkhtmltopdf.NewPDFGenerator()
	if err != nil {
		panic(exception.New("pdf初始化失败:" + err.Error()))
	}
	pdfg.AddPage(wkhtmltopdf.NewPageReader(strings.NewReader(html)))
	// 注意，这里会请求外链文件，注意网络
	err = pdfg.Create()
	if err != nil {
		panic(exception.New("pdf生成失败: " + err.Error()))
	}
	err = pdfg.WriteFile(destFile)
	if err != nil {
		panic(exception.New("pdf生成失败: " + err.Error()))
	}
}
