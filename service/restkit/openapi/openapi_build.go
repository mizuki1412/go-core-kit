package openapi

import (
	"github.com/mizuki1412/go-core-kit/class/const/httpconst"
	"github.com/mizuki1412/go-core-kit/class/exception"
	"github.com/mizuki1412/go-core-kit/cli/configkey"
	"github.com/mizuki1412/go-core-kit/library/arraykit"
	"github.com/mizuki1412/go-core-kit/library/stringkit"
	"github.com/mizuki1412/go-core-kit/service/configkit"
	"github.com/mizuki1412/go-core-kit/service/restkit/context"
	"reflect"
	"strings"
)

var Doc *ApiDocV3

func init() {
	Doc = &ApiDocV3{
		Paths: map[string]map[string]*ApiDocV3PathOperation{},
		Tags:  []*ApiDocV3Tag{},
		Components: &ApiDocV3ComponentObj{
			Schemas: map[string]*ApiDocV3Schema{},
		},
	}
	InitResParentSchema(context.RestRet{})
}

// Builder 单条路径的builder
type Builder struct {
	Path *ApiDocV3PathOperation
}

// BuildOpt 定义Functional Options
type BuildOpt func(*Builder)

// GenOperationId path转首字母大写后拼接，同时把其中的路径参数标识出来，暂时处理:
func GenOperationId(path, method string) (string, string) {
	res := strings.ToLower(method)
	arr := strings.Split(path, "/")
	for _, e := range arr {
		// 排除gin中路径参数匹配
		if e[0] == ':' {
			path = strings.ReplaceAll(path, e, "{"+e[1:]+"}")
		}
		if e == "" || e[0] == ':' || e[0] == '{' || e[0] == '*' {
			continue
		}
		res += stringkit.UpperFirst(e)
	}
	return res, path
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
	op.OperationId, path = GenOperationId(path, method)
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
func Tag(tag string) BuildOpt {
	return func(b *Builder) {
		b.Path.Tags = []string{tag}
		// 添加 doc.tags
		var target *ApiDocV3Tag
		for _, e := range Doc.Tags {
			if e.Name == tag {
				target = e
				break
			}
		}
		if target == nil {
			target = &ApiDocV3Tag{
				Name:        tag,
				Description: tag,
			}
			Doc.Tags = append(Doc.Tags, target)
		}
	}
}

// todo 统一封装schema
// return content-type(如果需要修改)
func reqPropTypeHandle(fieldType string, schema *ApiDocV3Schema, reqBody bool) string {
	key := ""
	switch {
	case fieldType == "file":
		if reqBody {
			key = httpconst.MimeMultipartPOSTForm
			schema.Type = "string"
			schema.Format = "binary"
		} else {
			panic(exception.New("file类型需要用ReqBody"))
		}
	case strings.Index(fieldType, "int") == 0:
		schema.Type = "integer"
		if fieldType == "int64" {
			schema.Format = "int64"
		} else {
			schema.Format = "int32"
		}
	case strings.Index(fieldType, "float") == 0:
		schema.Type = "number"
		if fieldType == "float32" {
			schema.Format = "float"
		} else {
			schema.Format = "double"
		}
	case strings.Index(fieldType, "bool") == 0:
		schema.Type = "boolean"
	case strings.Index(fieldType, "time") >= 0:
		schema.Type = "string"
		schema.Format = "date-time"
	default:
		schema.Type = "string"
	}
	return key
}

// ReqParam
// params struct的tags：
//   - comment: 注释
//   - validate:"required"
//   - default: 默认值
//   - in: query,path,header
func ReqParam(param any) BuildOpt {
	return func(b *Builder) {
		rt := reflect.TypeOf(param)
		if rt.Kind() != reflect.Struct {
			panic(exception.New("openapi param need struct"))
		}
		for i := 0; i < rt.NumField(); i++ {
			e := &ApiDocV3ReqParam{}
			tname := stringkit.LowerFirst(rt.Field(i).Type.Name())
			//println(tname)
			e.Schema = &ApiDocV3Schema{}
			reqPropTypeHandle(tname, e.Schema, false)
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

// ReqBody
// 关于schema：因为是请求，schema重复不多，所以就包含在path中，而不用ref。
// 使用body-json，schema统一type为object
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
		key := httpconst.MimeJSON
		rt := reflect.TypeOf(param)
		if rt.Kind() != reflect.Struct {
			panic(exception.New("openapi param need struct"))
		}
		for i := 0; i < rt.NumField(); i++ {
			e := &ApiDocV3Schema{}
			tname := stringkit.LowerFirst(rt.Field(i).Type.Name())
			// todo 内嵌对象？
			key0 := reqPropTypeHandle(tname, e, true)
			if key0 != "" {
				key = key0
			}
			e.Description = rt.Field(i).Tag.Get("comment")
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

// Response todo schema用ref引用，定义放在components中
func Response(bean any) BuildOpt {
	return func(b *Builder) {
		//rt := reflect.TypeOf(bean)
		//if rt.Kind() != reflect.Struct {
		//	panic(exception.New("openapi param need struct"))
		//}
		//schema := &ApiDocV3Schema{}
		//for i := 0; i < rt.NumField(); i++ {
		//	tname := stringkit.LowerFirst(rt.Field(i).Type.Name())
		//	reqPropTypeHandle(tname, e, false)
		//	e.Description = rt.Field(i).Tag.Get("comment")
		//	// 用name做key
		//	name := stringkit.LowerFirst(rt.Field(i).Name)
		//	println(name)
		//	//parent.Schema.Properties[name] = e
		//}
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

var resParentSchemaProps = map[string]*ApiDocV3Schema{}

// 返回的父类格式的data区域
var resParentSchemaPropData *ApiDocV3Schema

// InitResParentSchema 定义返回的父类格式，在response的时候绑定外部父类格式，实际格式在data中
func InitResParentSchema(obj any) {
	rt := reflect.TypeOf(obj)
	if rt.Kind() != reflect.Struct {
		panic(exception.New("InitResParentSchema need struct"))
	}
	for i := 0; i < rt.NumField(); i++ {
		e := &ApiDocV3Schema{}
		tname := stringkit.LowerFirst(rt.Field(i).Type.Name())
		reqPropTypeHandle(tname, e, true)
		e.Description = rt.Field(i).Tag.Get("comment")
		// 用name做key
		//name := stringkit.LowerFirst(rt.Field(i).Name)
		//parent.Schema.Properties[name] = e
	}
}
