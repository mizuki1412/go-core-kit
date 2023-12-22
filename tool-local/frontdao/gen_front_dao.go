package frontdao

import (
	"fmt"
	"github.com/mizuki1412/go-core-kit/library/arraykit"
	"github.com/mizuki1412/go-core-kit/library/filekit"
	"github.com/mizuki1412/go-core-kit/library/httpkit"
	"github.com/mizuki1412/go-core-kit/library/jsonkit"
	"github.com/mizuki1412/go-core-kit/library/stringkit"
	"github.com/mizuki1412/go-core-kit/service/restkit/openapi"
	"log"
	"sort"
	"strings"
)

type bean struct {
	Name     string
	Imports  []string
	Contents []string
}

func Gen(urlPrefix string) {
	ret, _ := httpkit.Request(httpkit.Req{
		Method: "GET",
		Url:    urlPrefix + "/v3/api-docs",
	})
	var keys []string
	doc := &openapi.ApiDocV3{}
	err := jsonkit.ParseObj(ret, doc)
	if err != nil {
		panic(err)
	}
	for key := range doc.Paths {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	beanMap := map[string]*bean{}
	for _, key := range keys {
		// path存在get/post
		for method, val := range doc.Paths[key] {
			if len(val.Tags) == 0 {
				log.Println("warning: path.tags is null", key)
				continue
			}
			name := stringkit.Split(val.Tags[0], ":")[0]
			if beanMap[name] == nil {
				beanMap[name] = &bean{Name: name}
			}
			b := beanMap[name]
			content := ""
			operationId, _ := openapi.GenOperationId(key, method)
			// 函数描述
			content += fmt.Sprintf("/// %s: %s", operationId, val.Summary)
			// 参数
			for _, e := range val.Parameters {
				require := ""
				if e.Required {
					require = "*"
				}
				// todo reqParam,  reqBody
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
			// todo upload, download
			//for _, ee := range all[key].Get("post.consumes").Array() {
			//	if ee.String() == httpconst.MimeMultipartPOSTForm {
			//		funcName = "upload"
			//		break
			//	}
			//}
			//for _, ee := range all[key].Get("post.produces").Array() {
			//	if ee.String() == httpconst.MimeStream {
			//		funcName = "download"
			//		break
			//	}
			//}
			switch funcName {
			case "request":
				content += fmt.Sprintf(`
export async function %s(%s){
	const {data} = await request('%s'%s)
	return data.data
}
`, operationId, param1, key, param2)
			case "upload":
				content += fmt.Sprintf(`
export async function %s(%s){
	await upload('%s'%s)
}
`, operationId, param1, key, param2)
			case "download":
				content += fmt.Sprintf(`
export async function %s(%s){
	await download('%s'%s)
}
`, operationId, param1, key, param2)
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
