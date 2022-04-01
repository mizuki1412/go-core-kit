package context

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/jmoiron/sqlx"
	"github.com/mizuki1412/go-core-kit/class"
	"github.com/mizuki1412/go-core-kit/class/exception"
	"github.com/mizuki1412/go-core-kit/library/jsonkit"
	"github.com/mizuki1412/go-core-kit/library/stringkit"
	"github.com/mizuki1412/go-core-kit/library/timekit"
	"github.com/mizuki1412/go-core-kit/service/logkit"
	"github.com/mizuki1412/go-core-kit/service/sqlkit"
	"github.com/spf13/cast"
	"net/http"
	"reflect"
	"strings"
	"time"
)

type Context struct {
	Proxy    *gin.Context
	Request  *http.Request
	Response gin.ResponseWriter
}

// Set msg per request
func (ctx *Context) Set(key string, val any) {
	ctx.Proxy.Set(key, val)
}

func (ctx *Context) Get(key string) any {
	r, _ := ctx.Proxy.Get(key)
	return r
}

func (ctx *Context) DBTx() *sqlx.Tx {
	var dao *sqlx.Tx
	if ctx.Get("dbtx") == nil {
		ctx.Set("dbtx", sqlkit.StartTx())
	} else {
		dao = ctx.Get("dbtx").(*sqlx.Tx)
	}
	return dao
}
func (ctx *Context) DBTxExist() bool {
	return ctx.Get("dbtx") != nil
}

// data: query, form, json/xml, param

// BindForm bean 指针
func (ctx *Context) BindForm(bean any) {
	switch bean.(type) {
	case map[string]any:
		// query会和form合并 post时 todo
		//allForm := ctx.Proxy.FormValues()
		//for k, v := range allForm {
		//	if len(v) == 1 {
		//		(bean.(map[string]any))[k] = v[0]
		//	} else if len(v) > 1 {
		//		(bean.(map[string]any))[k] = v
		//	}
		//}
		//if len(allForm) == 0 {
		//	// GET
		//	for i := 0; ; i++ {
		//		e := ctx.Proxy.Params().GetEntryAt(i)
		//		if e.Key == "" {
		//			break
		//		}
		//		(bean.(map[string]any))[e.Key] = e.ValueRaw
		//	}
		//}
	default:
		//err := ctx.Proxy.ReadForm(bean)
		//if err != nil {
		//	panic(exception.New("form解析错误"))
		//}
		ctx.bindStruct(bean)
		// validator
		err := Validator.Struct(bean)
		if err != nil {
			if _, ok := err.(*validator.InvalidValidationError); ok {
				panic(exception.New(err.Error()))
			}
			for _, err0 := range err.(validator.ValidationErrors) {
				panic(exception.New("validation failed: " + stringkit.LowerFirst(err0.Field()) + ", need " + err0.Tag()))
			}
		}
	}
}

