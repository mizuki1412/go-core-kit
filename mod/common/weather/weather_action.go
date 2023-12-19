package weather

import (
	"github.com/mizuki1412/go-core-kit/service-third/locationkit"
	"github.com/mizuki1412/go-core-kit/service/restkit/context"
	"github.com/mizuki1412/go-core-kit/service/restkit/openapi"
	"github.com/mizuki1412/go-core-kit/service/restkit/router"
)

func Init(router *router.Router) {
	tag := "common:公共模块"
	r := router.Group("/rest/common")
	{
		r.Post("/weather", weatherInfo).Api(openapi.Tag(tag), openapi.Summary("获取天气信息"), openapi.ReqParam(weatherInfoParams{}))
	}
}

type weatherInfoParams struct {
	CityCode string `validate:"required"`
}

func weatherInfo(ctx *context.Context) {
	params := weatherInfoParams{}
	ctx.BindForm(&params)
	ctx.JsonSuccess(locationkit.Weather(params.CityCode))
}
