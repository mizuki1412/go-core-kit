package frontdao

import (
	"fmt"
	"github.com/mizuki1412/go-core-kit/class/const/httpconst"
	"github.com/mizuki1412/go-core-kit/library/arraykit"
	"github.com/mizuki1412/go-core-kit/library/filekit"
	"github.com/mizuki1412/go-core-kit/library/httpkit"
	"github.com/mizuki1412/go-core-kit/library/stringkit"
	"github.com/tidwall/gjson"
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
	all := gjson.Get(ret, "paths").Map()
	for key := range all {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	beanMap := map[string]*bean{}
	for _, key := range keys {
		// todo get 暂无
		arr := all[key].Get("post.tags").Array()
		if len(arr) == 0 {
			continue
		}
		tag := arr[0].String()
		name := stringkit.Split(tag, ":")[0]
		if beanMap[name] == nil {
			beanMap[name] = &bean{Name: name}
		}
		b := beanMap[name]
		content := ""
		// 函数描述
		content += fmt.Sprintf("/// %s", all[key].Get("post.summary").String())
		// 参数
		for _, e := range all[key].Get("post.parameters").Array() {
			require := ""
			if e.Get("required").Bool() {
				require = "*"
			}
			content += fmt.Sprintf("\n// %s %s : %s : %s", require, e.Get("name").String(), e.Get("type").String(), e.Get("comment"))
		}
		// 函数内容
		var k string
		var k2 string
		param1 := ""
		param2 := ""
		if len(all[key].Get("post.parameters").Array()) > 0 {
			param1 = "params"
			param2 = ", params"
		}
		// key转function name
		if strings.Index(key, "/rest/") > -1 {
			k = key[6:]
		} else {
			k = key[1:]
		}
		for _, k1 := range stringkit.Split(k, "/") {
			k2 += stringkit.UpperFirst(k1)
		}
		// 区分函数
		funcName := "request"
		for _, ee := range all[key].Get("post.consumes").Array() {
			if ee.String() == httpconst.MimeMultipartPOSTForm {
				funcName = "upload"
				break
			}
		}
		for _, ee := range all[key].Get("post.produces").Array() {
			if ee.String() == httpconst.MimeStream {
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
`, k2, param1, key, param2)
		case "upload":
			content += fmt.Sprintf(`
export async function %s(%s){
	await upload('%s'%s)
}
`, k2, param1, key, param2)
		case "download":
			content += fmt.Sprintf(`
export async function %s(%s){
	await download('%s'%s)
}
`, k2, param1, key, param2)
		}
		if !arraykit.StringContains(b.Imports, funcName) {
			b.Imports = append(b.Imports, funcName)
		}
		b.Contents = append(b.Contents, content)
	}
	// 生成
	for _, e := range beanMap {
		final := fmt.Sprintf("import {%s} from '/lib/request'\n\n", strings.Join(e.Imports, ","))
		final += strings.Join(e.Contents, "\n")
		_ = filekit.WriteFile("./gen-front-dao/"+e.Name+".js", []byte(final))
	}

}
