package admindivision

import (
	"github.com/mizuki1412/go-core-kit/class"
	"github.com/mizuki1412/go-core-kit/mod/common/admindivision/dao/areadao"
	"github.com/mizuki1412/go-core-kit/mod/common/admindivision/dao/provincedao"
	"github.com/mizuki1412/go-core-kit/service/restkit/context"
	"github.com/mizuki1412/go-core-kit/service/restkit/router"
)

func Init(router *router.Router) {
	tag := "common:公共模块"
	r := router.Group("/rest/common/administrative")
	{
		r.Post("/listAllProvinceCity", ListAllProvinceCity).Swagger.Tag(tag).Summary("列表所有的省市")
		r.Get("/listAllProvinceCity", ListAllProvinceCity).Swagger.Tag(tag).Summary("列表所有的省市")
		r.Post("/listAreaByCity", listArea).Swagger.Tag(tag).Summary("按市列出区").Param(listAreaParam{})
		r.Get("/listAreaByCity", listArea).Swagger.Tag(tag).Summary("按市列出区").Param(listAreaParam{})
	}
}

func ListAllProvinceCity(ctx *context.Context) {
	ctx.JsonSuccess(provincedao.New().ListAll())
}

type listAreaParam struct {
	CityCode string `validate:"required"`
}

func listArea(ctx *context.Context) {
	params := listAreaParam{}
	ctx.BindForm(&params)
	ctx.JsonSuccess(areadao.New().ListByCity(class.String{String: params.CityCode, Valid: true}))
}
