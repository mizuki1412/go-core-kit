package swagger

import (
	"github.com/mizuki1412/go-core-kit/cli/configkey"
	"github.com/mizuki1412/go-core-kit/cli/httpconst"
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
	// 由于gin的关系，手动处理 basePath
	basePath := configkit.GetString(configkey.SwaggerBasePath)
	if strings.Index(path, basePath) == 0 {
		path = path[len(basePath):]
	}
	sp := &SwaggerPath{Path: path, Method: method}
	if _, ok := Doc.Paths[path]; !ok {
		Doc.Paths[path] = map[string]map[string]any{}
	}
	Doc.Paths[path][method] = map[string]any{}
	Doc.Paths[path][method]["consumes"] = []string{httpconst.MimePOSTForm}
	Doc.Paths[path][method]["produces"] = []string{httpconst.MimeJSON}
	Doc.Paths[path][method]["operationId"] = path + "-" + method
	Doc.Paths[path][method]["parameters"] = []map[string]any{}
	Doc.Paths[path][method]["responses"] = map[string]any{
		"200": map[string]any{
			"description": "OK",
		},
		"400": map[string]any{
			"description": "参数校验/业务逻辑错误",
		},
	}
	return sp
}

// Param params struct的tags：description，validate:"required"，default
func (swagger *SwaggerPath) Param(param any) *SwaggerPath {
	m := Doc.Paths[swagger.Path][swagger.Method]["parameters"]
	rt := reflect.TypeOf(param)
	for i := 0; i < rt.NumField(); i++ {
		e := map[string]any{}
		tname := stringkit.LowerFirst(rt.Field(i).Type.Name())
		//println(tname)
		switch {
		case tname == "string":
			e["type"] = "string"
			e["in"] = "formData"
		case tname == "file":
			e["type"] = "file"
			/// 参数数据所在位置 eg: query/formData/body
			e["in"] = "formData"
			//Doc.Paths[swagger.Path][swagger.Method]["consumes"] = []string{"multipart/form-data"}
		case strings.Index(tname, "int") == 0:
			e["type"] = "integer"
			e["in"] = "formData"
		case strings.Index(tname, "float") == 0:
			e["type"] = "number"
			e["in"] = "formData"
		case strings.Index(tname, "bool") == 0:
			e["type"] = "boolean"
			e["in"] = "formData"
		case strings.Index(tname, "time") == 0:
			e["type"] = "string"
			e["in"] = "formData"
		default:
			e["type"] = "string"
			e["in"] = "formData"
		}
		e["description"] = rt.Field(i).Tag.Get("description")
		if strings.Contains(rt.Field(i).Tag.Get("validate"), "required") {
			e["required"] = true
		}
		e["name"] = stringkit.LowerFirst(rt.Field(i).Name)
		if v, ok := rt.Field(i).Tag.Lookup("default"); ok {
			e["default"] = v
		}
		m = append(m.([]map[string]any), e)
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
	t := map[string]any{"title": title}
	if len(description) > 0 {
		t["description"] = description[0]
	}
	Doc.tags = append(Doc.tags, t)
	return swagger
}
func (swagger *SwaggerPath) Consume(mime string) *SwaggerPath {
	Doc.Paths[swagger.Path][swagger.Method]["consumes"] = []string{mime}
	return swagger
}
func (swagger *SwaggerPath) ConsumeMultipart() *SwaggerPath {
	return swagger.Consume(httpconst.MimeMultipartPOSTForm)
}
func (swagger *SwaggerPath) Produce(mime string) *SwaggerPath {
	Doc.Paths[swagger.Path][swagger.Method]["produces"] = []string{mime}
	return swagger
}
func (swagger *SwaggerPath) ProduceStream() *SwaggerPath {
	return swagger.Produce(httpconst.MimeStream)
}

var Doc SwaggerDoc

type SwaggerDoc struct {
	Swagger  string         `json:"swagger"`
	Info     map[string]any `json:"info"`
	Host     string         `json:"host"`
	BasePath string         `json:"basePath"`
	tags     []map[string]any
	Paths    map[string]map[string]map[string]any `json:"paths"`
}

func init() {
	Doc.Info = map[string]any{}
	Doc.tags = []map[string]any{}
	Doc.Paths = map[string]map[string]map[string]any{}
}

func (s *SwaggerDoc) ReadDoc() string {
	// match openapi 3.0
	s.Swagger = "2.0"
	s.Info["description"] = configkit.GetString(configkey.SwaggerDescription)
	s.Info["title"] = configkit.GetString(configkey.SwaggerTitle)
	s.Info["version"] = configkit.GetString(configkey.SwaggerVersion, "1.0.0")
	s.Host = configkit.GetString(configkey.SwaggerHost)
	// basePath已经在router中直接加上了，在NewPath中需要额外处理
	s.BasePath = configkit.GetString(configkey.SwaggerBasePath)
	return jsonkit.ToString(*s)
}
