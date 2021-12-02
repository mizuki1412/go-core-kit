package frontdao

import (
	"fmt"
	"github.com/mizuki1412/go-core-kit/library/arraykit"
	"github.com/mizuki1412/go-core-kit/library/filekit"
	"github.com/mizuki1412/go-core-kit/library/httpkit"
	"github.com/mizuki1412/go-core-kit/library/stringkit"
	"github.com/tidwall/gjson"
	"sort"
	"strings"
)

func Gen(urlPrefix string, next bool) {
	ret, _ := httpkit.Request(httpkit.Req{
		Method: "GET",
		Url:    urlPrefix + "/swagger/doc",
	})
	var keys []string
	all := gjson.Get(ret, "paths").Map()
	for key := range all {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	tagTemps := make([]string, 0)
	for _, key := range keys {
		var result string
		tag := all[key].Get("post.tags").Array()[0].String()
		name := stringkit.Split(tag, ":")[0]
		if all[key].Get("post.summary").String() != "" {
			result += fmt.Sprintf("\n/// %s", all[key].Get("post.summary").String())
		}
		for _, e := range all[key].Get("post.parameters").Array() {
			require := ""
			if e.Get("required").Bool() {
				require = "*"
			}
			result += fmt.Sprintf("\n// %s %s : %s : %s", require, e.Get("name").String(), e.Get("type").String(), e.Get("description"))
		}
		var k string
		var k2 string
		param1 := ""
		param2 := ""
		if len(all[key].Get("post.parameters").Array()) > 0 {
			if next {
				param1 = "params"
			} else {
				param1 = "params:any"
			}
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
		funName := "postService"
		if next {
			funName = "request"
		}
		tempFlag := ""
		if next {
			tempFlag = ".data"
		}
		result += fmt.Sprintf(`
export async function %s(%s){
	const {data} = await %s('%s'%s)
	return data%s
}
`, k2, param1, funName, key, param2, tempFlag)
		if !arraykit.StringContains(tagTemps, tag) {
			tagTemps = append(tagTemps, tag)
			if next {
				result = "import {request} from 'webkit1412/lib/request';\n" + result
			} else {
				result = "import {postService} from 'web-toolkit/src/case-main/index';\n" + result
			}
		}
		suffer := ".ts"
		if next {
			suffer = ".js"
		}
		_ = filekit.WriteFileAppend("./gen-front-dao/"+name+suffer, []byte(result))
	}
}
