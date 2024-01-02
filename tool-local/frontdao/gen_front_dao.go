package frontdao

import (
	"fmt"
	"github.com/mizuki1412/go-core-kit/class/const/httpconst"
	"github.com/mizuki1412/go-core-kit/class/exception"
	"github.com/mizuki1412/go-core-kit/library/filekit"
	"github.com/mizuki1412/go-core-kit/library/httpkit"
	"github.com/mizuki1412/go-core-kit/library/jsonkit"
	"github.com/mizuki1412/go-core-kit/library/stringkit"
	"github.com/mizuki1412/go-core-kit/service/restkit/openapi"
	"log"
	"strings"
)

// Dao 一个模块
type Dao struct {
	Name string // 模块名称
	Func []*DaoFunc
	// 标记
	FlagRequest  bool
	FlagUpload   bool
	FlagDownload bool
}

// DaoFunc 模块中的函数
type DaoFunc struct {
	OperationId string
	Summary     string
	Params      []*DaoFuncParam
	// 函数名
	FName string
}
type DaoFuncParam struct {
	Name        string
	In          string
	Require     bool
	Default     string
	Type        string
	Description string
}

func Gen(url string) {
	ret, _ := httpkit.Request(httpkit.Req{
		Method: "GET",
		Url:    url,
	})
	doc := &openapi.ApiDocV3{}
	err := jsonkit.ParseObj(ret, doc)
	if err != nil {
		panic(exception.New(err.Error()))
	}
	beanMap := map[string]*Dao{}
	// 解析
	for pathKey, pathVal := range doc.Paths {
		// 一个请求，path存在get/post
		for method, val := range pathVal {
			if len(val.Tags) == 0 {
				log.Println("warning: path.tags is null", pathKey, method)
				continue
			}
			name := stringkit.Split(val.Tags[0], ":")[0]
			if beanMap[name] == nil {
				beanMap[name] = &Dao{Name: name}
			}
			b := beanMap[name]
			f := &DaoFunc{}
			operationId, _ := openapi.GenOperationId(pathKey, method)
			f.OperationId = operationId
			f.Summary = val.Summary
			b.Func = append(b.Func, f)
			// 参数
			for _, e := range val.Parameters {
				p := &DaoFuncParam{}
				p.Require = e.Required
				p.In = e.In
				p.Type = e.Schema.Type
				p.Description = e.Description
				f.Params = append(f.Params, p)
			}
			// 区分函数
			funcName := "request"
			if _, ok := val.RequestBody.Content[httpconst.MimeMultipartPOSTForm]; ok {
				funcName = "upload"
				b.FlagUpload = true
			}
			for _, body := range val.Responses {
				if _, ok := body.Content[httpconst.MimeStream]; ok {
					funcName = "download"
					b.FlagDownload = true
					break
				}
			}
			if funcName == "request" {
				b.FlagRequest = true
			}
			f.FName = funcName
		}
	}
	// 生成
	for _, e := range beanMap {
		// content += fmt.Sprintf("/// %s: %s", operationId, val.Summary)
		// content += fmt.Sprintf("\n// %s %s : %s : %s", require, e.Name, e.Schema.Type, e.Description)
		//		switch funcName {
		//		case "request":
		//			content += fmt.Sprintf(`
		//export async function %s(%s){
		//	const {data} = await request('%s'%s)
		//	return data.data
		//}
		//`, operationId, param1, pathKey, param2)
		//		case "upload":
		//			content += fmt.Sprintf(`
		//export async function %s(%s){
		//	await upload('%s'%s)
		//}
		//`, operationId, param1, pathKey, param2)
		//		case "download":
		//			content += fmt.Sprintf(`
		//export async function %s(%s){
		//	await download('%s'%s)
		//}
		//`, operationId, param1, pathKey, param2)
		//		}
		final := fmt.Sprintf("import {%s} from '/lib/request'\n\n", strings.Join(e.Imports, ","))
		final += strings.Join(e.Contents, "\n")
		_ = filekit.WriteFile("./gen-front-dao/"+e.Name+".js", []byte(final))
	}

}
