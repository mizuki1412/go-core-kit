package admindivision

import (
	"github.com/mizuki1412/go-core-kit/class"
	"github.com/mizuki1412/go-core-kit/mod/common/admindivision/dao/areadao"
	"github.com/mizuki1412/go-core-kit/mod/common/admindivision/dao/provincedao"
	"github.com/mizuki1412/go-core-kit/mod/common/admindivision/model"
	"github.com/mizuki1412/go-core-kit/service/restkit/context"
	"github.com/mizuki1412/go-core-kit/service/restkit/openapi"
	"github.com/mizuki1412/go-core-kit/service/restkit/router"
)

func Init(router *router.Router) {
	tag := "common:公共模块"
	r := router.Group("/rest/common/administrative")
	{
		r.GetPost("/listAllProvinceCity", ListAllProvinceCity).Api(openapi.Tag(tag), openapi.Summary("列表所有的省市"), openapi.Response([]*model.Province{}))
		r.GetPost("/listAreaByCity", listArea).Api(openapi.Tag(tag), openapi.Summary("按市列出区"),
			openapi.ReqParam(listAreaParam{}), openapi.Response([]*model.Area{}))
	}
}

func ListAllProvinceCity(ctx *context.Context) {
	ctx.JsonSuccess(provincedao.New(provincedao.ResultDefault).ListAll())
}

type listAreaParam struct {
	CityCode string `validate:"required"`
}

func listArea(ctx *context.Context) {
	params := listAreaParam{}
	ctx.BindForm(&params)
	ctx.JsonSuccess(areadao.New().ListByCity(class.NewString(params.CityCode)))
}
