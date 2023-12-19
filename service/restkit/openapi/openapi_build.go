package openapi

import (
	"github.com/mizuki1412/go-core-kit/class/exception"
	"github.com/mizuki1412/go-core-kit/cli/configkey"
	"github.com/mizuki1412/go-core-kit/library/arraykit"
	"github.com/mizuki1412/go-core-kit/library/stringkit"
	"github.com/mizuki1412/go-core-kit/service/configkit"
	"reflect"
	"strings"
)

var Doc *ApiDocV3

func init() {
	Doc = &ApiDocV3{
		Paths: map[string]map[string]*ApiDocV3PathOperation{},
	}
}

type Builder struct {
	Path *ApiDocV3PathOperation
}

// BuildOpt 定义Functional Options
type BuildOpt func(*Builder)

// GenOperationId path转首字母大写后拼接
func GenOperationId(path, method string) string {
	res := strings.ToLower(method)
	arr := strings.Split(path, "/")
	for _, e := range arr {
		if e == "" {
			continue
		}
		res += stringkit.UpperFirst(e)
	}
	return res
}

func NewBuilder(path string, method string) *Builder {
	op := &ApiDocV3PathOperation{}
	// 初始化response
	op.Responses = map[string]*ApiDocV3ResBody{
		"200": {
			Description: "ok",
			Content: map[string]*ApiDocV3SchemaWrapper{
				"application/json": {
					Schema: &ApiDocV3Schema{
						Type: "string",
					},
				},
			},
		},
	}
	// 生成operationId
	op.OperationId = GenOperationId(path, method)
	if _, ok := Doc.Paths[path]; !ok {
		Doc.Paths[path] = map[string]*ApiDocV3PathOperation{}
	}
	Doc.Paths[path][method] = op
	return &Builder{Path: op}
}

func Description(val string) BuildOpt {
	return func(b *Builder) {
		b.Path.Description = val
	}
}
func Summary(val string) BuildOpt {
	return func(b *Builder) {
		b.Path.Summary = val
	}
}
func Tag(title string, description ...string) BuildOpt {
	return func(b *Builder) {
		b.Path.Tags = []string{title}
		if len(description) > 0 {
			b.Path.Description = description[0]
		}
	}
}

// ReqParam
// params struct的tags：
//
//	comment: 注释
//	validate:"required"
//	default: 默认值
//	in: query,path,header
func ReqParam(param any) BuildOpt {
	return func(b *Builder) {
		rt := reflect.TypeOf(param)
		for i := 0; i < rt.NumField(); i++ {
			e := &ApiDocV3ReqParam{}
			tname := stringkit.LowerFirst(rt.Field(i).Type.Name())
			//println(tname)
			e.Schema = &ApiDocV3Schema{}
			switch {
			case tname == "file":
				panic(exception.New("file类型需要用ReqBody"))
			case strings.Index(tname, "int") == 0:
				e.Schema.Type = "integer"
			case strings.Index(tname, "float") == 0:
				e.Schema.Type = "number"
			case strings.Index(tname, "bool") == 0:
				e.Schema.Type = "boolean"
			default:
				e.Schema.Type = "string"
			}
			e.Description = rt.Field(i).Tag.Get("comment")
			if strings.Contains(rt.Field(i).Tag.Get("validate"), "required") {
				e.Required = true
			}
			e.Name = stringkit.LowerFirst(rt.Field(i).Name)
			if v, ok := rt.Field(i).Tag.Lookup("default"); ok {
				e.Schema.Default = v
			}
			in := rt.Field(i).Tag.Get("in")
			if !arraykit.StringContains([]string{"query", "path", "header", "cookie"}, in) {
				in = "query"
			}
			e.In = in
			b.Path.Parameters = append(b.Path.Parameters, e)
		}
	}
}

func ReqBody(param any) BuildOpt {
	return func(b *Builder) {
		b.Path.RequestBody = &ApiDocV3ReqBody{Content: map[string]*ApiDocV3SchemaWrapper{}}
		// 默认都是json
		parent := &ApiDocV3SchemaWrapper{
			Schema: &ApiDocV3Schema{
				Type:       "object",
				Properties: map[string]*ApiDocV3Schema{},
			},
		}
		key := "application/json"
		rt := reflect.TypeOf(param)
		for i := 0; i < rt.NumField(); i++ {
			e := &ApiDocV3Schema{}
			tname := stringkit.LowerFirst(rt.Field(i).Type.Name())
			switch {
			case tname == "file":
				key = "multipart/form-data"
				e.Type = "string"
				e.Format = "binary"
			case strings.Index(tname, "int") == 0:
				e.Type = "integer"
			case strings.Index(tname, "float") == 0:
				e.Type = "number"
			case strings.Index(tname, "bool") == 0:
				e.Type = "boolean"
			default:
				e.Type = "string"
			}
			e.Description = rt.Field(i).Tag.Get("description")
			if v, ok := rt.Field(i).Tag.Lookup("default"); ok {
				e.Default = v
			}
			// 用name做key
			name := stringkit.LowerFirst(rt.Field(i).Name)
			if strings.Contains(rt.Field(i).Tag.Get("validate"), "required") {
				parent.Schema.Required = append(parent.Schema.Required, name)
			}
			parent.Schema.Properties[name] = e
		}
		b.Path.RequestBody.Content[key] = parent
	}
}

func Response(bean any) BuildOpt {
	// todo response
	return func(b *Builder) {

	}
}

// ResponseStream 返回字节流
func ResponseStream() BuildOpt {
	return func(b *Builder) {
		b.Path.Responses = map[string]*ApiDocV3ResBody{
			"200": {
				Description: "ok",
				Content: map[string]*ApiDocV3SchemaWrapper{
					"application/octet-stream": {},
				},
			},
		}
	}
}

// ReadDoc 返回 api-docs 结果
func (doc *ApiDocV3) ReadDoc() *ApiDocV3 {
	doc.Openapi = "3.1.0"
	if doc.Info == nil {
		doc.Info = &ApiDocV3Info{
			Title:       configkit.GetString(configkey.OpenApiTitle),
			Description: configkit.GetString(configkey.OpenApiDescription),
			License:     nil,
			Version:     configkit.GetString(configkey.OpenApiVersion),
		}
		if configkit.Exist(configkey.OpenApiContactEmail) || configkit.Exist(configkey.OpenApiContactName) || configkit.Exist(configkey.OpenApiContactUrl) {
			doc.Info.Contact = &ApiDocV3InfoContact{
				Name:  configkit.GetString(configkey.OpenApiContactName),
				Url:   configkit.GetString(configkey.OpenApiContactUrl),
				Email: configkit.GetString(configkey.OpenApiContactEmail),
			}
		}
		// todo servers

	}
	return doc
}

// SwaggerConfig swagger-ui不用？
func (doc *ApiDocV3) SwaggerConfig() map[string]any {
	return map[string]any{
		"configUrl":            "/v3/api-docs/swagger-config",
		"oauth2RedirectUrl":    "/swagger-ui/oauth2-redirect.html",
		"operationsSorter":     "alpha",
		"persistAuthorization": true,
		"tagsSorter":           "alpha",
		"url":                  "/v3/api-docs",
		"validatorUrl":         "",
	}
}
