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

var refPrefix = "#/components/schemas/"

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

const (
	fromReqParam int = iota
	fromReqBody
	fromResponse
)

// 统一封装对象的成员变量为schema，并回调处理
// return content-type(如果需要修改); callback - 回调每个field的处理结果
func buildSchemas(rt reflect.Type, from int, callBack func(s *ApiDocV3Schema, field reflect.StructField)) string {
	key := ""
	if rt.Kind() != reflect.Struct {
		panic(exception.New("buildSchemas need struct type"))
	}
	for i := 0; i < rt.NumField(); i++ {
		schema := &ApiDocV3Schema{}
		tname := stringkit.LowerFirst(rt.Field(i).Type.Name())
		switch {
		case tname == "file":
			if from == fromReqParam {
				panic(exception.New("file类型需要用ReqBody"))
			} else {
				key = httpconst.MimeMultipartPOSTForm
				schema.Type = "string"
				schema.Format = "binary"
			}
		case strings.Index(tname, "int") >= 0:
			if strings.Index(tname, "int64") >= 0 {
				schema.Format = "int64"
				schema.Type = "string"
			} else {
				schema.Format = "int32"
				schema.Type = "integer"
			}
		case strings.Index(tname, "float") == 0:
			if tname == "float64" {
				schema.Format = "double"
				schema.Type = "string"
			} else {
				schema.Format = "float"
				schema.Type = "number"
			}
		case strings.Index(tname, "bool") == 0:
			schema.Type = "boolean"
		case strings.Index(tname, "time") >= 0:
			schema.Type = "string"
			schema.Format = "date-time"
		case rt.Field(i).Type.Kind() == reflect.Pointer:
			// 内嵌对象
			if from == fromResponse {
				// $ref
				schema.Ref = buildComponentSchema(rt.Field(i).Type.Elem(), from)
			} else {
				schema, key = buildObjectSchema(rt.Field(i).Type.Elem(), from)
			}
		case tname == "":
			// any
			schema.Type = "object"
		default:
			schema.Type = "string"
		}
		schema.Description = rt.Field(i).Tag.Get("comment")
		if v, ok := rt.Field(i).Tag.Lookup("default"); ok {
			schema.Default = v
		}
		callBack(schema, rt.Field(i))
	}
	return key
}

// schema封装成一个type=object结构
func buildObjectSchema(rt reflect.Type, from int) (*ApiDocV3Schema, string) {
	schema := &ApiDocV3Schema{Properties: map[string]*ApiDocV3Schema{}}
	schema.Type = "object"
	key := buildSchemas(rt, from, func(s *ApiDocV3Schema, field reflect.StructField) {
		name := stringkit.LowerFirst(field.Name)
		if strings.Contains(field.Tag.Get("validate"), "required") {
			schema.Required = append(schema.Required, name)
		}
		schema.Properties[name] = s
	})
	return schema, key
}

// 将对象写入到components
func buildComponentSchema(rt reflect.Type, from int) string {
	if rt.Name() == "" {
		panic(exception.New("openapi components schema name is nil"))
	}
	if _, ok := Doc.Components.Schemas[rt.Name()]; ok {
		return refPrefix + rt.Name()
	}
	schema, _ := buildObjectSchema(rt, from)
	Doc.Components.Schemas[rt.Name()] = schema
	return refPrefix + rt.Name()
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
		buildSchemas(rt, fromReqParam, func(s *ApiDocV3Schema, field reflect.StructField) {
			e := &ApiDocV3ReqParam{}
			e.Schema = s
			e.Description = field.Tag.Get("comment")
			if strings.Contains(field.Tag.Get("validate"), "required") {
				e.Required = true
			}
			e.Name = stringkit.LowerFirst(field.Name)
			in := field.Tag.Get("in")
			if !arraykit.StringContains([]string{"query", "path", "header", "cookie"}, in) {
				in = "query"
			}
			e.In = in
			b.Path.Parameters = append(b.Path.Parameters, e)
		})
	}
}

// ReqBody
// 关于schema：因为是请求，schema重复不多，所以就包含在path中，而不用ref。
// 使用body-json，schema统一type为object
func ReqBody(param any) BuildOpt {
	return func(b *Builder) {
		b.Path.RequestBody = &ApiDocV3ReqBody{Content: map[string]*ApiDocV3SchemaWrapper{}}
		key := httpconst.MimeJSON
		rt := reflect.TypeOf(param)
		schema, key0 := buildObjectSchema(rt, fromReqBody)
		if key0 != "" {
			key = key0
		}
		b.Path.RequestBody.Content[key] = &ApiDocV3SchemaWrapper{Schema: schema}
	}
}

// Response schema用ref引用，定义放在components中
func Response(bean any) BuildOpt {
	return func(b *Builder) {
		rt := reflect.TypeOf(bean)
		ref := buildComponentSchema(rt, fromResponse)
		parent := resParentSchema
		parent.Properties[resParentSchemaDataKey].Ref = ref
		b.Path.Responses = map[string]*ApiDocV3ResBody{
			"200": {
				Description: "ok",
				Content: map[string]*ApiDocV3SchemaWrapper{
					"application/json": {
						Schema: &parent,
					},
				},
			},
		}
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

var resParentSchema ApiDocV3Schema
var resParentSchemaDataKey string

// InitResParentSchema 定义返回的父类格式，在response的时候绑定外部父类格式，实际格式在data中
func InitResParentSchema(obj any) {
	rt := reflect.TypeOf(obj)
	for i := 0; i < rt.NumField(); i++ {
		if rt.Field(i).Tag.Get("data") == "true" {
			resParentSchemaDataKey = rt.Field(i).Name
			break
		}
	}
	if resParentSchemaDataKey == "" {
		panic(exception.New("resParentSchemaDataKey cannot nil"))
	}
	s, _ := buildObjectSchema(rt, fromResponse)
	resParentSchema = *s
	println(123)
}
