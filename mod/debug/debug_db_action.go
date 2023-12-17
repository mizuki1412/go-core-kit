package debug

import (
	"github.com/mizuki1412/go-core-kit/mod/middleware"
	"github.com/mizuki1412/go-core-kit/service/logkit"
	"github.com/mizuki1412/go-core-kit/service/restkit/context"
	"github.com/mizuki1412/go-core-kit/service/restkit/router"
	"github.com/mizuki1412/go-core-kit/service/sqlkit"
)

func Init(router *router.Router) {
	tag := "debug:调试模块"
	r := router.Group("/rest/debug")
	r.Use(middleware.AuthUsernameAndPwd())
	{
		r.Post("/db/stat", db).Openapi.Tag(tag).Summary("db debug").ReqParam(dbParams{})
		r.Post("/db/ping", dbPing).Openapi.Tag(tag).Summary("db debug ping").ReqParam(dbParams{})
	}
}

type dbParams struct {
}

func db(ctx *context.Context) {
	params := dbParams{}
	ctx.BindForm(&params)
	dd := sqlkit.DefaultDataSource()
	ctx.JsonSuccess(dd.DBPool.Stats())
}

type dbPingParams struct {
	Phone string `comment:"手机号" default:"" trim:"true"`
	Pwd   string `validate:"required"`
}

func dbPing(ctx *context.Context) {
	params := dbPingParams{}
	ctx.BindForm(&params)
	dd := sqlkit.DefaultDataSource()
	err := dd.DBPool.Ping()
	if err != nil {
		logkit.Error(err.Error())
	}
	ctx.JsonSuccess(nil)
}
