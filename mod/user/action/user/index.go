package user

import (
	"github.com/mizuki1412/go-core-kit/mod/middleware"
	"github.com/mizuki1412/go-core-kit/service/restkit/router"
)

func Init(router *router.Router) {
	tag := "user:用户模块"
	router.Group("/rest/user/loginByUsername").Use(middleware.CreateSession()).Post("", loginByUsername).Swagger.Tag(tag).Summary("登录-用户名").Param(loginByUsernameParam{})
	router.Group("/rest/user/login").Use(middleware.CreateSession()).Post("", login).Swagger.Tag(tag).Summary("登录").Param(loginParam{})
	router.Group("/rest/user/info").Use(middleware.AuthUsernameAndPwd()).Post("", info).Swagger.Tag(tag).Summary("用户信息")
	r := router.Group("/rest/user", middleware.AuthUsernameAndPwd())
	{
		r.Post("/logout", logout).Swagger.Tag(tag).Summary("登出")
		r.Post("/updatePwd", updatePwd).Swagger.Tag(tag).Summary("密码修改").Param(updatePwdParam{})
		r.Post("/updateUserInfo", updateUserInfo).Swagger.Tag(tag).Summary("更新用户信息").Param(updateUserInfoParam{})
	}
	r1 := router.Group("/rest/user/admin", middleware.AuthUsernameAndPwd())
	{
		r1.Post("/list", listUsers).Swagger.Tag(tag).Summary("用户列表").Param(listUsersParams{})
		r1.Post("/listByRole", listByRole).Swagger.Tag(tag).Summary("用户列表 by role").Param(listByRoleParams{})
		r1.Post("/info", infoAdmin).Swagger.Tag(tag).Summary("用户信息").Param(infoAdminParams{})
	}
	r2 := router.Group("/rest/user/admin", middleware.AuthUsernameAndPwd())
	{
		r2.Post("/add", AddUser).Swagger.Tag(tag).Summary("添加用户").Param(AddUserParams{})
		r2.Post("/update", UpdateUser).Swagger.Tag(tag).Summary("修改用户").Param(UpdateParams{})
		r2.Post("/del", DelUser).Swagger.Tag(tag).Summary("删除冻结用户").Param(DelParams{})
	}
}
