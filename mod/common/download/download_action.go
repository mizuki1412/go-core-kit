package download

import (
	"github.com/mizuki1412/go-core-kit/class/exception"
	"github.com/mizuki1412/go-core-kit/init/configkey"
	"github.com/mizuki1412/go-core-kit/library/cryptokit"
	"github.com/mizuki1412/go-core-kit/mod/middleware"
	"github.com/mizuki1412/go-core-kit/service/configkit"
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
	r2 := router.Group("/rest/common")
	{
		r2.Post("/download", downloadPublic).Swagger.Tag(tag).Summary("公共下载").Param(downloadParams{})
		r2.Get("/download", downloadPublic).Swagger.Tag(tag).Summary("公共下载").Param(downloadParams{})
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
		filename = cryptokit.URLEncode(params.Name)
	}
	ctx.File(params.Name, filename)
}

func downloadPublic(ctx *context.Context) {
	params := downloadParams{}
	ctx.BindForm(&params)
	subDir := configkit.GetStringD(configkey.ProjectSubDir4PublicDownload)
	if params.Name[0] == '/' {
		params.Name = params.Name[1:]
	}
	subs := strings.Split(params.Name, "/")
	if subDir == "" || (subDir != "." && (len(subs) == 1 || strings.Index(subDir, subs[0]) < 0)) {
		panic(exception.New("未授权开放目录"))
	}
	dotIndex := strings.LastIndex(params.Name, ".")
	var filename string
	if dotIndex >= 0 {
		filename = cryptokit.URLEncode(params.Name[0:dotIndex]) + params.Name[dotIndex:]
	} else {
		filename = cryptokit.URLEncode(params.Name)
	}
	ctx.File(params.Name, filename)
}
