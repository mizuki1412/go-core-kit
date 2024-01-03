package frontdao

import (
	"fmt"
	"github.com/mizuki1412/go-core-kit/class/const/httpconst"
	"github.com/mizuki1412/go-core-kit/class/exception"
	"github.com/mizuki1412/go-core-kit/library/c"
	"github.com/mizuki1412/go-core-kit/library/filekit"
	"github.com/mizuki1412/go-core-kit/library/httpkit"
	"github.com/mizuki1412/go-core-kit/library/jsonkit"
	"github.com/mizuki1412/go-core-kit/library/stringkit"
	"github.com/mizuki1412/go-core-kit/service/restkit/openapi"
	"github.com/spf13/cast"
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
	Url         string
	Method      string
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
			f.Url = pathKey
			f.Method = method
			b.Func = append(b.Func, f)
			// 参数
			for _, e := range val.Parameters {
				p := &DaoFuncParam{}
				p.Require = e.Required
				p.In = e.In
				p.Type = e.Schema.Type
				p.Description = e.Description
				p.Name = e.Name
				p.Default = cast.ToString(e.Schema.Default)
				f.Params = append(f.Params, p)
				if p.In == openapi.ParamInPath {
					f.Url = strings.ReplaceAll(f.Url, "{"+p.Name+"}", "${params."+p.Name+"}")
				}
			}
			// 区分函数
			funcName := "request"
			if val.RequestBody != nil && val.RequestBody.Content != nil {
				if _, ok := val.RequestBody.Content[httpconst.MimeMultipartPOSTForm]; ok {
					funcName = "upload"
					b.FlagUpload = true
				}
			}
			if val.Responses != nil {
				for _, body := range val.Responses {
					if _, ok := body.Content[httpconst.MimeStream]; ok {
						funcName = "download"
						b.FlagDownload = true
						break
					}
				}
			}
			if funcName == "request" {
				b.FlagRequest = true
			}
			f.FName = funcName
		}
	}
	// 生成
	for _, b := range beanMap {
		var imports []string
		var contents []string
		if b.FlagRequest {
			imports = append(imports, "request")
		}
		if b.FlagUpload {
			imports = append(imports, "upload")
		}
		if b.FlagDownload {
			imports = append(imports, "download")
		}
		final := fmt.Sprintf("import {%s} from '/lib/request'\n", strings.Join(imports, ","))
		for _, f := range b.Func {
			content := fmt.Sprintf("\n/// %s: %s", f.OperationId, f.Summary)
			paramStr := "params = {"
			var paramStrs []string
			//var inPathVals []string
			for _, p := range f.Params {
				require := c.If[string](p.Require, "*", "")
				content += fmt.Sprintf("\n// %s %s : %s : %s", require, p.Name, p.Type, p.Description)
				paramStrs = append(paramStrs, p.Name+": "+c.If[string](p.Default == "", "null", p.Default))
				//if p.In == openapi.ParamInPath {
				//	inPathVals = append(inPathVals, p.Name)
				//}
			}
			paramStr += strings.Join(paramStrs, ", ") + "}"
			// 是否存在in-path
			//inPathStr := ""
			//if len(inPathVals) > 0 {
			//	inPathStr = "\n\t"
			//}
			// 函数体
			content += fmt.Sprintf("\nexport async function %s(%s){"+
				c.If[string](f.FName == "request", "\n\tconst {data} = await %s(`%s`, params, {method: '%s'})", "\n\tawait %s(`%s`, params, {method: '%s'})")+
				c.If[string](f.FName == "request", "\n\treturn data.data", "")+
				"\n}",
				f.OperationId, paramStr, f.FName, f.Url, f.Method)
			contents = append(contents, content)
		}
		final += strings.Join(contents, "\n")
		_ = filekit.WriteFile("./gen-front-dao/"+b.Name+".js", []byte(final))
	}

}
