package smscode

import (
	"github.com/mizuki1412/go-core-kit/service/restkit/context"
	"github.com/mizuki1412/go-core-kit/service/restkit/router"
)

func Init(router *router.Router) {
	tag := "user:用户模块"
	r := router.Group("/rest/user")
	{
		r.Post("/getVerifyCode", get).Swagger.Tag(tag).Summary("短信验证码获取").Param(getParams{})
	}
}

type getParams struct {
	Phone string `description:"手机号" validate:"required" trim:"true"`
}

func get(ctx *context.Context) {
	params := getParams{}
	ctx.BindForm(&params)
	//alismskit.Send(context2.Background(), params.Phone)
	ctx.JsonSuccess(nil)
}
