package swagger

import (
	"github.com/mizuki1412/go-core-kit/library/jsonkit"
	"github.com/mizuki1412/go-core-kit/library/stringkit"
	"github.com/mizuki1412/go-core-kit/service/configkit"
	"reflect"
	"strings"
)

type SwaggerPath struct {
	Path   string
	Method string
}

func NewPath(path string, method string) *SwaggerPath {
	sp := &SwaggerPath{Path: path, Method: method}
	if _, ok := Doc.Paths[path]; !ok {
		Doc.Paths[path] = map[string]map[string]interface{}{}
	}
	Doc.Paths[path][method] = map[string]interface{}{}
	Doc.Paths[path][method]["consumes"] = []string{"application/json"}
	Doc.Paths[path][method]["produces"] = []string{"*/*", "application/json"}
	Doc.Paths[path][method]["operationId"] = path
	Doc.Paths[path][method]["parameters"] = []map[string]interface{}{}
	Doc.Paths[path][method]["responses"] = map[string]interface{}{
		"200": map[string]interface{}{
			"description": "OK",
		},
		"400": map[string]interface{}{
			"description": "参数校验错误",
		},
		"401": map[string]interface{}{
			"description": "业务逻辑错误",
		},
	}
	return sp
}

/**
params struct的tags：description，validate:"required"，default
*/
func (swagger *SwaggerPath) Param(param interface{}) *SwaggerPath {
	m := Doc.Paths[swagger.Path][swagger.Method]["parameters"]
	rt := reflect.TypeOf(param)
	for i := 0; i < rt.NumField(); i++ {
		e := map[string]interface{}{}
		tname := rt.Field(i).Type.Name()
		//println(tname)
		if tname == "string" || tname == "String" {
			e["type"] = "string"
			e["in"] = "query"
		} else if tname == "File" {
			e["type"] = "file"
			/// 参数数据所在位置 eg: query/formData/body
			e["in"] = "formData"
			Doc.Paths[swagger.Path][swagger.Method]["consumes"] = []string{"multipart/form-data"}
		}
		e["description"] = rt.Field(i).Tag.Get("description")
		if strings.Contains(rt.Field(i).Tag.Get("validate"), "required") {
			e["required"] = true
		}
		e["name"] = stringkit.LowerFirst(rt.Field(i).Name)
		if v, ok := rt.Field(i).Tag.Lookup("default"); ok {
			e["default"] = v
		}
		m = append(m.([]map[string]interface{}), e)
	}
	Doc.Paths[swagger.Path][swagger.Method]["parameters"] = m
	return swagger
}
func (swagger *SwaggerPath) Description(val string) *SwaggerPath {
	m := Doc.Paths[swagger.Path][swagger.Method]
	m["description"] = val
	return swagger
}
func (swagger *SwaggerPath) Summary(val string) *SwaggerPath {
	m := Doc.Paths[swagger.Path][swagger.Method]
	m["summary"] = val
	return swagger
}
func (swagger *SwaggerPath) Tag(title string, description ...string) *SwaggerPath {
	m := Doc.Paths[swagger.Path][swagger.Method]
	if _, ok := m["tags"]; !ok {
		m["tags"] = []string{}
	}
	m["tags"] = append(m["tags"].([]string), title)
	for _, t := range Doc.tags {
		if t["title"] == title {
			return swagger
		}
	}
	// 新建tags
	t := map[string]interface{}{"title": title}
	if len(description) > 0 {
		t["description"] = description[0]
	}
	Doc.tags = append(Doc.tags, t)
	return swagger
}

var Doc SwaggerDoc

type SwaggerDoc struct {
	Swagger  string                 `json:"swagger"`
	Info     map[string]interface{} `json:"info"`
	Host     string                 `json:"host"`
	BasePath string                 `json:"basePath"`
	tags     []map[string]interface{}
	Paths    map[string]map[string]map[string]interface{} `json:"paths"`
}

func init() {
	Doc.Info = map[string]interface{}{}
	Doc.tags = []map[string]interface{}{}
	Doc.Paths = map[string]map[string]map[string]interface{}{}
}

func (s *SwaggerDoc) ReadDoc() string {
	s.Swagger = "2.0"
	s.Info["description"] = configkit.GetStringD(ConfigKeySwaggerDescription)
	s.Info["title"] = configkit.GetStringD(ConfigKeySwaggerTitle)
	s.Info["version"] = configkit.GetStringD(ConfigKeySwaggerVersion)
	s.Host = configkit.GetStringD(ConfigKeySwaggerHost)
	s.BasePath = configkit.GetStringD(ConfigKeySwaggerBasePath)
	return jsonkit.ToString(*s)
}

/**
标准：https://swagger.io/specification/v2/
swagger-ui可以单独部署，后端只提供doc.json
*/
