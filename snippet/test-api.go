package snippet

import (
	"github.com/mizuki1412/go-core-kit/library/jsonkit"
	"github.com/mizuki1412/go-core-kit/service/restkit/context"
	"github.com/mizuki1412/go-core-kit/service/restkit/openapi"
	"github.com/mizuki1412/go-core-kit/service/restkit/router"
)

func Init(router *router.Router) {
	tag := "user:用户模块"
	router.Post("", test).Api(openapi.Tag(tag), openapi.Summary("test1"), openapi.ReqParam(testParam{}))
}

type testParam struct {
	Id int32
}

func test(ctx *context.Context) {
	params := testParam{}
	ctx.BindForm(&params)
	println(jsonkit.ToString(params))
}
