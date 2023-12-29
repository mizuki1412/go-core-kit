package openapi

import (
	"github.com/mizuki1412/go-core-kit/class/exception"
	"github.com/mizuki1412/go-core-kit/cli/tag"
	"github.com/mizuki1412/go-core-kit/library/arraykit"
	"github.com/mizuki1412/go-core-kit/library/stringkit"
	"reflect"
	"strings"
)

var refPrefix = "#/components/schemas/"

// 对一个类型封装为schema
func buildSchemaByType(t reflect.Type) *ApiDocV3Schema {
	schema := &ApiDocV3Schema{}
	tname := strings.ToLower(t.Name())
	switch {
	case tname == "file":
		schema.Type = SchemaTypeString
		schema.Format = SchemaFormatBinary
	case strings.Contains(tname, "arrint"):
		schema.Type = SchemaTypeArray
		schema.Items = &ApiDocV3Schema{
			Type:   SchemaTypeInteger,
			Format: SchemaFormatInt32,
		}
	case strings.Contains(tname, "arrstring"):
		schema.Type = SchemaTypeArray
		schema.Items = &ApiDocV3Schema{
			Type: SchemaTypeString,
		}
	case strings.Contains(tname, "int"):
		if strings.Index(tname, "int64") >= 0 {
			schema.Format = SchemaFormatInt64
			schema.Type = SchemaTypeString
		} else {
			schema.Format = SchemaFormatInt32
			schema.Type = SchemaTypeInteger
		}
	case strings.Contains(tname, "float"):
		if tname == "float64" {
			schema.Format = SchemaFormatDouble
			schema.Type = SchemaTypeString
		} else {
			schema.Format = SchemaFormatFloat
			schema.Type = SchemaTypeNumber
		}
	case strings.Contains(tname, "decimal"):
		schema.Format = SchemaFormatDouble
		schema.Type = SchemaTypeString
	case strings.Contains(tname, "bool"):
		schema.Type = SchemaTypeBool
	case strings.Contains(tname, "time"):
		schema.Type = SchemaTypeString
		schema.Format = SchemaFormatDateTime
	case strings.Contains(tname, "string"):
		schema.Type = SchemaTypeString
	case strings.Contains(tname, "map") || strings.Contains(tname, "set"):
		schema.Type = SchemaTypeObject
	case t.Kind() == reflect.Pointer:
		schema.Type = SchemaTypeObject
		// field如果是对象，统一ref
		schema.Ref = buildComponentSchema(t.Elem())
	case t.Kind() == reflect.Struct:
		schema.Type = SchemaTypeObject
		schema.Ref = buildComponentSchema(t)
	case t.Kind() == reflect.Slice:
		schema.Type = SchemaTypeArray
		schema.Items = buildSchemaByType(t.Elem())
	default:
		schema.Type = SchemaTypeString
	}
	return schema
}

// 统一封装对象的成员变量为schema，并回调处理
// return content-type(如果需要修改); callback - 回调每个field的处理结果
func buildFieldSchemas(rt reflect.Type, callBack func(s *ApiDocV3Schema, field reflect.StructField)) {
	if rt.Kind() == reflect.Pointer || rt.Kind() == reflect.Slice {
		rt = rt.Elem()
	}
	// 针对pointer array
	if rt.Kind() == reflect.Pointer {
		rt = rt.Elem()
	}
	for i := 0; i < rt.NumField(); i++ {
		if tag.Schema.Contain(rt.Field(i).Tag, tag.SchemaIgnore) {
			continue
		}
		schema := buildSchemaByType(rt.Field(i).Type)
		schema.Description = rt.Field(i).Tag.Get(tag.Comment.Name)
		if v, ok := rt.Field(i).Tag.Lookup(tag.Default.Name); ok {
			schema.Default = v
		}
		callBack(schema, rt.Field(i))
	}
}

// schema封装成一个type=object结构
func buildObjectSchema(rt reflect.Type) *ApiDocV3Schema {
	schema := &ApiDocV3Schema{Properties: map[string]*ApiDocV3Schema{}}
	schema.Type = SchemaTypeObject
	buildFieldSchemas(rt, func(s *ApiDocV3Schema, field reflect.StructField) {
		name := stringkit.LowerFirst(field.Name)
		if tag.Validate.Contain(field.Tag, tag.ValidateRequired) {
			schema.Required = append(schema.Required, name)
		}
		schema.Properties[name] = s
	})
	return schema
}

// 将对象写入到components
func buildComponentSchema(rt reflect.Type) string {
	if rt.Name() == "" {
		panic(exception.New("openapi components schema name is nil"))
	}
	if _, ok := Doc.Components.Schemas[rt.Name()]; ok {
		return refPrefix + rt.Name()
	}
	// 防止循环嵌套调用，需要主动调用handleComponentsTodo
	Doc.Components.Schemas[rt.Name()] = &ApiDocV3Schema{}
	componentsTodoList = append(componentsTodoList, rt)
	return refPrefix + rt.Name()
}

var componentsTodoList []reflect.Type

// 处理components待处理的类型
func handleComponentsTodo() {
	for _, rt := range componentsTodoList {
		schema := buildObjectSchema(rt)
		Doc.Components.Schemas[rt.Name()] = schema
	}
}

func buildReqParamElement(s *ApiDocV3Schema, field reflect.StructField) *ApiDocV3ReqParam {
	e := &ApiDocV3ReqParam{}
	e.Schema = s
	e.Description = field.Tag.Get(tag.Comment.Name)
	if tag.Validate.Contain(field.Tag, tag.ValidateRequired) {
		e.Required = true
	}
	e.Name = stringkit.LowerFirst(field.Name)
	in := field.Tag.Get(tag.ParamIn.Name)
	if !arraykit.StringContains([]string{ParamInQuery, ParamInPath, ParamInHeader, ParamInCookie}, in) {
		in = ParamInQuery
	}
	e.In = in
	return e
}
