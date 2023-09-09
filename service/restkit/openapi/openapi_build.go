package openapi

import (
	"github.com/mizuki1412/go-core-kit/class/exception"
	"github.com/mizuki1412/go-core-kit/cli/configkey"
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
	Doc.Paths[path] = map[string]*ApiDocV3PathOperation{}
	Doc.Paths[path][method] = op
	return &Builder{Path: op}
}

func (b *Builder) Description(val string) *Builder {
	b.Path.Description = val
	return b
}
func (b *Builder) Summary(val string) *Builder {
	b.Path.Summary = val
	return b
}
func (b *Builder) Tag(title string, description ...string) *Builder {
	b.Path.Tags = []string{title}
	if len(description) > 0 {
		b.Path.Description = description[0]
	}
	return b
}

// ReqParam params struct的tags：description，validate:"required"，default
func (b *Builder) ReqParam(param any) *Builder {
	rt := reflect.TypeOf(param)
	for i := 0; i < rt.NumField(); i++ {
		e := &ApiDocV3Parameter{}
		tname := stringkit.LowerFirst(rt.Field(i).Type.Name())
		//println(tname)
		e.In = "query"
		e.Schema = &ApiDocV3Schema{}
		switch {
		case tname == "file":
			panic(exception.New("file类型需要用ReqBody"))
			//e.Schema.Type = "array"
			//e.Schema.Items = &ApiDocV3Schema{
			//	Type:   "string",
			//	Format: "binary",
			//}
		case strings.Index(tname, "int") == 0:
			e.Schema.Type = "integer"
		case strings.Index(tname, "float") == 0:
			e.Schema.Type = "number"
		case strings.Index(tname, "bool") == 0:
			e.Schema.Type = "boolean"
		default:
			e.Schema.Type = "string"
		}
		e.Description = rt.Field(i).Tag.Get("description")
		if strings.Contains(rt.Field(i).Tag.Get("validate"), "required") {
			e.Required = true
		}
		e.Name = stringkit.LowerFirst(rt.Field(i).Name)
		if v, ok := rt.Field(i).Tag.Lookup("default"); ok {
			e.Schema.Default = v
		}
		b.Path.Parameters = append(b.Path.Parameters, e)
	}
	return b
}

func (b *Builder) ReqBody(param any) *Builder {
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
			e.Type = "array"
			e.Items = &ApiDocV3Schema{
				Type:   "string",
				Format: "binary",
			}
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
	return b
}

func (b *Builder) Response(bean any) *Builder {
	// todo response
	return b
}

// ResponseStream 返回字节流
func (b *Builder) ResponseStream() *Builder {
	b.Path.Responses = map[string]*ApiDocV3ResBody{
		"200": {
			Description: "ok",
			Content: map[string]*ApiDocV3SchemaWrapper{
				"*/*": {
					Schema: &ApiDocV3Schema{
						Type:   "string",
						Format: "binary",
					},
				},
			},
		},
	}
	return b
}

// ReadDoc 返回 api-doc 结果
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
