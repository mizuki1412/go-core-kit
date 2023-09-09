package download

import (
	"github.com/mizuki1412/go-core-kit/class"
	"github.com/mizuki1412/go-core-kit/class/exception"
	"github.com/mizuki1412/go-core-kit/cli/configkey"
	"github.com/mizuki1412/go-core-kit/library/filekit"
	"github.com/mizuki1412/go-core-kit/mod/middleware"
	"github.com/mizuki1412/go-core-kit/service/configkit"
	"github.com/mizuki1412/go-core-kit/service/restkit/context"
	"github.com/mizuki1412/go-core-kit/service/restkit/router"
	"github.com/mizuki1412/go-core-kit/service/storagekit"
	"strings"
)

func Init(router *router.Router) {
	tag := "common:公共模块"
	r := router.Group("/rest")
	r.Use(middleware.AuthUsernameAndPwd())
	{
		r.Post("/download", download).Openapi.Tag(tag).Summary("私有下载").ReqParam(downloadParams{}).ResponseStream()
		r.Get("/download", download).Openapi.Tag(tag).Summary("私有下载").ReqParam(downloadParams{}).ResponseStream()
		r.Post("/upload", upload).Openapi.Tag(tag).Summary("私有上传").ReqBody(uploadParams{})
		r.Post("/file/list", fileList).Openapi.Tag(tag).Summary("文件列表").ReqParam(fileListParams{})
		r.Post("/file/del", fileDel).Openapi.Tag(tag).Summary("文件删除").ReqParam(fileListParams{})
	}
	r2 := router.Group("/rest/common")
	{
		r2.Post("/download", downloadPublic).Openapi.Tag(tag).Summary("公共下载").ReqParam(downloadParams{}).ResponseStream()
		r2.Get("/download", downloadPublic).Openapi.Tag(tag).Summary("公共下载").ReqParam(downloadParams{}).ResponseStream()
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
	Path class.String `description:"相对项目目录地址"`
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
	ctx.JsonSuccess(nil)
}

type fileListParams struct {
	Path string `description:"相对项目目录地址" validate:"required"`
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
	ctx.JsonSuccess(nil)
}