/// bean:指针
// 实现form/query/json中的数据合并获取。
/// description:"xxx" default:"" trim:"true"
func (ctx *Context) bindStruct(bean any) {
	rt := reflect.TypeOf(bean).Elem()
	rv := reflect.ValueOf(bean).Elem()
	// 取json和取form只能同时进行一次，取完，流被关闭了。
	jsonBody := map[string]any{}
	isJson := strings.Index(ctx.Request.Header.Get("content-type"), "application/json") >= 0
	if isJson {
		_ = ctx.Proxy.ShouldBindJSON(&jsonBody)
	}
	for i := 0; i < rt.NumField(); i++ {
		field := rt.Field(i)
		fieldV := rv.Field(i)
		typeString := field.Type.String()
		key := stringkit.LowerFirst(field.Name)
		// multipart file
		if typeString == "class.File" {
			file, err := ctx.Proxy.FormFile(key)
			// 如果文件流必须存在则检测
			if err != nil && strings.Index(field.Tag.Get("validate"), "required") > -1 {
				panic(exception.New(err.Error()))
			}
			if err == nil {
				f, e := file.Open()
				if e == nil {
					fieldV.Set(reflect.ValueOf(class.File{
						File:   f,
						Header: file,
					}))
				} else {
					logkit.Error(exception.New(e.Error()))
				}
			}
			continue
		}
		// bind struct key
		var val string
		var keyExist bool
		if isJson {
			switch jsonBody[key].(type) {
			case map[string]any:
				val = jsonkit.ToString(jsonBody[key])
			default:
				val = cast.ToString(jsonBody[key])
			}
			_, keyExist = jsonBody[key]
		} else {
			// 判断是否存在key，用于空字符串和无的区分
			val, keyExist = ctx.Proxy.GetPostForm(key)
			if val == "" {
				val = ctx.Proxy.Query(key)
			}
			if val == "" {
				// todo
				val = ctx.Proxy.Param(key)
			}
		}
		// 判断trim
		if field.Tag.Get("trim") == "true" {
			val = strings.TrimSpace(val)
		}
		if val == "" {
			if _, tagExist := field.Tag.Lookup("default"); tagExist {
				// default
				val = field.Tag.Get("default")
				keyExist = true
			}
		}
		switch typeString {
		case "string":
			fieldV.SetString(val)
		case "int32", "int", "int64", "int8", "int16", "byte":
			if !stringkit.IsNull(val) {
				fieldV.SetInt(cast.ToInt64(val))
			}
		case "float64":
			if !stringkit.IsNull(val) {
				fieldV.SetFloat(cast.ToFloat64(val))
			}
		case "bool":
			if !stringkit.IsNull(val) {
				fieldV.SetBool(cast.ToBool(val))
			}
		case "class.Int32":
			if !stringkit.IsNull(val) {
				tmp := class.Int32{Int32: cast.ToInt32(val), Valid: true}
				fieldV.Set(reflect.ValueOf(tmp))
			}
		case "class.Int64":
			if !stringkit.IsNull(val) {
				tmp := class.Int64{Int64: cast.ToInt64(val), Valid: true}
				fieldV.Set(reflect.ValueOf(tmp))
			}
		case "class.Float64":
			if !stringkit.IsNull(val) {
				tmp := class.Float64{Float64: cast.ToFloat64(val), Valid: true}
				fieldV.Set(reflect.ValueOf(tmp))
			}
		case "class.Bool":
			if !stringkit.IsNull(val) {
				tmp := class.Bool{Bool: cast.ToBool(val), Valid: true}
				fieldV.Set(reflect.ValueOf(tmp))
			}
		case "class.String":
			if keyExist {
				tmp := class.String{String: val, Valid: true}
				fieldV.Set(reflect.ValueOf(tmp))
			}
		case "class.ArrInt":
			if !stringkit.IsNull(val) {
				var p []int64
				_ = jsonkit.ParseObj(val, &p)
				tmp := class.ArrInt{Array: p, Valid: true}
				fieldV.Set(reflect.ValueOf(tmp))
			}
		case "class.ArrString":
			if !stringkit.IsNull(val) {
				var p []string
				_ = jsonkit.ParseObj(val, &p)
				tmp := class.ArrString{Array: p, Valid: true}
				fieldV.Set(reflect.ValueOf(tmp))
			}
		case "class.MapString":
			if !stringkit.IsNull(val) {
				var p map[string]any
				_ = jsonkit.ParseObj(val, &p)
				tmp := class.MapString{Map: p, Valid: true}
				fieldV.Set(reflect.ValueOf(tmp))
			}
		case "class.MapStringArr":
			if !stringkit.IsNull(val) {
				var p []map[string]any
				_ = jsonkit.ParseObj(val, &p)
				tmp := class.MapStringArr{Arr: p, Valid: true}
				fieldV.Set(reflect.ValueOf(tmp))
			}
		case "class.Time":
			if !stringkit.IsNull(val) {
				temp := class.Time{}
				if len(val) == 13 && strings.Index(val, "-") < 0 {
					s, err := cast.ToInt64E(val)
					if err == nil {
						//panic(exception.New("time cast error"))
						temp.Set(time.UnixMilli(s))
					}
				} else {
					s, err := timekit.Parse(val)
					if err == nil {
						//panic(exception.New("time cast error"))
						temp.Set(s)
					}
				}
				fieldV.Set(reflect.ValueOf(temp))
			}
		case "class.Decimal":
			if !stringkit.IsNull(val) {
				tmp := class.Decimal{}
				tmp.Set(val)
				fieldV.Set(reflect.ValueOf(tmp))
			}
		}

	}
}
