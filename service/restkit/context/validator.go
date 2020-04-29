package context

import (
	"database/sql/driver"
	"github.com/go-playground/validator/v10"
	"mizuki/project/core-kit/class"
	"reflect"
)

var Validator *validator.Validate

func init() {
	Validator = validator.New()
	Validator.RegisterCustomTypeFunc(validateValuer,
		class.ArrInt{},
		class.Int32{},
		class.String{},
		class.Float64{},
		class.ArrString{},
		class.Bool{},
		class.Int64{},
		class.MapString{},
		class.Time{})
}

func validateValuer(field reflect.Value) interface{} {
	if valuer, ok := field.Interface().(driver.Valuer); ok {
		val, _ := valuer.Value()
		if val != nil && reflect.ValueOf(val).IsZero() {
			// 零值判断，用于required todo
			return 1
		} else if val != nil {
			return val
		}
		// handle the error how you want
	}
	return nil
}

/**
note:
https://github.com/kataras/iris/wiki/Model-validation
https://github.com/go-playground/validator

*/
