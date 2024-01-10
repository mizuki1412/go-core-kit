package timeget

import (
	"github.com/mizuki1412/go-core-kit/v2/class"
	"github.com/mizuki1412/go-core-kit/v2/service/restkit/context"
	"github.com/mizuki1412/go-core-kit/v2/service/restkit/openapi"
	"github.com/mizuki1412/go-core-kit/v2/service/restkit/router"
)

func Init(router *router.Router) {
	tag := "common:公共模块"
	r := router.Group("/common")
	{
		r.Get("/time", timeGet).Api(openapi.Tag(tag), openapi.Summary("获取服务器时间"), openapi.Response(class.Time{}))
	}
}

func timeGet(ctx *context.Context) {
	ctx.JsonSuccess(class.NewTime())
}
