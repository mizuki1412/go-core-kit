package download

import (
	"github.com/mizuki1412/go-core-kit/v2/class"
	"github.com/mizuki1412/go-core-kit/v2/class/exception"
	"github.com/mizuki1412/go-core-kit/v2/cli/configkey"
	"github.com/mizuki1412/go-core-kit/v2/library/filekit"
	"github.com/mizuki1412/go-core-kit/v2/service/configkit"
	"github.com/mizuki1412/go-core-kit/v2/service/restkit/context"
	"github.com/mizuki1412/go-core-kit/v2/service/restkit/middleware"
	"github.com/mizuki1412/go-core-kit/v2/service/restkit/openapi"
	"github.com/mizuki1412/go-core-kit/v2/service/restkit/router"
	"github.com/mizuki1412/go-core-kit/v2/service/storagekit"
	"strings"
)

func Init(router *router.Router) {
	tag := "common:公共模块"
	r := router.Group("/common")
	r.Use(middleware.AuthJWT())
	{
		r.Get("/download/private", download).Api(openapi.Tag(tag),
			openapi.Summary("私有下载"), openapi.ReqParam(downloadParams{}), openapi.ResponseStream())
		r.Post("/upload/private", upload).Api(openapi.Tag(tag),
			openapi.Summary("私有上传"), openapi.ReqBody(uploadParams{}))
		r.Post("/file/list", fileList).Api(openapi.Tag(tag),
			openapi.Summary("文件列表"), openapi.ReqParam(fileListParams{}))
		r.Post("/file/del", fileDel).Api(openapi.Tag(tag),
			openapi.Summary("文件删除"), openapi.ReqParam(fileListParams{}))
	}
	r2 := router.Group("/common")
	{
		r2.Get("/download", downloadPublic).Api(openapi.Tag(tag),
			openapi.Summary("公共下载"), openapi.ReqParam(downloadParams{}), openapi.ResponseStream())
	}
}

type downloadParams struct {
	Name string `validate:"required"`
}

func download(ctx *context.Context) {
	params := downloadParams{}
	ctx.BindForm(&params)
	subDir := configkit.GetString(configkey.ProjectSubDir4PrivateDownload)
	if subDir == "" {
		// 默认开启
		subDir = "."
	}
	if params.Name[0] == '/' {
		params.Name = params.Name[1:]
	}
	subs := strings.Split(params.Name, "/")
	if subDir == "" || (subDir != "." && (len(subs) == 1 || strings.Index(subDir, subs[0]) < 0)) {
		panic(exception.New("未授权开放目录"))
	}
	ctx.File2(params.Name)
}

func downloadPublic(ctx *context.Context) {
	params := downloadParams{}
	ctx.BindForm(&params)
	subDir := configkit.GetString(configkey.ProjectSubDir4PublicDownload)
	if params.Name[0] == '/' {
		params.Name = params.Name[1:]
	}
	subs := strings.Split(params.Name, "/")
	if subDir == "" || (subDir != "." && (len(subs) == 1 || strings.Index(subDir, subs[0]) < 0)) {
		panic(exception.New("未授权开放目录"))
	}
	ctx.File2(params.Name)
}

type uploadParams struct {
	File class.File   `validate:"required"`
	Path class.String `comment:"相对项目目录地址"`
}

func upload(ctx *context.Context) {
	params := uploadParams{}
	ctx.BindForm(&params)
	if !params.Path.Valid {
		params.Path.Set("/")
	}
	if params.Path.String[len(params.Path.String)-1] != '/' {
		params.Path.String += "/"
	}
	storagekit.SaveInHome(&params.File, params.Path.String+params.File.Header.Filename)
	ctx.JsonSuccess()
}

type fileListParams struct {
	Path string `comment:"相对项目目录地址" validate:"required"`
}

func fileList(ctx *context.Context) {
	params := fileListParams{}
	ctx.BindForm(&params)
	fullPath := storagekit.GetFullPath(params.Path)
	files := filekit.ListFileNames(fullPath)
	ret := make([]string, 0, len(files))
	if fullPath[len(fullPath)-1] != '/' {
		fullPath += "/"
	}
	for _, e := range files {
		ret = append(ret, strings.ReplaceAll(e, fullPath, ""))
	}
	ctx.JsonSuccess(ret)
}

func fileDel(ctx *context.Context) {
	params := fileListParams{}
	ctx.BindForm(&params)
	fullPath := storagekit.GetFullPath(params.Path)
	err := filekit.DelFile(fullPath)
	if err != nil {
		panic(exception.New(err.Error()))
	}
	ctx.JsonSuccess()
}
