package openapi

import (
	"github.com/mizuki1412/go-core-kit/v2/class/const/httpconst"
	"github.com/mizuki1412/go-core-kit/v2/class/exception"
	"github.com/mizuki1412/go-core-kit/v2/cli/configkey"
	"github.com/mizuki1412/go-core-kit/v2/cli/tag"
	"github.com/mizuki1412/go-core-kit/v2/library/jsonkit"
	"github.com/mizuki1412/go-core-kit/v2/library/stringkit"
	"github.com/mizuki1412/go-core-kit/v2/service/configkit"
	"github.com/mizuki1412/go-core-kit/v2/service/restkit/context"
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
		if e == "" {
			continue
		}
		// 排除gin中路径参数匹配
		if e[0] == ':' {
			path = strings.ReplaceAll(path, e, "{"+e[1:]+"}")
		}
		if e[0] == ':' || e[0] == '{' || e[0] == '*' {
			continue
		}
		res += stringkit.UpperFirst(e)
	}
	return res, path
}

func NewBuilder(path string, method string) *Builder {
	op := &ApiDocV3PathOperation{}
	// 初始化response
	parent := &ApiDocV3Schema{}
	err := jsonkit.ParseObj(resParentSchema, parent)
	if err != nil {
		panic(exception.New(err.Error()))
	}
	parent.Properties[resParentSchemaDataKey].Type = SchemaTypeObject
	op.Responses = map[string]*ApiDocV3ResBody{
		"200": {
			Description: "ok",
			Content: map[string]*ApiDocV3SchemaWrapper{
				httpconst.MimeJSON: {
					Schema: parent,
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

// ReqParam
// params struct的tags：
//   - comment: 注释
//   - validate:"required"
//   - default: 默认值
//   - in: query,path,header
//   - schema: ignore
func ReqParam(param any) BuildOpt {
	return func(b *Builder) {
		rt := reflect.TypeOf(param)
		buildFieldSchemas(rt, func(s *ApiDocV3Schema, field reflect.StructField) {
			if s.Format == SchemaFormatBinary {
				panic(exception.New("file请使用reqBody"))
			}
			b.Path.Parameters = append(b.Path.Parameters, buildReqParamElement(s, field))
		})
		handleComponentsTodo()
	}
}

// ReqBody
// 关于schema：因为是请求，schema重复不多，所以就包含在path中，而不用ref。
// 使用body-json，schema统一type为object
// 支持 body json和form-data
func ReqBody(param any) BuildOpt {
	return func(b *Builder) {
		b.Path.RequestBody = &ApiDocV3ReqBody{Content: map[string]*ApiDocV3SchemaWrapper{}}
		keyList := []string{httpconst.MimeJSON}
		rt := reflect.TypeOf(param)
		// 逻辑同buildObjectSchema，但需要加上parameter的判断
		schema := &ApiDocV3Schema{Properties: map[string]*ApiDocV3Schema{}}
		schema.Type = SchemaTypeObject
		buildFieldSchemas(rt, func(s *ApiDocV3Schema, field reflect.StructField) {
			// 对file类型处理
			if s.Format == SchemaFormatBinary {
				keyList[0] = httpconst.MimeMultipartPOSTForm
			}
			// 如果存在in，则需要加入params
			if tag.ParamIn.Exist(field.Tag) {
				b.Path.Parameters = append(b.Path.Parameters, buildReqParamElement(s, field))
			} else {
				name := stringkit.LowerFirst(field.Name)
				if tag.Validate.Contain(field.Tag, tag.ValidateRequired) {
					schema.Required = append(schema.Required, name)
				}
				schema.Properties[name] = s
			}
		})
		b.Path.RequestBody.Content[keyList[0]] = &ApiDocV3SchemaWrapper{Schema: schema}
		handleComponentsTodo()
	}
}

// Response schema用ref引用，定义放在components中
func Response(bean any) BuildOpt {
	return func(b *Builder) {
		rt := reflect.TypeOf(bean)
		parent := &ApiDocV3Schema{}
		err := jsonkit.ParseObj(resParentSchema, parent)
		if err != nil {
			panic(exception.New(err.Error()))
		}
		if rt.Kind() == reflect.Slice {
			parent.Properties[resParentSchemaDataKey].Type = SchemaTypeArray
			parent.Properties[resParentSchemaDataKey].Items = buildSchemaByType(rt.Elem())
		} else {
			ref := buildComponentSchema(rt)
			parent.Properties[resParentSchemaDataKey].Type = SchemaTypeObject
			parent.Properties[resParentSchemaDataKey].Ref = ref
		}
		b.Path.Responses = map[string]*ApiDocV3ResBody{
			"200": {
				Description: "ok",
				Content: map[string]*ApiDocV3SchemaWrapper{
					httpconst.MimeJSON: {
						Schema: parent,
					},
				},
			},
		}
		handleComponentsTodo()
	}
}

// ResponseStream 返回字节流
func ResponseStream() BuildOpt {
	return func(b *Builder) {
		b.Path.Responses = map[string]*ApiDocV3ResBody{
			"200": {
				Description: "ok",
				Content: map[string]*ApiDocV3SchemaWrapper{
					httpconst.MimeStream: {Schema: &ApiDocV3Schema{Type: SchemaTypeString, Format: SchemaFormatBinary}},
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

var resParentSchema string
var resParentSchemaDataKey string

// InitResParentSchema 定义返回的父类格式，在response的时候绑定外部父类格式，实际格式在data中
func InitResParentSchema(obj any) {
	rt := reflect.TypeOf(obj)
	for i := 0; i < rt.NumField(); i++ {
		if tag.RetData.Hit(rt.Field(i).Tag) {
			resParentSchemaDataKey = stringkit.LowerFirst(rt.Field(i).Name)
			break
		}
	}
	if resParentSchemaDataKey == "" {
		panic(exception.New("resParentSchemaDataKey cannot nil"))
	}
	s := buildObjectSchema(rt)
	resParentSchema = jsonkit.ToString(s)
}
