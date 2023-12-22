package user

import (
	"github.com/mizuki1412/go-core-kit/service/restkit/context"
	"github.com/mizuki1412/go-core-kit/service/restkit/middleware"
	"github.com/mizuki1412/go-core-kit/service/restkit/openapi"
	"github.com/mizuki1412/go-core-kit/service/restkit/router"
	"log"
)

func Init(router *router.Router) {
	tag := "user:用户模块"
	router.Group("/rest/user/loginByUsername").Post("", loginByUsername).Api(openapi.Tag(tag),
		openapi.Summary("登录-用户名"), openapi.ReqParam(loginByUsernameParam{}))
	router.Group("/rest/user/login").Post("", login).Api(openapi.Tag(tag),
		openapi.Summary("登录"), openapi.ReqParam(loginParam{}))
	router.Group("/rest/user/info").Use(middleware.AuthJWT()).Post("", info).Api(openapi.Tag(tag),
		openapi.Summary("用户信息"))
	r := router.Group("/rest/user", middleware.AuthJWT())
	{
		r.Post("/logout", logout).Api(openapi.Tag(tag), openapi.Summary("登出"))
		r.Post("/updatePwd", updatePwd).Api(openapi.Tag(tag), openapi.Summary("密码修改"), openapi.ReqParam(updatePwdParam{}))
		r.Post("/updateUserInfo", updateUserInfo).Api(openapi.Tag(tag), openapi.Summary("更新用户信息"), openapi.ReqParam(updateUserInfoParam{}))
	}
	r1 := router.Group("/rest/user/admin", middleware.AuthJWT())
	{
		r1.Post("/list", listUsers).Api(openapi.Tag(tag),
			openapi.Summary("用户列表"), openapi.ReqParam(listUsersParams{}))
		r1.Post("/listByRole", listByRole).Api(openapi.Tag(tag),
			openapi.Summary("用户列表 by role"), openapi.ReqParam(listByRoleParams{}))
		r1.Post("/info", infoAdmin).Api(openapi.Tag(tag),
			openapi.Summary("用户信息"), openapi.ReqParam(infoAdminParams{}))
	}
	r2 := router.Group("/rest/user/admin", middleware.AuthJWT())
	{
		r2.Post("/add", AddUser).Api(openapi.Tag(tag), openapi.Summary("添加用户"), openapi.ReqParam(AddUserParams{}))
		r2.Post("/update", UpdateUser).Api(openapi.Tag(tag), openapi.Summary("修改用户"), openapi.ReqParam(UpdateParams{}))
		r2.Post("/del", DelUser).Api(openapi.Tag(tag), openapi.Summary("删除冻结用户"), openapi.ReqParam(DelParams{}))
	}

	// todo test
	router.Get("/rest/test2/:id", test2)
}

type test2Param struct {
	Id int32
}

func test2(ctx *context.Context) {
	params := test2Param{}
	ctx.BindForm(&params)
	log.Println(params.Id)
	ctx.JsonSuccess()
}
