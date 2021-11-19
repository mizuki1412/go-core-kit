package download

import (
	"github.com/mizuki1412/go-core-kit/library/cryptokit"
	"github.com/mizuki1412/go-core-kit/mod/middleware"
	"github.com/mizuki1412/go-core-kit/service/restkit/context"
	"github.com/mizuki1412/go-core-kit/service/restkit/router"
	"strings"
)

func Init(router *router.Router) {
	tag := "common:公共模块"
	r := router.Group("/rest")
	r.Use(middleware.AuthUsernameAndPwd())
	{
		r.Post("/download", download).Swagger.Tag(tag).Summary("私有下载").Param(downloadParams{})
		r.Get("/download", download).Swagger.Tag(tag).Summary("私有下载").Param(downloadParams{})
	}
}

type downloadParams struct {
	Name string `validate:"required"`
}

func download(ctx *context.Context) {
	params := downloadParams{}
	ctx.BindForm(&params)
	dotIndex := strings.LastIndex(params.Name, ".")
	var filename string
	if dotIndex >= 0 {
		filename = cryptokit.URLEncode(params.Name[0:dotIndex]) + params.Name[dotIndex:]
	} else {
		filename = params.Name
	}
	ctx.File(params.Name, filename)
}
