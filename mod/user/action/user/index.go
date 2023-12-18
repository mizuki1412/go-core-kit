package user

import (
	"github.com/mizuki1412/go-core-kit/service/restkit/middleware"
	"github.com/mizuki1412/go-core-kit/service/restkit/router"
)

func Init(router *router.Router) {
	tag := "user:用户模块"
	router.Group("/rest/user/loginByUsername").Post("", loginByUsername).Openapi.Tag(tag).Summary("登录-用户名").ReqParam(loginByUsernameParam{})
	router.Group("/rest/user/login").Post("", login).Openapi.Tag(tag).Summary("登录").ReqParam(loginParam{})
	router.Group("/rest/user/info").Use(middleware.AuthJWT()).Post("", info).Openapi.Tag(tag).Summary("用户信息")
	r := router.Group("/rest/user", middleware.AuthJWT())
	{
		r.Post("/logout", logout).Openapi.Tag(tag).Summary("登出")
		r.Post("/updatePwd", updatePwd).Openapi.Tag(tag).Summary("密码修改").ReqParam(updatePwdParam{})
		r.Post("/updateUserInfo", updateUserInfo).Openapi.Tag(tag).Summary("更新用户信息").ReqParam(updateUserInfoParam{})
	}
	r1 := router.Group("/rest/user/admin", middleware.AuthJWT())
	{
		r1.Post("/list", listUsers).Openapi.Tag(tag).Summary("用户列表").ReqParam(listUsersParams{})
		r1.Post("/listByRole", listByRole).Openapi.Tag(tag).Summary("用户列表 by role").ReqParam(listByRoleParams{})
		r1.Post("/info", infoAdmin).Openapi.Tag(tag).Summary("用户信息").ReqParam(infoAdminParams{})
	}
	r2 := router.Group("/rest/user/admin", middleware.AuthJWT())
	{
		r2.Post("/add", AddUser).Openapi.Tag(tag).Summary("添加用户").ReqParam(AddUserParams{})
		r2.Post("/update", UpdateUser).Openapi.Tag(tag).Summary("修改用户").ReqParam(UpdateParams{})
		r2.Post("/del", DelUser).Openapi.Tag(tag).Summary("删除冻结用户").ReqParam(DelParams{})
	}
}
