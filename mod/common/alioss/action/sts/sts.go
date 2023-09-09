package sts

import (
	"github.com/mizuki1412/go-core-kit/cli/configkey"
	"github.com/mizuki1412/go-core-kit/mod/middleware"
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
		r.Post("/get", get).Openapi.Tag(tag).Summary("ali sts 获取")
	}
}

var AdditionFunc = func(uid any, schema string) aliosskit.STSData {
	return aliosskit.GetSTS("user-"+cast.ToString(uid), configkit.GetString(configkey.AliOSSBucketName), "*")
}

func get(ctx *context.Context) {
	c := ctx.GetJwt()
	ctx.JsonSuccess(AdditionFunc(c.Id, c.Ext.GetString("schema")))
}
