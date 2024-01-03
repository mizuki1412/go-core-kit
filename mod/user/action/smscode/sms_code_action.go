package smscode

import (
	context2 "context"
	"github.com/mizuki1412/go-core-kit/class/exception"
	"github.com/mizuki1412/go-core-kit/library/mathkit"
	"github.com/mizuki1412/go-core-kit/library/regexkit"
	"github.com/mizuki1412/go-core-kit/service-third/alismskit"
	"github.com/mizuki1412/go-core-kit/service/rediskit"
	"github.com/mizuki1412/go-core-kit/service/restkit/context"
	"github.com/mizuki1412/go-core-kit/service/restkit/openapi"
	"github.com/mizuki1412/go-core-kit/service/restkit/router"
	"github.com/spf13/cast"
	"time"
)

func Init(router *router.Router) {
	tag := "user:用户模块"
	r := router.Group("/user")
	r.Post("/getVerifyCode", get).Api(openapi.Tag(tag), openapi.Summary("短信验证码获取"), openapi.ReqParam(getParams{}))
}

type getParams struct {
	Phone string `comment:"手机号" validate:"required" trim:"true"`
}

func get(ctx *context.Context) {
	params := getParams{}
	ctx.BindForm(&params)
	if !regexkit.IsPhone(params.Phone) {
		panic(exception.New("手机号码格式错误"))
	}
	sms := ""
	for i := 0; i < 4; i++ {
		sms += cast.ToString(mathkit.RandInt32(0, 10))
	}
	rediskit.Set(context2.Background(), rediskit.GetKeyWithPrefix("sms:"+params.Phone), sms, time.Duration(10)*time.Minute)
	alismskit.Send(alismskit.SendParams{Phone: params.Phone, Data: map[string]any{"code": sms}})
	ctx.JsonSuccess()
}
