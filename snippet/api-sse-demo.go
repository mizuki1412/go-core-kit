package snippet

import (
	"github.com/mizuki1412/go-core-kit/v2/service/restkit/context"
	"github.com/mizuki1412/go-core-kit/v2/service/restkit/router"
	"github.com/mizuki1412/go-core-kit/v2/service/restkit/ssehelper"
)

func InitSSE(router *router.Router) {
	r := router.Group("/sse")
	//r.Use(middleware.AuthUsernameAndPwd())
	{
		r.Get("/:id", sse)
		r.Get("/send", sseSend)
	}
}

type SSEParams struct {
	Id string
}

func sse(ctx *context.Context) {
	params := SSEParams{}
	ctx.BindForm(&params)
	ssehelper.ServiceClient(params.Id, ctx)
}

func sseSend(ctx *context.Context) {
	ssehelper.ToSend("123", "abcd")
	ctx.JsonSuccess()
}
