package snippet

import (
	"github.com/mizuki1412/go-core-kit/v2/class"
	"github.com/mizuki1412/go-core-kit/v2/library/jsonkit"
	"github.com/mizuki1412/go-core-kit/v2/mod/user/model"
	"github.com/mizuki1412/go-core-kit/v2/service/restkit/context"
	"github.com/mizuki1412/go-core-kit/v2/service/restkit/openapi"
	"github.com/mizuki1412/go-core-kit/v2/service/restkit/router"
)

func Init(router *router.Router) {
	tag := "demo:demo模块"
	router.Get("/param", test).Api(openapi.Tag(tag), openapi.Summary("test1"), openapi.ReqParam(testParam{}))
	router.Post("/param", test).Api(openapi.Tag(tag), openapi.Summary("test2"), openapi.ReqBody(testParam{}))
	router.Post("/path/:id", testPath).Api(openapi.Tag(tag), openapi.Summary("test path"), openapi.ReqBody(testPathParam{}))
	router.Post("/json", testBody).Api(openapi.Tag(tag), openapi.Summary("test json body"), openapi.ReqBody(testBodyParam{}))
	router.Post("/post/file", file).Api(openapi.Tag(tag), openapi.Summary("test-file"), openapi.ReqBody(fileParams{}))
	router.Delete("/delete", test).Api(openapi.Tag(tag), openapi.Summary("test4"), openapi.ReqParam(testParam{}))

	router.Get("/res/1", testRes).Api(openapi.Tag(tag), openapi.Summary("test response1"), openapi.Response(model.Role{}))
	router.Get("/res/2", testRes).Api(openapi.Tag(tag), openapi.Summary("test response2"), openapi.Response(model.UserList{}))
	router.Get("/res/20", testRes).Api(openapi.Tag(tag), openapi.Summary("test response2"), openapi.Response([]*model.User{}))
	router.Get("/res/3", testRes).Api(openapi.Tag(tag), openapi.Summary("test response3"), openapi.ResponseStream())
}

type testParam struct {
	Id        int32        `comment:"标识" validate:"required"`
	ValStr    class.String `comment:"数值"`
	ValLong   class.Int64
	ValDouble class.Float64
}

func test(ctx *context.Context) {
	params := testParam{}
	ctx.BindForm(&params)
	println(jsonkit.ToString(params))
	ctx.JsonSuccess()
}

type testBodyParam struct {
	Id          int32 `comment:"标识" validate:"required"`
	ValLong     class.Int64
	ValDouble   class.Float64 `schema:"ignore"`
	Param       *testParam
	User        model.User
	Params      []testParam
	ArrayInt    class.ArrInt
	ArrayString class.ArrString
}

func testBody(ctx *context.Context) {
	params := testBodyParam{}
	ctx.BindForm(&params)
	println(jsonkit.ToString(params))
	ctx.JsonSuccess()
}

type testPathParam struct {
	Id      int32        `comment:"标识" validate:"required" in:"path"`
	ValStr  class.String `comment:"数值"`
	ValLong class.Int64
}

func testPath(ctx *context.Context) {
	params := testPathParam{}
	ctx.BindForm(&params)
	println(jsonkit.ToString(params))
	ctx.JsonSuccess()
}

type fileParams struct {
	Id   string
	File class.File `validate:"required"`
}

func file(ctx *context.Context) {
	params := fileParams{}
	ctx.BindForm(&params)
	println(len(params.File.GetBytes()))
	ctx.JsonSuccess()
}

func testRes(ctx *context.Context) {
	ctx.JsonSuccess()
}
