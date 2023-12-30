package frontdao

import (
	"fmt"
	"github.com/mizuki1412/go-core-kit/class/const/httpconst"
	"github.com/mizuki1412/go-core-kit/class/exception"
	"github.com/mizuki1412/go-core-kit/library/arraykit"
	"github.com/mizuki1412/go-core-kit/library/filekit"
	"github.com/mizuki1412/go-core-kit/library/httpkit"
	"github.com/mizuki1412/go-core-kit/library/jsonkit"
	"github.com/mizuki1412/go-core-kit/library/stringkit"
	"github.com/mizuki1412/go-core-kit/service/restkit/openapi"
	"log"
	"strings"
)

type bean struct {
	Name     string   // 模块名称
	Imports  []string // 导入
	Contents []string // 描述
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
	beanMap := map[string]*bean{}
	for pathKey, pathVal := range doc.Paths {
		// 一个请求，path存在get/post
		for method, val := range pathVal {
			if len(val.Tags) == 0 {
				log.Println("warning: path.tags is null", pathKey, method)
				continue
			}
			name := stringkit.Split(val.Tags[0], ":")[0]
			if beanMap[name] == nil {
				beanMap[name] = &bean{Name: name}
			}
			b := beanMap[name]
			content := ""
			operationId, _ := openapi.GenOperationId(pathKey, method)
			// 函数描述
			content += fmt.Sprintf("/// %s: %s", operationId, val.Summary)
			// 参数
			for _, e := range val.Parameters {
				require := ""
				if e.Required {
					require = "*"
				}
				// todo in
				content += fmt.Sprintf("\n// %s %s : %s : %s", require, e.Name, e.Schema.Type, e.Description)
			}

			// 函数内容
			param1 := ""
			param2 := ""
			// todo reqParam,  reqBody
			if len(val.Parameters) > 0 {
				param1 = "params"
				param2 = ", params"
			}
			// 区分函数
			funcName := "request"
			if _, ok := val.RequestBody.Content[httpconst.MimeMultipartPOSTForm]; ok {
				funcName = "upload"
			}
			for _, body := range val.Responses {
				if _, ok := body.Content[httpconst.MimeStream]; ok {
					funcName = "download"
					break
				}
			}
			switch funcName {
			case "request":
				content += fmt.Sprintf(`
export async function %s(%s){
	const {data} = await request('%s'%s)
	return data.data
}
`, operationId, param1, pathKey, param2)
			case "upload":
				content += fmt.Sprintf(`
export async function %s(%s){
	await upload('%s'%s)
}
`, operationId, param1, pathKey, param2)
			case "download":
				content += fmt.Sprintf(`
export async function %s(%s){
	await download('%s'%s)
}
`, operationId, param1, pathKey, param2)
			}
			if !arraykit.StringContains(b.Imports, funcName) {
				b.Imports = append(b.Imports, funcName)
			}
			b.Contents = append(b.Contents, content)
		}
	}
	// 生成
	for _, e := range beanMap {
		final := fmt.Sprintf("import {%s} from '/lib/request'\n\n", strings.Join(e.Imports, ","))
		final += strings.Join(e.Contents, "\n")
		_ = filekit.WriteFile("./gen-front-dao/"+e.Name+".js", []byte(final))
	}

}
