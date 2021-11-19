package timeget

import (
	"github.com/mizuki1412/go-core-kit/class"
	"github.com/mizuki1412/go-core-kit/service/restkit/context"
	"github.com/mizuki1412/go-core-kit/service/restkit/router"
	"time"
)

func Init(router *router.Router) {
	tag := "common:公共模块"
	r := router.Group("/rest/common")
	{
		r.Post("/time", timeGet).Swagger.Tag(tag).Summary("获取服务器时间")
	}
}

func timeGet(ctx *context.Context) {
	now := class.Time{Time: time.Now(), Valid: true}
	ctx.JsonSuccess(now)
}
