package sts

import (
	"github.com/mizuki1412/go-core-kit/init/configkey"
	"github.com/mizuki1412/go-core-kit/mod/middleware"
	"github.com/mizuki1412/go-core-kit/mod/user/model"
	"github.com/mizuki1412/go-core-kit/service-third/aliosskit"
	"github.com/mizuki1412/go-core-kit/service/configkit"
	"github.com/mizuki1412/go-core-kit/service/restkit/context"
	"github.com/mizuki1412/go-core-kit/service/restkit/router"
	"github.com/spf13/cast"
)

func Init(router *router.Router) {
	tag := "common:公共模块"
	r := router.Group("/rest/sts")
	r.Use(middleware.AuthUsernameAndPwd())
	{
		r.Post("/get", get).Swagger.Tag(tag).Summary("ali sts 获取")
	}
}

var AdditionFunc = func(user *model.User, schema string) aliosskit.STSData {
	return aliosskit.GetSTS("user"+cast.ToString(user.Id), configkit.GetStringD(configkey.AliOSSBucketName), "*")
}

func get(ctx *context.Context) {
	ctx.JsonSuccess(AdditionFunc(ctx.SessionGetUser(), ctx.SessionGetSchema()))
}
